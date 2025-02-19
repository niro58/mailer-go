package contract

type Email struct {
	Subject string
	Body    string
}
type EmailService interface {
	Send(senderKey string, recipient string, email Email) error
}
