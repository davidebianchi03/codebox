package serializers

import (
	"fmt"
	"time"

	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
)

type WorkspaceContainerPort struct {
	ServiceName string    `json:"service_name"`
	PortNumber  uint      `json:"port_number"`
	Public      bool      `json:"public"`
	PortUrl     string    `json:"port_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func LoadWorkspaceContainerPort(port *models.WorkspaceContainerPort) *WorkspaceContainerPort {
	if port == nil {
		return nil
	}

	portUrl := ""
	if config.Environment.UseSubDomains {
		portUrl = fmt.Sprintf(
			"http://codebox--%d--%s--%d.%s",
			port.Container.WorkspaceID,
			port.Container.ContainerName,
			port.PortNumber,
			config.Environment.WildcardDomain,
		)
	} else {
		portUrl = fmt.Sprintf(
			"%s/views/port-forward/workspace/%d/container/%s/port/%d",
			config.Environment.ExternalUrl,
			port.Container.WorkspaceID,
			port.Container.ContainerName,
			port.PortNumber,
		)
	}

	return &WorkspaceContainerPort{
		ServiceName: port.ServiceName,
		PortNumber:  port.PortNumber,
		Public:      port.Public,
		PortUrl:     portUrl,
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
