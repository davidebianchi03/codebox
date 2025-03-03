package bgtasks

import (
	"errors"
	"fmt"
	"time"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	runnerinterface "github.com/davidebianchi03/codebox/runner-interface"
	"github.com/gocraft/work"
)

func (jobContext *Context) StopWorkspace(job *work.Job) error {
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

	if err := ri.StopWorkpace(&workspace); err != nil {
		workspace.AppendLogs(fmt.Sprintf("failed to stop workspace, %s", err.Error()))
		workspace.Status = models.WorkspaceStatusError
		return errors.New("failed to stop workspace")
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
			logs = logs[logsIndex:]
			workspace.AppendLogs(logs)
			logsIndex += len(logs)
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
	return nil
}
