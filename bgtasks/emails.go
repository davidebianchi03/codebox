package bgtasks

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/config"
	gomail "gopkg.in/mail.v2"
)

/*
SendEmailTask is the task used to send emails in background
It takes the following arguments from the job context:
- subject: the subject of the message
- recipient: the address of the recipient, can be a csv list
- textBody: the email text content
- htmlBody: the body of the email rendered as an html document
*/
func (jobContext *Context) SendEmailTask(job *work.Job) error {
	subject := job.ArgString("subject")
	recipient := job.ArgString("recipient")
	htmlBody := job.ArgString("htmlBody")
	textBody := job.ArgString("textBody")

	if err := SendEmailMessage(
		strings.Split(recipient, ","),
		subject,
		htmlBody,
		textBody,
	); err != nil {
		// TODO: log error
		log.Println(err)
	}

	return nil
}

/*
This function send emails, this function can be invoked with sendemailtask to send
messages using a background task, it can be also directly invokne
*/
func SendEmailMessage(
	recipients []string,
	subject,
	htmlBody,
	textBody string,
) error {
	if config.Environment.EmailSMTPHost == "" || config.Environment.EmailSMTPPort == 0 ||
		config.Environment.EmailSMTPUser == "" || config.Environment.EmailSMTPPassword == "" {
		return fmt.Errorf("Email sender is not configured")
	}

	message := gomail.NewMessage()
	message.SetHeader("From", config.Environment.EmailSMTPUser)
	message.SetHeader("To", recipients...)
	message.SetHeader("Subject", subject)

	message.SetBody("text/plain", textBody)
	message.AddAlternative("text/html", htmlBody)

	dialer := gomail.NewDialer(
		config.Environment.EmailSMTPHost,
		config.Environment.EmailSMTPPort,
		config.Environment.EmailSMTPUser,
		config.Environment.EmailSMTPPassword,
	)

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("Failed to send email: %s", err)
	}

	return nil
}
