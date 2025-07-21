package config

type RunnerChoice struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	SupportedTypes []WorkspaceType `json:"supported_types"`
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
			},
		},
	}
}

func RetrieveRunnerTypeByID(id string) *RunnerChoice {
	for _, runner := range ListAvailableRunnerTypes() {
		if runner.ID == id {
			return &runner
		}
	}
	return nil
}
