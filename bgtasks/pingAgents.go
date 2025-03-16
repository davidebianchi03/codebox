package bgtasks

import (
	"time"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	runnerinterface "github.com/davidebianchi03/codebox/runner-interface"
	"github.com/gocraft/work"
)

func (jobContext *Context) PingAgents(job *work.Job) error {
	var containers []models.WorkspaceContainer
	if err := db.DB.Preload("Workspace.Runner").Find(&containers).Error; err != nil {
		return err
	}

	for _, container := range containers {
		ri := runnerinterface.RunnerInterface{Runner: container.Workspace.Runner}
		if ri.PingAgent(&container) {
			container.AgentLastContact = time.Now()
			db.DB.Save(&container)
		}
	}

	return nil
}
