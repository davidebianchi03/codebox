package bgtasks

import (
	"time"

	dbconn "github.com/davidebianchi03/codebox/db/connection"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/davidebianchi03/codebox/runnerinterface"
	"github.com/gocraft/work"
)

func (jobContext *Context) PingAgents(job *work.Job) error {
	var containers []models.WorkspaceContainer
	if err := dbconn.DB.Preload("Workspace.Runner").Find(&containers).Error; err != nil {
		return err
	}

	for _, container := range containers {
		ri := runnerinterface.RunnerInterface{Runner: container.Workspace.Runner}
		if ri.PingAgent(&container) {
			now := time.Now()
			container.AgentLastContact = &now
			dbconn.DB.Save(&container)
		}
	}

	return nil
}
