package mail

import (
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "INFOBANKDATA <noreplay@gmail.com>"
const CONFIG_AUTH_EMAIL = "mihidev47@gmail.com"
const CONFIG_AUTH_PASSWORD = "qxnlhizxpsfkbsbl"

func SendMail(to, title, msg string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", title)
	mailer.SetBody("text/html", msg)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return err
}
