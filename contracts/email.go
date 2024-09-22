package contract

type Email struct {
	Subject string
	Body    string
}
type EmailService interface {
	Send(email Email) error
}