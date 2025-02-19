package handler

import (
	contract "mailer-go/internal/contracts"

	"github.com/gin-gonic/gin"
)

func (a *App) Health(c *gin.Context) {
	c.JSON(200, CreateReply(
		"ok",
		nil,
	))
}
func (a *App) Send(c *gin.Context) {
	type Request struct {
		SenderKey string `json:"senderKey"`
		Recipient string `json:"recipient"`
		Subject   string `json:"subject"`
		Body      string `json:"body"`
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
		Subject: req.Subject,
		Body:    req.Body,
	}

	err := a.EmailService.Send(req.SenderKey, req.Recipient, email)

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
