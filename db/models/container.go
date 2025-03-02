package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkspaceContainer struct {
	gorm.Model
	ID                uint `gorm:"primarykey"`
	WorkspaceID       uint
	Workspace         Workspace `gorm:"constraint:OnDelete:CASCADE;"`
	ContainerID       string    `gorm:"size:255"`
	ContainerName     string    `gorm:"size:255"`
	ContainerImage    string    `gorm:"size:255"`
	ContainerUserID   uint
	ContainerUserName string `gorm:"size:255"`
	AgentLastContact  time.Time
}
