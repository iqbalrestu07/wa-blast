package services

import (
	"wa-blast/loggers"
	"wa-blast/repositories"
)

// Logger
var log = loggers.Get()

// User ...
var User UserService
var Message MessageService
var Whatsapp WhatsappService

// Init ...
func Init() {
	User = &userService{repositories.User}
	Message = &messageService{repositories.User}
	Whatsapp = &whatsappService{repositories.User, repositories.Message}
}
