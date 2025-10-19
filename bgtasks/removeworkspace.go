package bgtasks

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

func removeWorkspace(workspace models.Workspace, skipErrors bool) error {
	if workspace.Runner != nil {
		ri := runnerinterface.RunnerInterface{
			Runner: workspace.Runner,
		}

		err := ri.RemoveWorkspace(&workspace)
		if err != nil && !skipErrors {
			workspace.AppendLogs(fmt.Sprintf("failed to remove workspace, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			dbconn.DB.Save(&workspace)
			return errors.New("failed to remove workspace")
		}

		if err == nil {
			// fetch workspace details and logs
			starting := true
			logsIndex := 0
			for starting {
				details, err := ri.GetDetails(&workspace)
				if err != nil {
					workspace.AppendLogs(fmt.Sprintf("failed to fetch workspace details, %s", err.Error()))
					if skipErrors {
						break
					} else {
						workspace.Status = models.WorkspaceStatusError
						dbconn.DB.Save(&workspace)
						return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
					}
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
		}

		details, err := ri.GetDetails(&workspace)
		if err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to fetch workspace details, %s", err.Error()))
			if !skipErrors {
				workspace.Status = models.WorkspaceStatusError
				dbconn.DB.Save(&workspace)
				return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
			}
			workspace.Status = models.WorkspaceStatusDeleting
		} else {
			workspace.Status = details.Status
		}
	}

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

	// remove configuration files if the source is git
	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		if workspace.GitSource != nil {
			if workspace.GitSource.Sources != nil {
				os.RemoveAll(workspace.GitSource.Sources.GetAbsolutePath())
			}
			dbconn.DB.Unscoped().Delete(&workspace.GitSource)
		}
	}
	workspace.ClearLogs()

	dbconn.DB.Unscoped().Delete(&workspace)
	return nil
}

func (jobContext *Context) DeleteWorkspaceTask(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")
	skipErrors := job.ArgBool("skip_errors")

	var workspace models.Workspace
	result := dbconn.DB.Preload("Runner").Preload("GitSource").First(&workspace, map[string]interface{}{"ID": workspaceId})
	if result.Error != nil {
		return fmt.Errorf("failed to retrieve workspace from db %s", result.Error)
	}

	if workspace.ID <= 0 {
		return errors.New("workspace not found")
	}

	return removeWorkspace(workspace, skipErrors)
}
