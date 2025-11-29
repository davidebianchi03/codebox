package bgtasks

import (
	"errors"
	"fmt"
	"time"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

/*
Stop a running workspace, this is the background task
*/
func (jobContext *Context) StopWorkspaceTask(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	workspace, err := models.RetrieveWorkspaceById(uint(workspaceId))
	if err != nil {
		return nil
	}

	if workspace == nil {
		return nil
	}
	defer dbconn.DB.Save(&workspace)

	StopWorkspace(workspace, false)
	return nil
}

/*
Stop a running workspace,
this is a separate function so it can be called from multiple places
*/
func StopWorkspace(workspace *models.Workspace, skipErrors bool) error {
	defer dbconn.DB.Save(&workspace)

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}

	stopping := true
	if err := ri.StopWorkpace(workspace); err != nil {
		stopping = false
		if !skipErrors {
			workspace.AppendLogs(fmt.Sprintf("failed to stop workspace, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			return errors.New("failed to stop workspace")
		}
	}

	// fetch workspace details and logs
	logsIndex := 0
	for stopping {
		details, err := ri.GetDetails(workspace)
		if err != nil {
			if !skipErrors {
				workspace.AppendLogs(fmt.Sprintf("failed to fetch workspace details, %s", err.Error()))
				workspace.Status = models.WorkspaceStatusError
				return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
			} else {
				break
			}
		}

		if details.Status == models.WorkspaceStatusStopping {
			stopping = true
		} else {
			stopping = false
		}

		logs, err := ri.GetLogs(workspace)
		if err == nil {
			if len(logs) > logsIndex {
				logs = logs[logsIndex:]
				workspace.AppendLogs(logs)
				logsIndex += len(logs)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	details, err := ri.GetDetails(workspace)
	if err != nil {
		if !skipErrors {
			workspace.AppendLogs(fmt.Sprintf("failed to fetch workspace details, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
		}
	}

	workspace.Status = details.Status

	var containers []models.WorkspaceContainer
	dbconn.DB.Find(&containers, map[string]interface{}{
		"workspace_id": workspace.ID,
	})
	for _, container := range containers {
		dbconn.DB.Unscoped().Delete(&[]models.WorkspaceContainerPort{}, map[string]interface{}{
			"container_id": container.ID,
		})
		dbconn.DB.Unscoped().Delete(&container)
	}

	return nil
}
