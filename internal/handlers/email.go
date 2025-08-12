package handler

import (
	contract "mailer-go/internal/contracts"

	"github.com/gin-gonic/gin"
)

func (a *App) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
func (a *App) Send(c *gin.Context) {
	var req contract.Email
	if err := c.ShouldBind(&req); err != nil {
		Respond(c, nil, err)
		return
	}

	if req.ContentType == "" {
		req.ContentType = "text/plain"
	}

	err := a.EmailService.AddJob(req)

	Respond(c, nil, err)
}
func (a *App) Status(c *gin.Context) {
	Respond(c, a.EmailService.Count(), nil)
}
func (a *App) SendTemplate(c *gin.Context) {
	var req contract.EmailTemplate
	if err := c.ShouldBind(&req); err != nil {
		Respond(c, nil, err)
		return
	}

	err := a.EmailService.AddTemplateJob(req)

	Respond(c, nil, err)
}
