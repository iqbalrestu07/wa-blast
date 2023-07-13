package models

import "time"

type Outbox struct {
	ID           string
	IDCompanies  string
	IDDevice     string
	To           string
	Message      string
	MsgSuccess   string
	MsgError     string
	IsSending    bool
	SendDatetime time.Time
}
