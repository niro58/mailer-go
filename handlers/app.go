package handler

import (
	contract "email-sender/contracts"
	service "email-sender/services"
)

type App struct {
	EmailService contract.EmailService
}

func CreateApplication () App {
	var app App
	app.EmailService = service.NewEmailService()

	return app
}