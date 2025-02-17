package handler

import (
	contract "mailer-go/internal/contracts"
	service "mailer-go/internal/services"
)

type App struct {
	EmailService contract.EmailService
}

func CreateApplication() App {
	var app App
	app.EmailService = service.NewEmailService()

	return app
}
