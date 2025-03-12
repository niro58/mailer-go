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
	"strings"
	"sync"
)

var (
	ErrClientNotFound   = errors.New("client not found")
	ErrTemplateNotFound = errors.New("client not found")
)

type ClientConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	workers = 10
	jobs    = make(chan contract.Email, workers)
	wg      = sync.WaitGroup{}
)

type EmailService struct {
	Configs   map[string]*ClientConfig
	Templates map[string]Template
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

	return result, nil
}

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
	templates, err := getTemplates()
	if err != nil {
		panic(err)
	}

	return EmailService{configs, templates}
}

func (e EmailService) StartPool() {
	for w := 0; w <= workers; w++ {
		wg.Add(1)
		go e.Send(w, jobs, &wg)
	}
}
func (e EmailService) AddJob(email contract.Email) error {
	jobs <- email
	return nil
}
func (e EmailService) AddTemplateJob(template contract.EmailTemplate) error {
	var email contract.Email

	templateConfig, ok := e.Templates[template.TemplateKey]
	if !ok {
		return ErrTemplateNotFound
	}
	err := templateConfig.Validate(template.Variables)
	if err != nil {
		return err
	}

	for k, v := range template.Variables {
		templateConfig.Body = strings.ReplaceAll(templateConfig.Body, k, v)
		templateConfig.Subject = strings.ReplaceAll(templateConfig.Subject, k, v)
	}

	email.Subject = templateConfig.Subject
	email.Body = templateConfig.Body
	jobs <- email

	return nil
}
func (e *EmailService) Send(id int, jobs <-chan contract.Email, wg *sync.WaitGroup) error {
	defer wg.Done()

	for job := range jobs {
		config, exists := e.Configs[job.SenderKey]
		if !exists {
			return ErrClientNotFound
		}

		client, err := createSMTPClient(*config)
		if err != nil {
			return err
		}
		defer client.Quit()

		for _, recipient := range job.Recipients {
			if err := client.Rcpt(recipient); err != nil {
				return err
			}
		}

		w, err := client.Data()
		if err != nil {
			return err
		}
		msg := "From: " + config.Host + "\r\n" +
			"To: " + strings.Join(job.Recipients, ", ") + "\r\n" +
			"Bcc: " + strings.Join(job.Bcc, ", ") + "\r\n" +
			"Subject: " + job.Subject + "\r\n" +
			"\r\n" +
			job.Body + "\r\n"

		if _, err = w.Write([]byte(msg)); err != nil {
			return err
		}
		if err = w.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (EmailService) Count() int {
	return len(jobs)
}
