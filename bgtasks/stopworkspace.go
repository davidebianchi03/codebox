package bgtasks

import (
	"errors"
	"fmt"
	"time"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/notifications"
	"gitlab.com/codebox4073715/codebox/logging"
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

	notifications.SendWorkspaceStopNotification(*workspace)
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

	// create a backup of the exposed ports before stopping the workspace
	containers, err := models.ListWorkspaceContainersByWorkspace(*workspace)
	if err != nil {
		if !skipErrors {
			workspace.AppendLogs(fmt.Sprintf("failed to list workspace containers, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			return fmt.Errorf("failed to list workspace containers, %s", err.Error())
		}
	}

	for _, container := range containers {
		ports, err := models.ListContainerPortsByWorkspaceContainer(container)
		if err != nil {
			if !skipErrors {
				logging.Error(
					"failed to list container ports for container %s, %s",
					container.ContainerName,
					err.Error(),
				)
				workspace.AppendLogs("internal server error")
				workspace.Status = models.WorkspaceStatusError
				return fmt.Errorf("failed to list container ports for container %s, %s", container.ContainerName, err.Error())
			}
		}

		for _, port := range ports {
			_, err := models.CreateContainerPortBackup(
				container,
				port.ServiceName,
				port.PortNumber,
				port.Public,
			)
			if err != nil {
				if !skipErrors {
					logging.Error(
						"failed to create container port backup for container %s, %s",
						container.ContainerName,
						err.Error(),
					)
					workspace.AppendLogs("internal server error")
					workspace.Status = models.WorkspaceStatusError
					return fmt.Errorf("failed to create container port backup for container %s, %s", container.ContainerName, err.Error())
				}
			}
		}
	}

	stopping := true
	if err := ri.StopWorkspace(workspace); err != nil {
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
		details, err := ri.GetWorkspaceDetails(workspace)
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

		logs, err := ri.GetWorkspaceLogs(workspace)
		if err == nil {
			if len(logs) > logsIndex {
				logs = logs[logsIndex:]
				workspace.AppendLogs(logs)
				logsIndex += len(logs)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	details, err := ri.GetWorkspaceDetails(workspace)
	if err != nil {
		if !skipErrors {
			workspace.AppendLogs(fmt.Sprintf("failed to fetch workspace details, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
		}
	}

	workspace.Status = details.Status

	for _, container := range containers {
		if err := models.DeleteContainer(container); err != nil {
			logging.Error(
				"failed to delete container %s of workspace %s, %s",
				container.ContainerName,
				workspace.Name,
				err.Error(),
			)
			return err
		}
	}

	notifications.SendWorkspaceStoppedNotification(*workspace)

	return nil
}
