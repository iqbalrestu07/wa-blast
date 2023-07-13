package services

import (
	"wa-blast/configs"
	"wa-blast/models"
	"wa-blast/repositories"
	"wa-blast/request"
	"wa-blast/response"
	"wa-blast/util"
)

// MessageService ...
type MessageService interface {
	SyncAccount(id string, req request.SyncAccount) (*response.Success, error)
	Message(id string, req request.Message) (*response.Success, error)
	Blast(id string, req request.BlastMessage) (*response.Success, error)
}

type messageService struct {
	user repositories.UserRepository
}

func (s *messageService) SyncAccount(id string, req request.SyncAccount) (*response.Success, error) {

	device := models.CompaniesDevice{
		ID:              util.NewId(),
		IDCompaniesAuth: id,
		DeviceName:      req.DeviceName,
		Active:          false,
	}

	s.user.CreateDevice(device)

	defer Whatsapp.SyncAccount(device.ID, device)

	return &response.Success{
		Result: &response.Device{
			DeviceName: device.DeviceName,
			Url:        configs.MustGetString("server.url") + "/files/" + device.ID + ".png",
		},
	}, nil
}

func (s *messageService) Message(id string, req request.Message) (*response.Success, error) {

	device, _ := s.user.GetDevice("device_name = ? AND active = ?", req.DeviceName, true)
	if device.ID == "" {
		return nil, util.NewError("-1002")
	}

	msgID, err := Whatsapp.SendMessage(id, device.ID, req)

	return &response.Success{
		Result: msgID,
	}, err
}

func (s *messageService) Blast(id string, req request.BlastMessage) (*response.Success, error) {

	device, _ := s.user.GetDevice("device_name = ? AND active = ?", req.DeviceName, true)
	if device.ID == "" {
		return nil, util.NewError("-1002")
	}
	go func(id string, device string, req request.BlastMessage) {
		Whatsapp.Blast(id, device, req)
	}(id, device.ID, req)

	return &response.Success{
		Result: "running blast..",
	}, nil
}
