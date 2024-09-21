package devcontainer

import (
	"fmt"
	"os"
	"path"

	"codebox.com/db"
	"codebox.com/utils"
	"codebox.com/workspaces/common"
)

type DevcontainerWorkspace struct {
	common.CommonWorkspace
}

func (dw *DevcontainerWorkspace) StartWorkspace() {
	// ottengo la configurazione dal repo git se non esiste il file
	if dw.Workspace.WorkspaceConfigurationFiles == "" {
		err := common.RetrieveWorkspaceConfigFilesFromGitRepo(dw.Workspace)
		if err != nil {
			dw.Workspace.Logs += err.Error() + "\n"
			dw.Workspace.Status = db.WorkspaceStatusError
			db.DB.Save(&dw.Workspace)
			return
		}
		dw.Workspace.Logs += "\n"
		db.DB.Save(&dw.Workspace)
	}

	if dw.Workspace.WorkspaceConfigurationFiles == "" {
		dw.Workspace.Logs += "missing workspace configuration files, if problem persists contact us\n"
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// estraggo i file di configurazione in una cartella temporanea
	workingDir, err := os.MkdirTemp("", fmt.Sprintf("tmp_workspace_%d", dw.Workspace.ID))
	if err != nil {
		dw.Workspace.Logs += "cannot create temp working directory, " + err.Error() + "\n"
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}
	defer os.RemoveAll(workingDir)

	configFilesLocation := path.Join(workingDir, dw.Workspace.GitRepoConfigurationFolder)
	err = utils.ExtractTarGz(dw.Workspace.WorkspaceConfigurationFiles, configFilesLocation)
	if err != nil {
		dw.Workspace.Logs += "cannot extract configuration files from targz archive, " + err.Error() + "\n"
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// caricamento della configurazione dal file .devcontainer.json
	devcontainerConfig := InitDevcontainerJson(dw.Workspace, path.Join(configFilesLocation, "devcontainer.json"))
	err = devcontainerConfig.LoadConfigFromFiles()
	if err != nil {
		dw.Workspace.Logs += err.Error() + "\n"
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// aggiustamento dei file di configurazione
	err = devcontainerConfig.FixConfigFiles()
}
