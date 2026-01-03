package emails

import (
	"errors"

	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
SendUserWaitingForApprovalEmail sends an email to notify all the superusers
that there is a user that is waiting for approval
*/
func SendUserWaitingForApprovalEmail(user models.User) error {
	html, err := RenderHtmlEmailTemplate(
		"email_user_waiting_for_approval.html",
		map[string]any{
			"user_email": user.Email,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render html body for 'email_user_waiting_for_approval'")
	}

	text, err := RenderTextEmailTemplate(
		"email_user_waiting_for_approval.txt",
		map[string]any{
			"user_email": user.Email,
		},
	)

	if err != nil {
		// TODO: log error
		return errors.New("failed to render text body for 'email_user_waiting_for_approval'")
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
		"subject":   "User Waiting For Approval",
		"recipient": recipientsCsvString,
		"htmlBody":  html,
		"textBody":  text,
	})

	return nil
}
