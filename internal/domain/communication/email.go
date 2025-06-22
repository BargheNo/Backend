package communication

type EmailService interface {
	SendEmail(toEmail string, subject string, templateFile string, data interface{}) error
}
