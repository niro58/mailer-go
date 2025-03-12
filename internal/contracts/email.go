package contract

type EmailHeaders struct {
	SenderKey  string   `json:"senderKey"`
	Recipients []string `json:"recipients"`
	Bcc        []string `json:"bcc"`
}

type Email struct {
	EmailHeaders

	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailTemplate struct {
	EmailHeaders
	TemplateKey string            `json:"string"`
	Variables   map[string]string `json:"variables"`
}

type EmailService interface {
	StartPool()
	AddJob(email Email) error
	AddTemplateJob(template EmailTemplate) error
	Count() int
}
