package config

type WorkspaceType struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	SupportedConfigSources []string `json:"supported_config_sources"`
}

func ListWorkspaceTypes() []WorkspaceType {
	return []WorkspaceType{
		{
			ID:   "docker_compose",
			Name: "Docker Compose",
			SupportedConfigSources: []string{
				"git",
				"template",
			},
		},
	}
}
