package test

import (
	contract "email-sender/contracts"
	service "email-sender/services"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestEmailSend(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	emailService := service.NewEmailService()
	email := contract.Email{
		Subject: "Test",
		Body:    "Test",
	}
	err = emailService.Send(email)
	if err != nil {
		fmt.Println(err)
		t.Errorf("Error sending email")
	}

}