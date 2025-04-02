package bgtasks

import (
	"errors"
	"fmt"
	"time"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/davidebianchi03/codebox/runnerinterface"
	"github.com/gocraft/work"
)

func (jobContext *Context) DeleteWorkspace(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	var workspace models.Workspace
	result := db.DB.Preload("Runner").First(&workspace, map[string]interface{}{"ID": workspaceId})
	if result.Error != nil {
		return fmt.Errorf("failed to retrieve workspace from db %s", result.Error)
	}

	if workspace.ID <= 0 {
		return errors.New("workspace not found")
	}
	defer db.DB.Save(&workspace)
	workspace.ClearLogs()

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}

	if err := ri.RemoveWorkspace(&workspace); err != nil {
		workspace.AppendLogs(fmt.Sprintf("failed to remove workspace, %s", err.Error()))
		workspace.Status = models.WorkspaceStatusError
		return errors.New("failed to remove workspace")
	}

	// fetch workspace details and logs
	starting := true
	logsIndex := 0
	for starting {
		details, err := ri.GetDetails(&workspace)
		if err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to fetch workspace details, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
		}

		if details.Status == models.WorkspaceStatusStopping {
			starting = true
		} else {
			starting = false
		}

		logs, err := ri.GetLogs(&workspace)
		if err == nil {
			if len(logs) > logsIndex {
				logs = logs[logsIndex:]
				workspace.AppendLogs(logs)
				logsIndex += len(logs)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	details, err := ri.GetDetails(&workspace)
	if err != nil {
		workspace.AppendLogs(fmt.Sprintf("failed to fetch workspace details, %s", err.Error()))
		workspace.Status = models.WorkspaceStatusError
		return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
	}

	workspace.Status = details.Status

	var containers []models.WorkspaceContainer
	db.DB.Find(&containers, map[string]interface{}{
		"workspace_id": workspace.ID,
	})
	for _, container := range containers {
		db.DB.Unscoped().Delete(&[]models.WorkspaceContainerPort{}, map[string]interface{}{
			"container_id": container.ID,
		})
		db.DB.Unscoped().Delete(&container)
	}
	db.DB.Delete(&workspace)
	return nil
}
