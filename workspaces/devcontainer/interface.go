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
			dw.Workspace.AppendLogs(err.Error() + "\n")
			dw.Workspace.Status = db.WorkspaceStatusError
			db.DB.Save(&dw.Workspace)
			return
		}
		dw.Workspace.AppendLogs("\n")
		db.DB.Save(&dw.Workspace)
	}

	if dw.Workspace.WorkspaceConfigurationFiles == "" {
		dw.Workspace.AppendLogs("missing workspace configuration files, if problem persists contact us\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// estraggo i file di configurazione in una cartella temporanea
	workingDir, err := os.MkdirTemp("", fmt.Sprintf("tmp_workspace_%d", dw.Workspace.ID))
	if err != nil {
		dw.Workspace.AppendLogs("cannot create temp working directory, " + err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}
	defer os.RemoveAll(workingDir)

	workingDir = path.Join(workingDir, fmt.Sprintf("codebox_workspace_%d", dw.Workspace.ID))
	err = os.MkdirAll(workingDir, 0777)
	if err != nil {
		dw.Workspace.AppendLogs("cannot create temp working directory, " + err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	configFilesLocation := path.Join(workingDir, dw.Workspace.GitRepoConfigurationFolder)
	err = utils.ExtractTarGz(dw.Workspace.WorkspaceConfigurationFiles, configFilesLocation)
	if err != nil {
		dw.Workspace.AppendLogs("cannot extract configuration files from targz archive, " + err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// caricamento della configurazione dal file .devcontainer.json
	devcontainerConfig := InitDevcontainerJson(dw.Workspace, configFilesLocation)
	err = devcontainerConfig.LoadConfigFromFiles()
	if err != nil {
		dw.Workspace.AppendLogs(err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// aggiustamento dei file di configurazione
	err = devcontainerConfig.FixConfigFiles()
	if err != nil {
		dw.Workspace.AppendLogs(err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// creazione e avvio dei containers
	err = devcontainerConfig.GoUp()
	if err != nil {
		dw.Workspace.AppendLogs(err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// mapping dei containers
	err = devcontainerConfig.MapContainers()
	if err != nil {
		dw.Workspace.AppendLogs(err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// install and start agents
	dw.Workspace.AppendLogs("Starting agents...")
	err = devcontainerConfig.StartAgents()
	if err != nil {
		dw.Workspace.AppendLogs(err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// ping agents
	dw.Workspace.AppendLogs("Checking status of agents...")
	err = devcontainerConfig.CheckAgents()
	if err != nil {
		dw.Workspace.AppendLogs(err.Error() + "\n")
		dw.Workspace.Status = db.WorkspaceStatusError
		db.DB.Save(&dw.Workspace)
		return
	}

	// configure reverse proxy

	dw.Workspace.Status = db.WorkspaceStatusRunning
	db.DB.Save(&dw.Workspace)
}
