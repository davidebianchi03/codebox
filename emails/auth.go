package emails

import (
	"errors"

	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
SendUserAlreadyExistsEmail renders and sends an email to notify that
the email address is already associated with an existing account.
*/
func SendUserAlreadyExistsEmail(
	emailAddress string,
) error {
	html, err := RenderHtmlEmailTemplate(
		"email_notify_existing.html",
		map[string]any{},
	)
	if err != nil {
		// TODO: log error
		return errors.New("failed to render html body for 'email_notify_existing'")
	}

	text, err := RenderTextEmailTemplate(
		"email_notify_existing.txt",
		map[string]any{},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render text body for 'email_notify_existing'")
	}

	bgtasks.BgTasksEnqueuer.Enqueue("send_email", work.Q{
		"subject":   "Sign-up Attempt Notice",
		"recipient": emailAddress,
		"htmlBody":  html,
		"textBody":  text,
	})

	return nil
}

/*
SendEmailVerificationEmail renders and sends an email verification email
to the given user with the provided verification URL.
*/
func SendEmailVerificationEmail(
	user models.User,
	verificationUrl string,
) error {
	html, err := RenderHtmlEmailTemplate(
		"email_verify_address.html",
		map[string]any{
			"verification_url": verificationUrl,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render html body for 'email_verify_address'")
	}

	text, err := RenderTextEmailTemplate(
		"email_verify_address.txt",
		map[string]any{
			"verification_url": verificationUrl,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render text body for 'email_verify_address'")
	}

	bgtasks.BgTasksEnqueuer.Enqueue("send_email", work.Q{
		"subject":   "Please verify your email address",
		"recipient": user.Email,
		"htmlBody":  html,
		"textBody":  text,
	})

	return nil
}
