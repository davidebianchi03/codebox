package bgtasks

import (
	"errors"

	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
Delete a runner:
- stop all workspaces running on this runner (ignoring errors)
- set workspace runner as nil
- delete the runner
*/
func (jobContext *Context) DeleteRunnerTask(job *work.Job) error {
	runnerId := job.ArgInt64("runner_id")

	runner, err := models.RetrieveRunnerByID(uint(runnerId))
	if err != nil {
		// TODO: log error
		return err
	}

	if runner == nil {
		// TODO: log error
		return errors.New("runner not found")
	}

	workspaces, err := models.ListWorkspacesByRunner(*runner)
	if err != nil {
		runner.DeletionInProgress = false
		models.UpdateRunner(*runner)
		// TODO: log error
		return err
	}

	// stop all workspaces and set runner as null
	for _, w := range workspaces {
		err := StopWorkspace(&w, true)
		if err != nil {
			runner.DeletionInProgress = false
			models.UpdateRunner(*runner)
			// TODO: log error
			return err
		}

		_, err = models.UpdateWorkspace(
			&w,
			w.Name,
			models.WorkspaceStatusStopped,
			nil,
			w.ConfigSource,
			w.TemplateVersion,
			w.GitSource,
			w.EnvironmentVariables,
		)
		if err != nil {
			runner.DeletionInProgress = false
			models.UpdateRunner(*runner)
			// TODO: log error
			return err
		}
	}

	if err := models.DeleteRunner(*runner); err != nil {
		runner.DeletionInProgress = false
		models.UpdateRunner(*runner)
		// TODO: log error
		return err
	}

	return nil
}
