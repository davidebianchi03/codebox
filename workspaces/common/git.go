package common

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"codebox.com/db"
	"codebox.com/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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

	// manage authentication
	var gitAuth transport.AuthMethod
	if strings.HasPrefix(workspace.GitRepoUrl, "http") {
		gitAuth = nil // TODO: support for authentication with token
	} else {
		gitAuth, err = ssh.NewPublicKeys("git", []byte(workspace.Owner.SshPrivateKey), "")
		if err != nil {
			return fmt.Errorf("Git authentication failure %s", err)
		}
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:             workspace.GitRepoUrl,
		InsecureSkipTLS: true,
		Progress:        &cloneLogsBuf,
		Depth:           1, // retrieve only latest commit
		SingleBranch:    true,
		Auth:            gitAuth,
		// TODO: retrieve configuration
	})
	cloning = false

	// retrieve dei log rimanenti
	newBytes := cloneLogsBuf.Bytes()[logsEndIndex:]
	if len(newBytes) > 0 {
		workspace.AppendLogs(string(newBytes))
		db.DB.Save(workspace)
	}

	if err != nil {
		if strings.HasPrefix(workspace.GitRepoUrl, "http") {
			return fmt.Errorf("Failed to clone remote repository %s", err)
		} else {
			return fmt.Errorf("Have you added the Codebox SSH public key to the remote Git server? %s", err)
		}
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
