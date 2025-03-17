package test

import (
	"fmt"
	contract "mailer-go/internal/contracts"
	service "mailer-go/internal/services"
	util "mailer-go/internal/utils"
	"path"
	"testing"

	"github.com/joho/godotenv"
)

func Configure(t *testing.T) {
	err := godotenv.Load(path.Join(util.Root, "/.env"))
	if err != nil {
		t.Errorf("Error loading .env file")
	}
}
func TestEmailSend(t *testing.T) {
	Configure(t)

	emailService := service.NewEmailService()
	emailService.StartPool()

	email := contract.Email{
		EmailHeaders: contract.EmailHeaders{
			SenderKey:  "lucky",
			Recipients: []string{"nichita.roilean@gmail.com"},
		},
		ContentType: "text/html",
		Subject:     "Test",
		Body:        "<h1>This is a test email</h1><p>With HTML content.</p>",
	}

	emailService.AddJob(email)

	fmt.Println(emailService.Count())

	emailService.Shutdown()
}
