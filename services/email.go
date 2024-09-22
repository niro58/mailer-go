package service

import (
	contract "email-sender/contracts"
	"net/smtp"
	"os"
)

type EmailService struct {
	auth smtp.Auth
	sendTo string
}

func NewEmailService() EmailService {
	auth := smtp.PlainAuth("", os.Getenv("TO_EMAIL"), os.Getenv("MAILTRAP_SECRET_KEY"), "sandbox.smtp.mailtrap.io")
	to := os.Getenv("TO_EMAIL")
	return EmailService{auth, to}
}

func (e EmailService) Send(email contract.Email) error {
	sendTo := []string{e.sendTo}
	msg := []byte("To: " + e.sendTo + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" + email.Body + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", e.auth, e.sendTo, sendTo,msg)

	if err != nil {
		return err
	}

	return nil
}