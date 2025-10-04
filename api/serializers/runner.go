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

// AdminRunnersSerializer is used for admin-specific runner information
type AdminRunnersSerializer struct {
	ID                 uint       `json:"id"`
	Name               string     `json:"name"`
	Type               string     `json:"type"`
	LastContact        *time.Time `json:"last_contact"`
	UsePublicUrl       bool       `json:"use_public_url"`
	PublicUrl          string     `json:"public_url"`
	DeletionInProgress bool       `json:"deletion_in_progress"`
}

func LoadAdminRunnerSerializer(runner *models.Runner) *AdminRunnersSerializer {
	if runner == nil {
		return nil
	}
	return &AdminRunnersSerializer{
		ID:                 runner.ID,
		Name:               runner.Name,
		Type:               runner.Type,
		LastContact:        runner.LastContact,
		UsePublicUrl:       runner.UsePublicUrl,
		PublicUrl:          runner.PublicUrl,
		DeletionInProgress: runner.DeletionInProgress,
	}
}

func LoadMultipleAdminRunnerSerializer(runners []models.Runner) []AdminRunnersSerializer {
	serializers := make([]AdminRunnersSerializer, len(runners))
	for i, runner := range runners {
		serializers[i] = *LoadAdminRunnerSerializer(&runner)
	}
	return serializers
}
