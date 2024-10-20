package common

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"codebox.com/db"
	"codebox.com/utils"
	"github.com/go-git/go-git/v5"
)

func RetrieveWorkspaceConfigFilesFromGitRepo(workspace *db.Workspace) error {
	dir, err := os.MkdirTemp("", fmt.Sprintf("workspace_%s_%d", workspace.Name, workspace.ID))
	if err != nil {
		return fmt.Errorf("failed to create a temporary directory for cloning the Git repository, preventing retrieval of workspace configuration files")
	}
	defer os.RemoveAll(dir)

	// clone git repository
	var cloneLogsBuf bytes.Buffer
	cloning := true
	logsEndIndex := 0
	go func() {
		for cloning {
			newBytes := cloneLogsBuf.Bytes()[logsEndIndex:]
			if len(newBytes) > 0 {
				workspace.AppendLogs(string(newBytes))
				db.DB.Save(&workspace)
				logsEndIndex += len(newBytes)
			}
		}
	}()
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:             workspace.GitRepoUrl,
		InsecureSkipTLS: true,
		Progress:        &cloneLogsBuf,
		// TODO: auth with ssh keys
	})
	cloning = false

	// retrieve dei log rimanenti
	newBytes := cloneLogsBuf.Bytes()[logsEndIndex:]
	if len(newBytes) > 0 {
		workspace.AppendLogs(string(newBytes))
		db.DB.Save(workspace)
	}

	if err != nil {
		return fmt.Errorf("Failed to clone remote repository %s", err)
	}

	configurationFolderPath := filepath.Join(dir, workspace.GitRepoConfigurationFolder)
	pathInfo, err := os.Stat(configurationFolderPath)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Configuration folder path does not exists")
		} else {
			return fmt.Errorf("Unknown error occured while retrieving workspace configuration: %s", err)
		}
	}

	if !pathInfo.IsDir() {
		return fmt.Errorf("Configuration folder is not a directory")
	}

	outputFilePath, err := workspace.GetConfigFilePath()
	if err != nil {
		return err
	}

	err = utils.CreateNewTarGzArchive(configurationFolderPath, outputFilePath)
	if err != nil {
		return fmt.Errorf("cannot create targz archive: %s", err)
	}

	workspace.WorkspaceConfigurationFiles = outputFilePath
	db.DB.Save(workspace)

	return nil
}
