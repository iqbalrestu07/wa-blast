package midtrans

import (
	"wa-blast/configs"
	"wa-blast/loggers"

	"github.com/veritrans/go-midtrans"
)

var midtransClient *midtrans.Client

var MidtransCore *midtrans.CoreGateway

var SnapGateway *midtrans.SnapGateway

var Resp *midtrans.ChargeReqWithMap

var log = loggers.Get()

// Init ...
func Init() {
	midclient := midtrans.NewClient()
	midclient.ServerKey = configs.MustGetString("midtrans.serverkey")
	midclient.ClientKey = configs.MustGetString("midtrans.clientkey")
	midclient.APIEnvType = midtrans.Sandbox

	if configs.MustGetString("midtrans.environment") != "sandbox" {
		midclient.APIEnvType = midtrans.Production
	}

	coreGateway := midtrans.CoreGateway{
		Client: midclient,
	}

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	midtransClient = midtransClient
	MidtransCore = &coreGateway
	SnapGateway = &snapGateway

	log.Infof("Initiate midtrans payment on %s", configs.MustGetString("midtrans.environment"))
}

func ReqSnap(code string, price float64, username, phone, email string) (midtrans.SnapResponse, error) {

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  code,
			GrossAmt: int64(price),
		},
		CustomerDetail: &midtrans.CustDetail{
			FName: username,
			Email: email,
			Phone: phone,
		},
	}

	snapTokenResp, err := SnapGateway.GetToken(snapReq)

	return snapTokenResp, err
}
