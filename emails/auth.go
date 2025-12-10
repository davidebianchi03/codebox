package emails

import (
	"errors"

	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/db/models"
)

func SendEmailVerificationEmail(
	user models.User,
	verificationUrl string,
) error {
	html, err := RenderHtmlEmailTemplate(
		"email_verify_address",
		map[string]any{
			"verification_url": verificationUrl,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render html body")
	}

	text, err := RenderTextEmailTemplate(
		"email_verify_address",
		map[string]any{
			"verification_url": verificationUrl,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render text body")
	}

	bgtasks.BgTasksEnqueuer.Enqueue("send_email", work.Q{
		"subject":   "Please verify your email address",
		"recipient": user.Email,
		"htmlBody":  html,
		"textBody":  text,
	})

	return nil
}
