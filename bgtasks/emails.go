package bgtasks

import (
	"log"

	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/config"
	gomail "gopkg.in/mail.v2"
)

func (jobContext *Context) SendEmailTask(job *work.Job) error {
	if config.Environment.EmailSMTPHost == "" || config.Environment.EmailSMTPPort == 0 ||
		config.Environment.EmailSMTPUser == "" || config.Environment.EmailSMTPPassword == "" {
		log.Println("Email sender is not configured")
		return nil
	}

	subject := job.ArgString("subject")
	recipient := job.ArgString("recipient")
	htmlBody := job.ArgString("htmlBody")
	textBody := job.ArgString("textBody")

	message := gomail.NewMessage()
	message.SetHeader("From", config.Environment.EmailSMTPUser)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", subject)

	message.SetBody("text/html", htmlBody)
	message.AddAlternative("text/plain", textBody)

	dialer := gomail.NewDialer(
		config.Environment.EmailSMTPHost,
		config.Environment.EmailSMTPPort,
		config.Environment.EmailSMTPUser,
		config.Environment.EmailSMTPPassword,
	)

	if err := dialer.DialAndSend(message); err != nil {
		// TODO: log error
		log.Println("Failed to send email:", err)
		return nil
	}

	return nil
}
