package bgtasks

import (
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

/*
Update workspace configuration files,
this task fetches the latest configuration files from the git repository
or updates the template version if the config source is template
then restarts the workspace
*/
func (jobContext *Context) UpdateWorkspaceConfigFilesTask(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	workspace, err := models.RetrieveWorkspaceById(uint(workspaceId))
	if err != nil {
		return nil
	}

	if workspace == nil {
		return nil
	}

	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		if workspace.GitSource.Sources == nil {
			gitSources := models.File{
				Filepath: path.Join("git-sources", fmt.Sprintf("%s.tar.gz", uuid.New().String())),
			}
			dbconn.DB.Save(&gitSources)

			workspace.GitSource.SourcesID = &gitSources.ID
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

	jobContext.StartWorkspaceTask(job)
	return nil
}
