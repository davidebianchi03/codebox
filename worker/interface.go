package worker

import "codebox.com/db"

type WorkerInterface struct {
	workspace *db.Workspace
}

func (i *WorkerInterface) GoUp() {

}
