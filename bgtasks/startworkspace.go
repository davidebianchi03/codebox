package bgtasks

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/davidebianchi03/codebox/git"
	"github.com/davidebianchi03/codebox/runnerinterface"
	"github.com/davidebianchi03/codebox/utils/targz"
	"github.com/gocraft/work"
)

func (jobContext *Context) StartWorkspace(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	var workspace models.Workspace
	db.DB.Model(&models.Workspace{}).
		Preload("User").
		Preload("Runner").
		Preload("GitSource").
		Preload("TemplateVersion").
		First(&workspace, map[string]interface{}{
			"ID": workspaceId,
		})

	if workspace.ID <= 0 {
		return errors.New("workspace not found")
	}
	defer db.DB.Save(&workspace)

	workspace.ClearLogs()

	// if workspace config source is a git repository retrieve latest version
	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		// TODO: check if config file exists
		if workspace.GitSource != nil {
			if workspace.GitSource.Files == "" {
				tempDirPath, err := os.MkdirTemp("", fmt.Sprintf("codebox-%d", workspace.ID))
				if err != nil {
					workspace.AppendLogs(fmt.Sprintf("failed to create tmp folder, %s", err.Error()))
					workspace.Status = models.WorkspaceStatusError
					return nil
				}
				defer os.RemoveAll(tempDirPath)

				gitSourcesFile, err := workspace.GitSource.GetConfigFileAbsPath()
				if err != nil {
					workspace.AppendLogs(fmt.Sprintf("failed to retrieve configuration file path, %s", err.Error()))
					workspace.Status = models.WorkspaceStatusError
					return nil
				}

				if err = git.CloneRepo(
					workspace.GitSource.RepositoryURL,
					workspace.GitSource.RefName,
					tempDirPath,
					[]byte(workspace.User.SshPrivateKey),
					1,
				); err != nil {
					workspace.AppendLogs(fmt.Sprintf("failed to clone git repository, %s", err.Error()))
					workspace.Status = models.WorkspaceStatusError
					return nil
				}

				// create targz archive
				if err = targz.CreateArchive(tempDirPath, gitSourcesFile); err != nil {
					workspace.AppendLogs(fmt.Sprintf("failed to create targz archive, %s", err.Error()))
					workspace.Status = models.WorkspaceStatusError
					return nil
				}

				gitSource := workspace.GitSource
				gitSource.Files, _ = gitSource.GetConfigFileAbsPath()
				db.DB.Save(&gitSource)

				workspace.AppendLogs("the git repository has been cloned")
			}
		} else {
			workspace.AppendLogs("git source is nil")
			workspace.Status = models.WorkspaceStatusError
			return nil
		}
	}

	if workspace.Runner == nil {
		workspace.AppendLogs("runner does not exist")
		workspace.Status = models.WorkspaceStatusError
		return nil
	}

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}

	if err := ri.StartWorkspace(&workspace); err != nil {
		workspace.AppendLogs(fmt.Sprintf("failed to start workspace, %s", err.Error()))
		workspace.Status = models.WorkspaceStatusError
		return errors.New("failed to start workspace")
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

		if details.Status == models.WorkspaceStatusStarting {
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

	// map container
	for _, c := range details.Containers {
		containerUserId, err := strconv.Atoi(c.ContainerUserID)
		if err != nil {
			containerUserId = 0
		}

		workspaceContainer := models.WorkspaceContainer{
			Workspace:         workspace,
			ContainerID:       c.ID,
			ContainerName:     c.Name,
			ContainerImage:    c.Image,
			ContainerUserID:   uint(containerUserId),
			ContainerUserName: c.ContainerUserName,
			WorkspacePath:     c.WorkspacePath,
		}

		db.DB.Create(&workspaceContainer)

		// map ports
		for _, p := range c.ExposedPorts {
			port := models.WorkspaceContainerPort{
				Container:   workspaceContainer,
				ServiceName: p.ServiceName,
				PortNumber:  uint(p.PortNumber),
				Public:      p.Public,
			}

			db.DB.Create(&port)
		}

		// ping agent
		if ri.PingAgent(&workspaceContainer) {
			now := time.Now()
			workspaceContainer.AgentLastContact = &now
			db.DB.Save(&workspaceContainer)
		}
	}

	workspace.Status = details.Status
	return nil
}
