package bgtasks

import (
	"time"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

func (jobContext *Context) PingAgentsTask(job *work.Job) error {
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
