package serializers

import "gitlab.com/codebox4073715/codebox/config"

type WorkspaceTypeSerializer struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	SupportedConfigSources []string `json:"supported_config_sources"`
	ConfigFilesDefaultPath string   `json:"config_files_default_path"`
}

func LoadWorkspaceTypeSerializer(workspaceType *config.WorkspaceType) *WorkspaceTypeSerializer {
	if workspaceType == nil {
		return nil
	}

	return &WorkspaceTypeSerializer{
		ID:                     workspaceType.ID,
		Name:                   workspaceType.Name,
		SupportedConfigSources: workspaceType.SupportedConfigSources,
		ConfigFilesDefaultPath: workspaceType.ConfigFilesDefaultPath,
	}
}

func LoadMultipleWorkspaceTypeSerializer(workspaceTypes []config.WorkspaceType) []WorkspaceTypeSerializer {
	serializers := make([]WorkspaceTypeSerializer, len(workspaceTypes))
	for i, workspaceType := range workspaceTypes {
		serializers[i] = *LoadWorkspaceTypeSerializer(&workspaceType)
	}
	return serializers
}
