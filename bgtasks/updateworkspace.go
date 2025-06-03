package bgtasks

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/gocraft/work"
	"github.com/google/uuid"
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
		Preload("GitSource.Sources").
		Preload("TemplateVersion").
		Preload("TemplateVersion.Template").
		First(&workspace, map[string]interface{}{"ID": workspaceId})

	if result.Error != nil {
		return fmt.Errorf("failed to retrieve workspace from db %s", result.Error)
	}

	if workspace == nil {
		return errors.New("workspace not found")
	}

	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		if workspace.GitSource.Sources == nil {
			gitSources := models.File{
				Filepath: path.Join("git-sources", fmt.Sprintf("%s.tar.gz", uuid.New().String())),
			}
			dbconn.DB.Save(&gitSources)

			workspace.GitSource.SourcesID = gitSources.ID
			workspace.GitSource.Sources = &gitSources
			dbconn.DB.Save(&workspace)
		}

		// remove existsing files and clone repository again
		os.RemoveAll(workspace.GitSource.Sources.GetAbsolutePath())

		tempDirPath, err := os.MkdirTemp("", fmt.Sprintf("codebox-%d", workspace.ID))
		if err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to create tmp folder, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			dbconn.DB.Save(&workspace)
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
			dbconn.DB.Save(&workspace)
			return nil
		}

		// create targz archive
		tgm := targz.TarGZManager{Filepath: workspace.GitSource.Sources.GetAbsolutePath()}
		if err = tgm.CompressFolder(tempDirPath); err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to create targz archive, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			dbconn.DB.Save(&workspace)
			return nil
		}
	} else {
		latestVersion, err := models.RetrieveLatestTemplateVersionByTemplate(*workspace.TemplateVersion.Template)
		if err != nil {
			workspace.AppendLogs(fmt.Sprintf("failed to create targz archive, %s", err.Error()))
			workspace.Status = models.WorkspaceStatusError
			dbconn.DB.Save(&workspace)
			return nil
		}

		workspace.TemplateVersionID = &latestVersion.ID
		workspace.TemplateVersion = latestVersion
		dbconn.DB.Save(&workspace)
	}

	dbconn.DB.Save(&workspace)
	workspace.AppendLogs("Config files have been updated")

	return jobContext.StartWorkspace(job)
}
