package bgtasks

import (
	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
Background task that deletes a user,
this task deletes the user and all his workspaces
*/
func (jobContext *Context) DeleteUserTask(job *work.Job) error {
	// TODO: send emails to admin for errors
	userEmail := job.ArgString("user_email")

	user, err := models.RetrieveUserByEmail(userEmail)
	if err != nil {
		// TODO: log error
		return err
	}

	// delete all the workspaces
	workspaces, err := models.ListUserWorkspaces(*user)
	if err != nil {
		// TODO: log error
		return err
	}

	for _, w := range workspaces {
		err := removeWorkspace(w, true)
		if err != nil {
			// TODO: log error
		}
	}

	// delete the user
	if err := models.DeleteUser(user); err != nil {
		// TODO: log error
		return err
	}

	return nil
}
