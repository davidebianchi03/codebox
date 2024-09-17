package common

import (
	"fmt"

	"codebox.com/db"
)

type CommonWorkspace struct {
	Workspace *db.Workspace
}

func (w *CommonWorkspace) StartWorkspace() error {
	return fmt.Errorf("not implemented")
}

func (w *CommonWorkspace) StopWorkspace() error {
	return fmt.Errorf("not implemented")
}

func (w *CommonWorkspace) DeleteWorkspace() error {
	return fmt.Errorf("not implemented")
}
