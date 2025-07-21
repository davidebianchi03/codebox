package serializers

import (
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type WorkspaceSerializer struct {
	ID                   uint                                `json:"id"`
	Name                 string                              `json:"name"`
	User                 *UserSerializer                     `json:"user"`
	Status               string                              `json:"status"`
	Type                 string                              `json:"type"`
	Runner               *RunnerSerializer                   `json:"runner"`
	ConfigSource         string                              `json:"config_source"`
	TemplateVersion      *WorkspaceTemplateVersionSerializer `json:"template_version"`
	GitSource            *GitWorkspaceSourceSerializer       `json:"git_source"`
	EnvironmentVariables []string                            `json:"environment_variables"`
	CreatedAt            time.Time                           `json:"created_at"`
	UpdatedAt            time.Time                           `json:"updated_at"`
}

func LoadWorkspaceSerializer(workspace *models.Workspace) *WorkspaceSerializer {
	if workspace == nil {
		return nil
	}

	return &WorkspaceSerializer{
		ID:                   workspace.ID,
		Name:                 workspace.Name,
		User:                 LoadUserSerializer(workspace.User),
		Status:               workspace.Status,
		Type:                 workspace.Type,
		Runner:               LoadRunnerSerializer(workspace.Runner),
		ConfigSource:         workspace.ConfigSource,
		TemplateVersion:      LoadWorkspaceTemplateVersionSerializer(workspace.TemplateVersion),
		GitSource:            LoadGitWorkspaceSourceSerializer(workspace.GitSource),
		EnvironmentVariables: workspace.EnvironmentVariables,
		CreatedAt:            workspace.CreatedAt,
		UpdatedAt:            workspace.UpdatedAt,
	}
}

func LoadMultipleWorkspaceSerializer(workspaces []models.Workspace) []WorkspaceSerializer {
	serializers := make([]WorkspaceSerializer, len(workspaces))
	for i, workspace := range workspaces {
		serializers[i] = *LoadWorkspaceSerializer(&workspace)
	}
	return serializers
}
