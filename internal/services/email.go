package service

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	contract "mailer-go/internal/contracts"
	"net/smtp"
	"os"
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

type Job struct {
	Email  contract.Email
	Config *ClientConfig
}

var (
	workers = 10
)

type EmailService struct {
	Configs   map[string]*ClientConfig
	Templates map[string]Template
	jobs      chan Job
	wg        sync.WaitGroup
}

func getClientConfigs() (map[string]*ClientConfig, error) {
	clientsPath := "./clients.json"
	fmt.Println("Loading client configs from", clientsPath)
	jsonFile, err := os.Open(clientsPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]*ClientConfig
	json.Unmarshal([]byte(byteValue), &result)

	fmt.Println("Loaded client configs:", len(result))
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

func NewEmailService() *EmailService {
	configs, err := getClientConfigs()
	if err != nil {
		panic(err)
	}
	templates, err := getTemplates()
	if err != nil {
		panic(err)
	}

	return &EmailService{
		Configs:   configs,
		Templates: templates,
		jobs:      make(chan Job, workers),
	}
}

func (e *EmailService) StartPool() {
	for w := 0; w < workers; w++ {
		e.wg.Add(1)
		go e.Send(w, e.jobs, &e.wg)
	}
}

func (e *EmailService) AddJob(email contract.Email) error {
	config, exists := e.Configs[email.SenderKey]
	if !exists {
		return ErrClientNotFound
	}

	e.jobs <- Job{email, config}

	return nil
}

func (e *EmailService) AddTemplateJob(template contract.EmailTemplate) error {
	email := contract.Email{
		EmailHeaders: template.EmailHeaders,
	}

	templateConfig, ok := e.Templates[template.TemplateKey]
	if !ok {
		return ErrTemplateNotFound
	}
	config, exists := e.Configs[template.SenderKey]
	if !exists {
		return ErrClientNotFound
	}

	err := templateConfig.Validate(template.Variables)
	if err != nil {
		return err
	}

	for k, v := range template.Variables {
		k = "{{" + k + "}}"
		templateConfig.Body = strings.ReplaceAll(templateConfig.Body, k, v)
		templateConfig.Subject = strings.ReplaceAll(templateConfig.Subject, k, v)
	}

	email.Subject = templateConfig.Subject
	email.Body = templateConfig.Body
	e.jobs <- Job{email, config}

	return nil
}

func (e *EmailService) Send(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Println("Worker", id, "sending email to", job.Email.Recipients)

		client, err := createSMTPClient(*job.Config)
		if err != nil {
			fmt.Println("Error creating SMTP client:", err)
			continue
		}

		for _, recipient := range job.Email.Recipients {
			if err := client.Rcpt(recipient); err != nil {
				fmt.Println("Error adding recipient:", err)
				continue
			}
		}

		w, err := client.Data()
		if err != nil {
			fmt.Println("Error getting data writer:", err)
			continue
		}
		msg := fmt.Sprintf(
			"From: %s\r\n"+
				"To: %s\r\n"+
				"Bcc: %s\r\n"+
				"Subject: %s\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: %s; charset=\"UTF-8\"\r\n\r\n"+
				"%s\r\n",
			job.Config.Host,
			strings.Join(job.Email.Recipients, ", "),
			strings.Join(job.Email.Bcc, ", "),
			job.Email.Subject,
			job.Email.ContentType,
			job.Email.Body)

		if _, err = w.Write([]byte(msg)); err != nil {
			fmt.Println("Error writing message:", err)
			continue
		}
		if err = w.Close(); err != nil {
			fmt.Println("Error closing writer:", err)
			continue
		}

		client.Quit()
	}
}

func (e *EmailService) Count() int {
	return len(e.jobs)
}

func (e *EmailService) Wait() {
	e.wg.Wait()
}

func (e *EmailService) Shutdown() {
	close(e.jobs)
	e.Wait()
}
