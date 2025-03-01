package bgtasks

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/davidebianchi03/codebox/git"
	runnerinterface "github.com/davidebianchi03/codebox/runner-interface"
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

				if err = git.CloneRepo(workspace.GitSource.RepositoryURL, tempDirPath, []byte(workspace.User.SshPrivateKey), 1); err != nil {
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
		return errors.New("failed to start workspace")
	}

	// fetch workspace details and logs
	starting := true
	logsIndex := 0
	for starting {
		details, err := ri.GetDetails(&workspace)
		if err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to fetch workspace details, %s", err.Error()))
			return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
		}

		if details.Status == models.WorkspaceStatusStarting {
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
		return fmt.Errorf("failed to fetch workspace details, %s", err.Error())
	}

	workspace.Status = details.Status
	return nil
}

// 	workspace.ClearLogs()

// 	if workspace.Type == db.WorkspaceTypeDevcontainer {
// 		workspaceInterface := devcontainer.DevcontainerWorkspace{}
// 		workspaceInterface.Workspace = workspace
// 		workspaceInterface.StartWorkspace()
// 	} else {
// 		return fmt.Errorf("%s: unsupported workspace type", workspace.Type)
// 	}

// 	return nil
// }

// func (ctx *WorkspaceTaskContext) StopWorkspace(job *work.Job) error {
// 	workspaceId := job.ArgInt64("workspace_id")

// 	var workspace *db.Workspace
// 	result := db.DB.Where("ID=?", workspaceId).Preload("Owner").First(&workspace)
// 	if result.Error != nil {
// 		return fmt.Errorf("failed to retrieve workspace from db %s", result.Error)
// 	}

// 	workspace.ClearLogs()

// 	if workspace.Type == db.WorkspaceTypeDevcontainer {
// 		workspace.AppendLogs("Stopping workspace...")
// 		workspace.Status = db.WorkspaceStatusStopping
// 		db.DB.Save(&workspace)
// 		workspaceInterface := devcontainer.DevcontainerWorkspace{}
// 		workspaceInterface.Workspace = workspace
// 		workspaceInterface.StopWorkspace()
// 	} else {
// 		return fmt.Errorf("%s: unsupported workspace type", workspace.Type)
// 	}

// 	return nil
// }

// func (ctx *WorkspaceTaskContext) RestartWorkspace(job *work.Job) error {
// 	return nil
// }

// func (ctx *WorkspaceTaskContext) DeleteWorkspace(job *work.Job) error {
// 	workspaceId := job.ArgInt64("workspace_id")

// 	var workspace *db.Workspace
// 	result := db.DB.Where("ID=?", workspaceId).Preload("Owner").First(&workspace)
// 	if result.Error != nil {
// 		return fmt.Errorf("failed to retrieve workspace from db %s", result.Error)
// 	}

// 	workspace.ClearLogs()

// 	if workspace.Type == db.WorkspaceTypeDevcontainer {
// 		workspace.AppendLogs("Deleting workspace...")
// 		workspace.Status = db.WorkspaceStatusDeleting
// 		db.DB.Save(&workspace)
// 		workspaceInterface := devcontainer.DevcontainerWorkspace{}
// 		workspaceInterface.Workspace = workspace
// 		workspaceInterface.StopWorkspace()
// 		workspace.Status = db.WorkspaceStatusDeleting
// 		db.DB.Save(&workspace)
// 		workspaceInterface.DeleteWorkspace()
// 	} else {
// 		return fmt.Errorf("%s: unsupported workspace type", workspace.Type)
// 	}

// 	return nil
// }
