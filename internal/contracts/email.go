package contract

type EmailHeaders struct {
	SenderKey  string   `json:"senderKey" form:"senderKey" binding:"required"`
	Recipients []string `json:"recipients" form:"recipients" binding:"required"`
	Bcc        []string `json:"bcc" form:"bcc"`
}

type Email struct {
	EmailHeaders
	ContentType string `json:"contentType" form:"contentType"`
	Subject     string `json:"subject" form:"subject" binding:"required"`
	Body        string `json:"body" form:"body" binding:"required"`
}

type EmailTemplate struct {
	EmailHeaders
	TemplateKey string            `json:"templateKey"`
	Variables   map[string]string `json:"variables"`
}

type EmailService interface {
	StartPool()
	AddJob(email Email) error
	AddTemplateJob(template EmailTemplate) error
	Count() int
	Wait()
}
