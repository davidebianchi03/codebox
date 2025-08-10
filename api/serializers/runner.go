package serializers

import (
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type RunnerSerializer struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	LastContact *time.Time `json:"last_contact"`
}

func LoadRunnerSerializer(runner *models.Runner) *RunnerSerializer {
	if runner == nil {
		return nil
	}
	return &RunnerSerializer{
		ID:          runner.ID,
		Name:        runner.Name,
		Type:        runner.Type,
		LastContact: runner.LastContact,
	}
}

func LoadMultipleRunnerSerializer(runners []models.Runner) []RunnerSerializer {
	serializers := make([]RunnerSerializer, len(runners))
	for i, runner := range runners {
		serializers[i] = *LoadRunnerSerializer(&runner)
	}
	return serializers
}
