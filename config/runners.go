package config

type RunnerChoice struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	SupportedTypes []WorkspaceType `json:"supported_type"`
}

func ListAvailableRunnerTypes() []RunnerChoice {
	return []RunnerChoice{
		{
			ID:          "docker",
			Name:        "Docker",
			Description: "Runner for docker containers based environments",
			SupportedTypes: []WorkspaceType{
				{
					ID:   "docker_compose",
					Name: "Docker Compose",
					SupportedConfigSources: []string{
						"git",
						"template",
					},
				},
			},
		},
	}
}
