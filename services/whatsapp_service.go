package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
	"wa-blast/configs"
	"wa-blast/models"
	"wa-blast/repositories"
	"wa-blast/request"
	"wa-blast/util"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// UserService ...
type WhatsappService interface {
	SyncAccount(id string, data models.CompaniesDevice)
	SendMessage(company, id string, req request.Message) (string, error)
	Blast(company, id string, req request.BlastMessage)
}

type whatsappService struct {
	user    repositories.UserRepository
	message repositories.MessageRepository
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func (s *whatsappService) SyncAccount(id string, data models.CompaniesDevice) {
	var qr string
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "./db/"+id+".db", dbLog)
	if err != nil {
		fmt.Println(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		fmt.Println(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		go func() {
			// No ID stored, new login
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				fmt.Println(err)
			}

			for evt := range qrChan {
				switch evt.Event {
				case "success":
					{
						s.ActivationAccount(client, id, qr, data)
					}
				case "timeout":
					{
						s.user.DeleteDevice(id)
					}
				case "code":
					{
						qr = evt.Code
						err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, configs.MustGetString("file.images")+"/"+id+".png")
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}()
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			fmt.Println(err)
		}
	}

}

func (s *whatsappService) ActivationAccount(client *whatsmeow.Client, ID, qr string, data models.CompaniesDevice) {
	data.Active = true
	data.Qrcode = qr
	data.SyncDatetime = time.Now()
	data.LastActivity = time.Now()

	s.user.UpdateDevice(ID, data)

	os.Remove(configs.MustGetString("file.images") + "/" + ID + ".png")
}

func (s *whatsappService) SendMessage(company, id string, req request.Message) (string, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "./db/"+id+".db", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		return "", util.NewError("-1002")
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			return "", util.NewError("-1002")
		}
	}
	msgID := whatsmeow.GenerateMessageID()
	checkValid, _ := client.IsOnWhatsApp([]string{"+" + util.ClearNumber(req.To)})
	if len(checkValid) > 0 {
		if checkValid[0].IsIn {

			send, err := client.SendMessage(context.Background(), types.NewJID(util.ClearNumber(req.To), util.TypeTarget(req.To)), &waProto.Message{
				Conversation: proto.String(req.Message),
			})
			if err != nil {
				s.message.InsertMessageOutbox(models.Outbox{
					ID:           msgID,
					IDCompanies:  company,
					IDDevice:     id,
					To:           req.To,
					Message:      req.Message,
					MsgSuccess:   err.Error(),
					IsSending:    false,
					SendDatetime: time.Now(),
				})
				client.Disconnect()
				return "", util.NewError("-1003")
			}

			sendResp, _ := json.Marshal(send)
			s.message.InsertMessageOutbox(models.Outbox{
				ID:           msgID,
				IDCompanies:  company,
				IDDevice:     id,
				To:           req.To,
				Message:      req.Message,
				MsgSuccess:   string(sendResp),
				IsSending:    true,
				SendDatetime: time.Now(),
			})

			client.Disconnect()
			return msgID, nil
		}
	}

	return "", errors.New("Phone not found")
}

func (s *whatsappService) Blast(company, id string, req request.BlastMessage) {
	rand.Seed(time.Now().Unix())
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "./db/"+id+".db", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	targetTelk := []string{"6281211"}
	minNumber := 10000
	maxNumber := 999999
	max := 0
	for {
		if max >= int(req.Max) {
			break
		}
		number := fmt.Sprintf("%d", rand.Intn(maxNumber-minNumber+1)+minNumber)
		for _, t := range targetTelk {
			genPhone := t + number
			msgID := whatsmeow.GenerateMessageID()
			checkValid, _ := client.IsOnWhatsApp([]string{"+" + util.ClearNumber(genPhone)})
			if len(checkValid) > 0 {
				if checkValid[0].IsIn {
					msg := req.Message[rand.Intn(len(req.Message))]
					send, err := client.SendMessage(context.Background(), types.NewJID(util.ClearNumber(genPhone), util.TypeTarget(genPhone)), &waProto.Message{
						Conversation: proto.String(msg),
					})
					if err != nil {
						fmt.Println(util.ClearNumber(genPhone) + " : Not Found")
					}

					sendResp, _ := json.Marshal(send)
					s.message.InsertMessageOutbox(models.Outbox{
						ID:           msgID,
						IDCompanies:  company,
						IDDevice:     id,
						To:           genPhone,
						Message:      msg,
						MsgSuccess:   string(sendResp),
						IsSending:    true,
						SendDatetime: time.Now(),
					})
					fmt.Println(util.ClearNumber(genPhone) + " : Found")
					max++
				}
			}
		}

	}

	client.Disconnect()
}
