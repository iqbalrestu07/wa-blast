package models

import "time"

type CompaniesAuth struct {
	ID         string
	Perusahaan string
	SecretKey  string
	Active     bool
}

type CompaniesDevice struct {
	ID              string
	IDCompaniesAuth string
	DeviceName      string
	Qrcode          string
	Active          bool
	SyncDatetime    time.Time
	LastActivity    time.Time
}
