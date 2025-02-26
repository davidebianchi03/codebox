package config

type RunnerChoice struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	SupportedTypes []WorkspaceType `json:"supported_type"`
}

func ListAvailableWorkspaceTypes() []RunnerChoice {
	return []RunnerChoice{
		RunnerChoice{
			ID:          "docker",
			Name:        "Docker",
			Description: "Runner for docker containers based environments",
			SupportedTypes: []WorkspaceType{
				WorkspaceType{
					ID:   "docker_compose",
					Name: "Docker Compose",
				},
			},
		},
	}
}
