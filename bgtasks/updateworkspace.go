package bgtasks

import (
	"errors"
	"fmt"
	"os"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/davidebianchi03/codebox/git"
	"github.com/davidebianchi03/codebox/utils/targz"
	"github.com/gocraft/work"
)

func (jobContext *Context) UpdateWorkspaceConfigFiles(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	var workspace *models.Workspace
	result := db.DB.Preload("Runner").
		Preload("User").
		Preload("Runner").
		Preload("GitSource").
		Preload("TemplateVersion").
		First(&workspace, map[string]interface{}{"ID": workspaceId})

	if result.Error != nil {
		return fmt.Errorf("failed to retrieve workspace from db %s", result.Error)
	}

	if workspace == nil {
		return errors.New("workspace not found")
	}

	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		tempDirPath, err := os.MkdirTemp("", fmt.Sprintf("codebox-%d", workspace.ID))
		if err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to create tmp folder, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			db.DB.Save(&workspace)
			return nil
		}
		defer os.RemoveAll(tempDirPath)

		gitSourcesFile, err := workspace.GitSource.GetConfigFileAbsPath()
		if err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to retrieve configuration file path, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			db.DB.Save(&workspace)
			return nil
		}
		os.RemoveAll(gitSourcesFile)

		if err = git.CloneRepo(
			workspace.GitSource.RepositoryURL,
			workspace.GitSource.RefName,
			tempDirPath,
			[]byte(workspace.User.SshPrivateKey),
			1,
		); err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to clone git repository, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			db.DB.Save(&workspace)
			return nil
		}

		// create targz archive
		if err = targz.CreateArchive(tempDirPath, gitSourcesFile); err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to create targz archive, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			db.DB.Save(&workspace)
			return nil
		}

		gitSource := workspace.GitSource
		gitSource.Files, _ = gitSource.GetConfigFileAbsPath()
		db.DB.Save(&gitSource)
	} else {
		panic("not implemented")
	}
	db.DB.Save(&workspace)
	workspace.AppendLogs("Config files have been updated")

	return jobContext.StartWorkspace(job)
}
