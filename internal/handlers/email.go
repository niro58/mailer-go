package handler

import (
	contract "mailer-go/internal/contracts"
	"os"

	"github.com/gin-gonic/gin"
)

func (a *App) Health(c *gin.Context) {
	c.JSON(200, CreateReply(
		"ok2",
		nil,
	))
}
func (a *App) Send(c *gin.Context) {
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
		Sender    string `json:"sender"`
		Reason    string `json:"reason"`
		Subject   string `json:"subject"`
		Body      string `json:"body"`
		Recipient string `json:"recipient"`
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

	err := a.EmailService.Send(req.Sender, req.Reason, email, req.Recipient)

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
