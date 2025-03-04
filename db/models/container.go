package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkspaceContainer struct {
	ID                uint           `gorm:"primarykey" json:"-"`
	WorkspaceID       uint           `json:"-"`
	Workspace         Workspace      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContainerID       string         `gorm:"size:255" json:"container_id"`
	ContainerName     string         `gorm:"size:255" json:"container_name"`
	ContainerImage    string         `gorm:"size:255" json:"container_image"`
	ContainerUserID   uint           `json:"container_user_id"`
	ContainerUserName string         `gorm:"size:255" json:"container_user_name"`
	AgentLastContact  time.Time      `json:"agent_last_contact"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
