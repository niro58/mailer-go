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
	var req contract.Email
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, CreateReply(
			nil,
			err,
		))
		return
	}

	a.EmailService.AddJob(req)

	c.JSON(200, CreateReply(
		nil,
		nil,
	))
}
func (a *App) Status(c *gin.Context) {
	c.JSON(200, CreateReply(
		a.EmailService.Count(),
		nil,
	))
}
func (a *App) SendTemplate(c *gin.Context) {
	var req contract.EmailTemplate
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, CreateReply(
			nil,
			err,
		))
		return
	}

	a.EmailService.AddTemplateJob(req)

	c.JSON(200, CreateReply(
		nil,
		nil,
	))
}
