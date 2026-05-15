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

type ContainerFileInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"is_dir"`
	Size     int64  `json:"size"`
	Mode     string `json:"mode"`
	ModTime  int64  `json:"mod_time"`
	Owner    string `json:"owner,omitempty"`
	Group    string `json:"group,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

type ContainerReadFileResponse struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

type ContainerSystemInfo struct {
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	HomeDir     string `json:"home_dir"`
	CurrentUser string `json:"current_user"`
	UserID      string `json:"user_id,omitempty"`
	GroupID     string `json:"group_id,omitempty"`
	NumCPU      int    `json:"num_cpu"`
}

type ExecuteCommandResponse struct {
	Command    string `json:"command"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	WasSuccess bool   `json:"was_success"`
}
