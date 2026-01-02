package emails

import (
	"errors"

	"gitlab.com/codebox4073715/codebox/bgtasks"
)

/*
SendTestEmail renders and sends an email to test
that the email sending service is working.
The message is not sent with a bg task like other email.
In this way it's possibile get more details about failures.
*/
func SendTestEmail(
	emailAddress string,
) (err error) {
	html, err := RenderHtmlEmailTemplate(
		"test_email.html",
		map[string]any{},
	)
	if err != nil {
		// TODO: log error
		return errors.New("failed to render html body for 'test_email'")
	}

	text, err := RenderTextEmailTemplate(
		"test_email.txt",
		map[string]any{},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render text body for 'test_email'")
	}

	if err := bgtasks.SendEmailMessage(
		emailAddress,
		"Test Email",
		html,
		text,
	); err != nil {
		return err
	}

	return nil
}
