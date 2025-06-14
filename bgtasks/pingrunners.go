package bgtasks

import (
	"time"

	"github.com/gocraft/work"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

func (jobContext *Context) PingRunners(job *work.Job) error {
	var runners []models.Runner
	if err := dbconn.DB.Find(&runners).Error; err != nil {
		return err
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
