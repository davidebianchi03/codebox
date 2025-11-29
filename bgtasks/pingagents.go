package bgtasks

import (
	"time"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

/*
Ping all workspaces agents to check if they are online
*/
func (jobContext *Context) PingAgentsTask(job *work.Job) error {
	containers, err := models.ListAllWorkspaceContainers()
	if err != nil {
		// TODO: log error
		return nil
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
