package request

type SyncAccount struct {
	DeviceName string `json:"device_name"`
}

type BlastMessage struct {
	DeviceName string   `json:"device_name"`
	Max        int64    `json:"max"`
	Message    []string `json:"message"`
}

type Message struct {
	DeviceName string `json:"device_name"`
	To         string `json:"to"`
	Message    string `json:"message"`
}

type BlastMsg struct {
	DeviceName string `json:"device_name"`
	Max        int64  `json:"max"`
	Message    string `json:"message"`
}
