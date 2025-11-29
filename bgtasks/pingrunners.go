package bgtasks

import (
	"time"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

/*
Ping all runners to check if they are online and get their version
*/
func (jobContext *Context) PingRunnersTask(job *work.Job) error {
	runners, err := models.ListRunners(-1, 0)
	if err != nil {
		// TODO: log error
		return nil
	}

	for _, runner := range runners {
		ri := runnerinterface.RunnerInterface{Runner: &runner}

		version, err := ri.GetRunnerVersion()
		if err == nil {
			now := time.Now()
			runner.Version = version
			runner.LastContact = &now
			dbconn.DB.Save(&runner)
		}
	}

	return nil
}
