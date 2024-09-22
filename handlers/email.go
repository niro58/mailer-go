package handler

import (
	contract "email-sender/contracts"
	"os"

	"github.com/gin-gonic/gin"
)

func (a *App) Send(c *gin.Context){
	
	auth := c.Request.Header.Get("Authorization")
	apiAuth := os.Getenv("API_AUTH")
	if auth != apiAuth {
		c.JSON(401, CreateReply(
			nil,
			ErrUnauthorized,
		))
			return
	}
	type Request struct {
		Subject string `json:"subject"`
		Body string `json:"body"`
	}

	var req Request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, CreateReply(
			nil,
			err,
		))
		return
	}
	email := contract.Email{
		To:  os.Getenv("EMAIL_TO"),
		Subject: req.Subject,
		Body: req.Body,
	}
	err := a.EmailService.Send(email)

	if err != nil {
		c.JSON(200, CreateReply(
			nil,
			err,
		))
		return
	}
	
	c.JSON(200, CreateReply(
		nil,
		nil,
	))
}

