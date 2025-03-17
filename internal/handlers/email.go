package handler

import (
	contract "mailer-go/internal/contracts"

	"github.com/gin-gonic/gin"
)

func (a *App) Health(c *gin.Context) {
	Respond(c, "OK", nil)
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

	a.EmailService.AddJob(req)

	Respond(c, nil, nil)
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

	a.EmailService.AddTemplateJob(req)

	Respond(c, nil, nil)
}
