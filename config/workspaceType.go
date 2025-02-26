package config

type WorkspaceType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ListWorkspaceTypes() []WorkspaceType {
	return []WorkspaceType{
		WorkspaceType{
			ID:   "docker_compose",
			Name: "Docker Compose",
		},
	}
}
