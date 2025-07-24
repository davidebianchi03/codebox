package serializers

import (
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type WorkspaceContainerSerializer struct {
	ContainerID       string     `gorm:"column:container_id; size:255" json:"container_id"`
	ContainerName     string     `gorm:"column:container_name; size:255" json:"container_name"`
	ContainerImage    string     `gorm:"column:container_image; size:255" json:"container_image"`
	ContainerUserID   uint       `gorm:"column:container_user_id;" json:"container_user_id"`
	ContainerUserName string     `gorm:"size:255" json:"container_user_name"`
	AgentLastContact  *time.Time `gorm:"column:agent_last_contact;" json:"agent_last_contact"`
	WorkspacePath     string     `gorm:"column:workspace_path; size:255" json:"workspace_path"`
	CreatedAt         time.Time  `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;" json:"updated_at"`
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
