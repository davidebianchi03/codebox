package bgtasks

import (
	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/logging"
)

/*
bg task that deletes old logs
*/
func (jobContext *Context) RotateLogsTask(job *work.Job) error {
	logging.RotateLogs()
	return nil
}
