package bgtasks

import (
	"time"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
Restart a running workspace, this is the background task
First stops the workspace, then starts it again
*/
func (jobContext *Context) RestartWorkspaceTask(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	workspace, err := models.RetrieveWorkspaceById(uint(workspaceId))
	if err != nil {
		return nil
	}

	if workspace == nil {
		return nil
	}
	defer dbconn.DB.Save(&workspace)

	// First stop the workspace, skip errors during stop to proceed with restart
	workspace.AppendLogs("Stopping workspace...")
	if err := StopWorkspace(workspace, false); err != nil {
		return nil
	}

	workspace.Status = models.WorkspaceStatusStarting
	dbconn.DB.Save(&workspace)

	// Give a moment for the workspace to fully stop before restarting
	time.Sleep(1 * time.Second)

	// Now start the workspace again
	workspace.Status = models.WorkspaceStatusStarting
	workspace.AppendLogs("Starting workspace...")

	if err := StartWorkspace(workspace); err != nil {
		return nil
	}

	return nil
}
