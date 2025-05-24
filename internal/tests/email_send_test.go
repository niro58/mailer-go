package test

import (
	"fmt"
	contract "mailer-go/internal/contracts"
	env "mailer-go/internal/environment"
	service "mailer-go/internal/services"
	"testing"
)

func TestEmailSend(t *testing.T) {
	env.NewEnv()

	emailService := service.NewEmailService()
	emailService.StartPool()

	email := contract.Email{
		EmailHeaders: contract.EmailHeaders{
			SenderKey:  "invoice-recogniser",
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
