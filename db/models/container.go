package models

import (
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type WorkspaceContainer struct {
	ID                uint           `gorm:"primarykey" json:"-"`
	WorkspaceID       uint           `gorm:"column:workspace_id;" json:"-"`
	Workspace         Workspace      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContainerID       string         `gorm:"column:container_id; size:255" json:"container_id"`
	ContainerName     string         `gorm:"column:container_name; size:255" json:"container_name"`
	ContainerImage    string         `gorm:"column:container_image; size:255" json:"container_image"`
	ContainerUserID   uint           `gorm:"column:container_user_id;" json:"container_user_id"`
	ContainerUserName string         `gorm:"size:255" json:"container_user_name"`
	AgentLastContact  *time.Time     `gorm:"column:agent_last_contact;" json:"agent_last_contact"`
	WorkspacePath     string         `gorm:"column:workspace_path; size:255" json:"workspace_path"`
	CreatedAt         time.Time      `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

/*
ListWorkspaceContainersByWorkspace retrieves all containers for a given workspace.
*/
func ListWorkspaceContainersByWorkspace(workspace Workspace) ([]WorkspaceContainer, error) {
	var containers []WorkspaceContainer
	result := dbconn.DB.Where("workspace_id = ?", workspace.ID).Find(&containers)
	if result.Error != nil {
		return nil, result.Error
	}
	return containers, nil
}

/*
RetrieveWorkspaceContainerByName retrieves a specific container by name in a workspace.
*/
func RetrieveWorkspaceContainerByName(workspace Workspace, containerName string) (*WorkspaceContainer, error) {
	var container WorkspaceContainer
	result := dbconn.DB.Where("workspace_id = ? AND container_name = ?", workspace.ID, containerName).First(&container)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &container, nil
}
