package serializers

import (
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type WorkspaceContainerPort struct {
	ServiceName string    `json:"service_name"`
	PortNumber  uint      `json:"port_number"`
	Public      bool      `json:"public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func LoadWorkspaceContainerPort(port *models.WorkspaceContainerPort) *WorkspaceContainerPort {
	if port == nil {
		return nil
	}

	return &WorkspaceContainerPort{
		ServiceName: port.ServiceName,
		PortNumber:  port.PortNumber,
		Public:      port.Public,
		CreatedAt:   port.CreatedAt,
		UpdatedAt:   port.UpdatedAt,
	}
}

func LoadMultipleWorkspaceContainerPorts(ports []models.WorkspaceContainerPort) []WorkspaceContainerPort {
	serialized := make([]WorkspaceContainerPort, len(ports))
	for i, port := range ports {
		serialized[i] = *LoadWorkspaceContainerPort(&port)
	}
	return serialized
}
