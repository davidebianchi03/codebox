package config

type WorkspaceType struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	SupportedConfigSources []string `json:"supported_config_sources"`
	ConfigFilesDefaultPath string   `json:"config_files_default_path"`
}

// Retrieve the list of all workspace types
func ListWorkspaceTypes() []WorkspaceType {
	return []WorkspaceType{
		{
			ID:   "docker_compose",
			Name: "Docker Compose",
			SupportedConfigSources: []string{
				"git",
				"template",
			},
			ConfigFilesDefaultPath: "docker-compose.yml",
		},
		{
			ID:   "devcontainer",
			Name: "Dev Container",
			SupportedConfigSources: []string{
				"git",
			},
			ConfigFilesDefaultPath: ".devcontainer",
		},
	}
}

// Retrieve a specific workspace type by ID
// Returns nil if the type is not found
func RetrieveWorkspaceType(id string) *WorkspaceType {
	for _, t := range ListWorkspaceTypes() {
		if t.ID == id {
			return &t
		}
	}
	return nil
}
