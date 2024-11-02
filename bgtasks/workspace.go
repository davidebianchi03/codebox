package bgtasks

import (
	"fmt"

	"codebox.com/db"
	"codebox.com/workspaces/devcontainer"
	"github.com/gocraft/work"
)

type WorkspaceTaskContext struct {
	WorkspaceId uint
}

func (ctx *WorkspaceTaskContext) StartWorkspace(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	var workspace *db.Workspace
	result := db.DB.Where("ID=?", workspaceId).Preload("Owner").First(&workspace)
	if result.Error != nil {
		return fmt.Errorf("failed to retrieve workspace from db %s", result.Error)
	}

	workspace.ClearLogs()

	if workspace.Type == db.WorkspaceTypeDevcontainer {
		workspaceInterface := devcontainer.DevcontainerWorkspace{}
		workspaceInterface.Workspace = workspace
		workspaceInterface.StartWorkspace()
	} else {
		return fmt.Errorf("%s: unsupported workspace type", workspace.Type)
	}

	return nil
}

func (ctx *WorkspaceTaskContext) StopWorkspace(job *work.Job) error {
	workspaceId := job.ArgInt64("workspace_id")

	var workspace *db.Workspace
	result := db.DB.Where("ID=?", workspaceId).Preload("Owner").First(&workspace)
	if result.Error != nil {
		return fmt.Errorf("failed to retrieve workspace from db %s", result.Error)
	}

	workspace.ClearLogs()

	if workspace.Type == db.WorkspaceTypeDevcontainer {
		workspaceInterface := devcontainer.DevcontainerWorkspace{}
		workspaceInterface.Workspace = workspace
		workspaceInterface.StopWorkspace()
	} else {
		return fmt.Errorf("%s: unsupported workspace type", workspace.Type)
	}

	return nil
}

func (ctx *WorkspaceTaskContext) RestartWorkspace(job *work.Job) error {
	return nil
}

func (ctx *WorkspaceTaskContext) DeleteWorkspace(job *work.Job) error {
	return nil
}
