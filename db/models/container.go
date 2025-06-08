package models

import (
	"time"

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
