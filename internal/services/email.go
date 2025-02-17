package service

import (
	"crypto/tls"
	"mailer-go/internal/config"
	contract "mailer-go/internal/contracts"
	"net/smtp"
)

type EmailService struct {
	Clients map[string]*smtp.Client
}

func CreateClients() map[string]*smtp.Client {
	clients := make(map[string]*smtp.Client)
	for key, sender := range config.Senders {
		tlsConfig := &tls.Config{
			ServerName: sender.Host,
		}
		conn, err := tls.Dial("tcp", sender.Host+":"+sender.Port, tlsConfig)
		if err != nil {
			panic(err)
		}

		client, err := smtp.NewClient(conn, sender.Host)
		if err != nil {
			panic(err)
		}

		auth := smtp.PlainAuth("", sender.Username, sender.Password, sender.Host)
		if err := client.Auth(auth); err != nil {
			panic(err)
		}
		if err := client.Mail(sender.Username); err != nil {
			panic(err)
		}

		clients[key] = client
	}
	return clients
}
func NewEmailService() EmailService {
	return EmailService{Clients: CreateClients()}
}

func (e EmailService) Send(sender string, contactReason string, email contract.Email, recipient string) error {
	client := e.Clients[sender]
	if client == nil {
		return config.ErrSenderNotFound
	}
	if err := client.Rcpt(recipient); err != nil {
		panic(err)
	}

	w, err := client.Data()
	if err != nil {
		panic(err)
	}
	username := ""
	msg := "From: " + username + "\r\n" +
		"To: " + recipient + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" + // Blank line to separate headers from body
		email.Body + "\r\n"

	_, err = w.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	client.Quit()
	return nil
}
