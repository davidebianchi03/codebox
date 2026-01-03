package emails

import (
	"errors"

	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
SendRatelimitExceededMultipleTimesEmail sends an email to notify
all admins that the rate limit on an endpoint has been
exceeded multiple times by an IP address.
*/
func SendRatelimitExceededMultipleTimesEmail(
	ipAddress string,
	endpoint string,
) error {
	html, err := RenderHtmlEmailTemplate(
		"email_ratelimit_notification.html",
		map[string]any{
			"ip_address": ipAddress,
			"endpoint":   endpoint,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render html body for 'email_ratelimit_notification'")
	}

	text, err := RenderTextEmailTemplate(
		"email_ratelimit_notification.txt",
		map[string]any{
			"ip_address": ipAddress,
			"endpoint":   endpoint,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render text body for 'email_ratelimit_notification'")
	}

	adminUsers, err := models.ListSuperUsers()
	if err != nil {
		// TODO: log error
		return errors.New("failed to list superusers")
	}

	recipientsCsvString := ""
	for _, admin := range adminUsers {
		if recipientsCsvString != "" {
			recipientsCsvString += ","
		}

		recipientsCsvString += admin.Email
	}

	bgtasks.BgTasksEnqueuer.Enqueue("send_email", work.Q{
		"subject":   "Ratelimit Excedeed Multiple Times",
		"recipient": recipientsCsvString,
		"htmlBody":  html,
		"textBody":  text,
	})

	return nil
}
