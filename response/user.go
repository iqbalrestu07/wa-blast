package response

type Device struct {
	DeviceName string `json:"device_name"`
	Url        string `json:"qrcode_url"`
}
