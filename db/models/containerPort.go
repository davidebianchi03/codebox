package models

import "gorm.io/gorm"

type WorkspaceContainerPort struct {
	gorm.Model
	ID          uint `gorm:"primarykey"`
	ContainerID uint
	Container   WorkspaceContainer `gorm:"constraint:OnDelete:CASCADE;"`
	PortNumber  uint               `gorm:"not null;"`
	Public      bool               `gorm:"default:false;"`
}
