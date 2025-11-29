package bgtasks

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gocraft/work"
	"github.com/google/uuid"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/git"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
	"gitlab.com/codebox4073715/codebox/utils/targz"
)

func (jobContext *Context) StartWorkspaceTask(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	workspace, err := models.RetrieveWorkspaceById(uint(workspaceId))
	if err != nil {
		return nil
	}

	if workspace == nil {
		return nil
	}
	defer dbconn.DB.Save(&workspace)

	// if workspace config source is a git repository retrieve latest version
	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		if workspace.GitSource != nil {
			if workspace.GitSource.Sources == nil {
				gitSources := models.File{
					Filepath: path.Join("git-sources", fmt.Sprintf("%s.tar.gz", uuid.New().String())),
				}
				dbconn.DB.Save(&gitSources)

				workspace.GitSource.SourcesID = &gitSources.ID
				workspace.GitSource.Sources = &gitSources
				dbconn.DB.Save(&workspace.GitSource)
				dbconn.DB.Save(&workspace)
			}

			// check if config files exists, clone them if not exist
			if !workspace.GitSource.Sources.Exists() {
				tempDirPath, err := os.MkdirTemp("", fmt.Sprintf("codebox-%d", workspace.ID))
				if err != nil {
					workspace.AppendLogs(fmt.Sprintf("failed to create tmp folder, %s", err.Error()))
					workspace.Status = models.WorkspaceStatusError
					return nil
				}
				defer os.RemoveAll(tempDirPath)

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
				tgm := targz.TarGZManager{Filepath: workspace.GitSource.Sources.GetAbsolutePath()}
				if err = tgm.CompressFolder(tempDirPath); err != nil {
					workspace.AppendLogs(fmt.Sprintf("failed to create targz archive, %s", err.Error()))
					workspace.Status = models.WorkspaceStatusError
					return nil
				}

				workspace.AppendLogs("the git repository has been cloned")
			}
		} else {
			workspace.AppendLogs("git source is nil")
			workspace.Status = models.WorkspaceStatusError
			return nil
		}
	} else {
		// check if config files exist
		if workspace.TemplateVersion.Sources == nil {
			workspace.AppendLogs("Template version has no sources")
			workspace.Status = models.WorkspaceStatusError
			return nil
		}

		if !workspace.TemplateVersion.Sources.Exists() {
			workspace.AppendLogs("Template version has no sources")
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

	if err := ri.StartWorkspace(workspace); err != nil {
		workspace.AppendLogs(fmt.Sprintf("failed to start workspace, %s", err.Error()))
		workspace.Status = models.WorkspaceStatusError
		return errors.New("failed to start workspace")
	}

	// fetch workspace details and logs
	starting := true
	logsIndex := 0
	for starting {
		details, err := ri.GetDetails(workspace)
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
			Workspace:         *workspace,
			ContainerID:       c.ID,
			ContainerName:     c.Name,
			ContainerImage:    c.Image,
			ContainerUserID:   uint(containerUserId),
			ContainerUserName: c.ContainerUserName,
			WorkspacePath:     c.WorkspacePath,
		}

		dbconn.DB.Create(&workspaceContainer)

		// map ports
		for _, p := range c.ExposedPorts {
			port := models.WorkspaceContainerPort{
				Container:   workspaceContainer,
				ServiceName: p.ServiceName,
				PortNumber:  uint(p.PortNumber),
				Public:      p.Public,
			}

			dbconn.DB.Create(&port)
		}

		// ping agent
		if ri.PingAgent(&workspaceContainer) {
			now := time.Now()
			workspaceContainer.AgentLastContact = &now
			dbconn.DB.Save(&workspaceContainer)
		}
	}

	workspace.Status = details.Status
	return nil
}
