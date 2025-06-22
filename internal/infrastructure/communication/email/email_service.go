package email

import (
	"bytes"
	"html/template"
	"net/smtp"

	"github.com/BargheNo/Backend/bootstrap"
)

type EmailService struct {
	senderAccount *bootstrap.EmailAccount
	templatesPath *bootstrap.EmailTemplates
}

func NewEmailService(senderAccount *bootstrap.EmailAccount, templatesPath *bootstrap.EmailTemplates) *EmailService {
	return &EmailService{
		senderAccount: senderAccount,
		templatesPath: templatesPath,
	}
}

func (emailService *EmailService) SendEmail(toEmail string, subject string, templateFile string, data interface{}) error {
	from := emailService.senderAccount.EmailFrom
	password := emailService.senderAccount.EmailPassword
	smtpHost := emailService.senderAccount.SMTPHost
	smtpPort := emailService.senderAccount.SMTPPort

	tmpl, err := template.ParseFiles(emailService.templatesPath.Path + templateFile)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	body.Write([]byte("To: " + toEmail + "\r\n"))
	body.Write([]byte("Subject: " + subject + "\r\n"))
	body.Write([]byte("MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"))
	err = tmpl.Execute(&body, data)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, body.Bytes())
	if err != nil {
		return err
	}

	return nil
}
