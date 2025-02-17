package contract

type Email struct {
	Subject string
	Body    string
}

type EmailService interface {
	Send(sender string, contactReason string, email Email, recipient string) error
}
