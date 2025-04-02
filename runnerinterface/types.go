package runnerinterface

type RunnerExposedPort struct {
	PortNumber  int    `json:"port_number"`
	ServiceName string `json:"service_name"`
	Public      bool   `json:"public"`
}

type RunnerContainer struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	State             string              `json:"state"`
	Image             string              `json:"image"`
	ContainerUserID   string              `json:"container_user"`
	ContainerUserName string              `json:"container_user_name"`
	ExposedPorts      []RunnerExposedPort `json:"exposed_ports"`
	WorkspacePath     string              `json:"workspace_path"`
}

type RunnerWorkspaceStatusResponse struct {
	Status     string            `json:"status"`
	Containers []RunnerContainer `json:"containers"`
}
