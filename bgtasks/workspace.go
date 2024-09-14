package bgtasks

import (
	"github.com/gocraft/work"
)

type WorkspaceTaskContext struct {
	// actionRequestedBy *db.User
	// actionRequestedOn time.Time
	workspaceId int
}

func (ctx *WorkspaceTaskContext) StartWorkspace(job *work.Job) error {
	return nil
}

func (ctx *WorkspaceTaskContext) StopWorkspace(job *work.Job) error {
	return nil
}

func (ctx *WorkspaceTaskContext) RestartWorkspace(job *work.Job) error {
	return nil
}

func (ctx *WorkspaceTaskContext) DeleteWorkspace(job *work.Job) error {
	return nil
}
