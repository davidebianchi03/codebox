package devcontainer

import "codebox.com/workspaces/common"

type DevcontainerWorkspace struct {
	common.CommonWorkspace
}

func (dw *DevcontainerWorkspace) StartWorkspace() {
	_ = common.RetrieveWorkspaceConfigFilesFromGitRepo(dw.Workspace)
}
