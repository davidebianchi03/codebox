package emails

import (
	"errors"
	"fmt"

	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
)

/*
SendPasswordResetEmail renders and sends a password reset email
to the given email address with the provided token.
The message is sent as a background task.
*/
func SendPasswordResetEmail(
	emailAddress string,
	token string,
) (err error) {
	resetPasswordUrl := fmt.Sprintf("%s/reset-password?token=%s",
		config.Environment.ExternalUrl,
		token,
	)
	html, err := RenderHtmlEmailTemplate(
		"password_reset.html",
		map[string]any{
			"reset_url": resetPasswordUrl,
		},
	)
	if err != nil {
		// TODO: log error
		return errors.New("failed to render html body for 'password_reset'")
	}

	text, err := RenderTextEmailTemplate(
		"password_reset.txt",
		map[string]any{
			"reset_url": resetPasswordUrl,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render text body for 'password_reset'")
	}

	if err := bgtasks.SendEmailMessage(
		[]string{emailAddress},
		"Password Reset",
		html,
		text,
	); err != nil {
		return err
	}

	return nil
}

/*
SendUserNotRegisteredEmail renders and sends an email to inform
the recipient that a password reset was requested for an email
not registered in the system.
The message is sent as a background task.
*/
func SendUserNotRegisteredEmail(
	emailAddress string,
) (err error) {
	html, err := RenderHtmlEmailTemplate(
		"user_not_registered.html",
		map[string]any{},
	)
	if err != nil {
		// TODO: log error
		return errors.New("failed to render html body for 'user_not_registered'")
	}

	text, err := RenderTextEmailTemplate(
		"user_not_registered.txt",
		map[string]any{},
	)

	if err != nil {
		return errors.New("failed to render text body for 'user_not_registered'")
	}

	if err := bgtasks.SendEmailMessage(
		[]string{emailAddress},
		"Password Reset Request",
		html,
		text,
	); err != nil {
		return err
	}

	return nil
}
