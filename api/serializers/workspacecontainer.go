package serializers

import (
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type WorkspaceContainerSerializer struct {
	ContainerID       string     `json:"container_id"`
	ContainerName     string     `json:"container_name"`
	ContainerImage    string     `json:"container_image"`
	ContainerUserID   uint       `json:"container_user_id"`
	ContainerUserName string     `json:"container_user_name"`
	AgentLastContact  *time.Time `json:"agent_last_contact"`
	WorkspacePath     string     `json:"workspace_path"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func LoadWorkspaceContainerSerializer(container *models.WorkspaceContainer) *WorkspaceContainerSerializer {
	if container == nil {
		return nil
	}

	return &WorkspaceContainerSerializer{
		ContainerID:       container.ContainerID,
		ContainerName:     container.ContainerName,
		ContainerImage:    container.ContainerImage,
		ContainerUserID:   container.ContainerUserID,
		ContainerUserName: container.ContainerUserName,
		AgentLastContact:  container.AgentLastContact,
		WorkspacePath:     container.WorkspacePath,
		CreatedAt:         container.CreatedAt,
		UpdatedAt:         container.UpdatedAt,
	}
}

func LoadMultipleWorkspaceContainerSerializers(containers []models.WorkspaceContainer) []WorkspaceContainerSerializer {
	serialized := make([]WorkspaceContainerSerializer, len(containers))
	for i, container := range containers {
		serialized[i] = *LoadWorkspaceContainerSerializer(&container)
	}
	return serialized
}
