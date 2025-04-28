package bgtasks

import (
	"errors"
	"fmt"
	"os"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/git"
	"gitlab.com/codebox4073715/codebox/utils/targz"
)

func (jobContext *Context) UpdateWorkspaceConfigFiles(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	var workspace *models.Workspace
	result := dbconn.DB.Preload("Runner").
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
			dbconn.DB.Save(&workspace)
			return nil
		}
		defer os.RemoveAll(tempDirPath)

		gitSourcesFile, err := workspace.GitSource.GetConfigFileAbsPath()
		if err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to retrieve configuration file path, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			dbconn.DB.Save(&workspace)
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
			dbconn.DB.Save(&workspace)
			return nil
		}

		// create targz archive
		tgm := targz.TarGZManager{Filepath: gitSourcesFile}
		if err = tgm.CompressFolder(tempDirPath); err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to create targz archive, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			dbconn.DB.Save(&workspace)
			return nil
		}

		gitSource := workspace.GitSource
		gitSource.Files, _ = gitSource.GetConfigFileAbsPath()
		dbconn.DB.Save(&gitSource)
	} else {
		panic("not implemented")
	}
	dbconn.DB.Save(&workspace)
	workspace.AppendLogs("Config files have been updated")

	return jobContext.StartWorkspace(job)
}
