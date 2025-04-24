package bgtasks

import (
	"time"

	dbconn "github.com/davidebianchi03/codebox/db/connection"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/davidebianchi03/codebox/runnerinterface"
	"github.com/gocraft/work"
)

func (jobContext *Context) PingRunners(job *work.Job) error {
	var runners []models.Runner
	if err := dbconn.DB.Find(&runners).Error; err != nil {
		return err
	}

	for _, runner := range runners {
		ri := runnerinterface.RunnerInterface{Runner: &runner}

		_, err := ri.GetRunnerVersion()
		if err == nil {
			now := time.Now()
			runner.LastContact = &now
			dbconn.DB.Save(&runner)
		}
	}

	return nil
}
