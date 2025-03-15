package bgtasks

import (
	"time"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	runnerinterface "github.com/davidebianchi03/codebox/runner-interface"
	"github.com/gocraft/work"
)

func (jobContext *Context) PingRunners(job *work.Job) error {
	var runners []models.Runner
	if err := db.DB.Find(&runners).Error; err != nil {
		return err
	}

	for _, runner := range runners {
		ri := runnerinterface.RunnerInterface{Runner: &runner}

		_, err := ri.GetRunnerVersion()
		if err == nil {
			runner.LastContact = time.Now()
			db.DB.Save(&runner)
		}
	}

	return nil
}
