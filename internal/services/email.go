package service

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	contract "mailer-go/internal/contracts"
	"net/smtp"
	"os"
	"path"
	"sync"
)

type ClientConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Mutex    *sync.Mutex
}
type EmailService struct {
	Configs map[string]*ClientConfig
}

func getClientConfigs() (map[string]*ClientConfig, error) {
	dirname, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	jsonFile, err := os.Open(path.Join(dirname, "/clients.json"))
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]*ClientConfig
	json.Unmarshal([]byte(byteValue), &result)
	for _, config := range result {
		config.Mutex = &sync.Mutex{}
	}

	return result, nil
}

var ErrClientNotFound = errors.New("client not found")

func createSMTPClient(sender ClientConfig) (*smtp.Client, error) {
	tlsConfig := &tls.Config{
		ServerName: sender.Host,
	}

	conn, err := tls.Dial("tcp", sender.Host+":"+sender.Port, tlsConfig)
	if err != nil {
		return nil, err
	}

	client, err := smtp.NewClient(conn, sender.Host)
	if err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", sender.Username, sender.Password, sender.Host)
	if err := client.Auth(auth); err != nil {
		return nil, err
	}
	if err := client.Mail(sender.Username); err != nil {
		return nil, err
	}
	return client, nil
}

func NewEmailService() EmailService {
	configs, err := getClientConfigs()
	if err != nil {
		panic(err)
	}
	return EmailService{configs}
}

func (e EmailService) Send(senderKey string, recipient string, email contract.Email) error {
	config, exists := e.Configs[senderKey]
	if !exists {
		return ErrClientNotFound
	}

	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	client, err := createSMTPClient(*config)
	if err != nil {
		return err
	}
	if err := client.Rcpt(recipient); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	msg := "From: " + config.Host + "\r\n" +
		"To: " + recipient + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" +
		email.Body + "\r\n"

	if _, err = w.Write([]byte(msg)); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	return nil
}
