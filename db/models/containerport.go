package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkspaceContainerPort struct {
	ID          uint               `gorm:"primarykey" json:"-"`
	ContainerID uint               `json:"-"`
	Container   WorkspaceContainer `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ServiceName string             `gorm:"size:255; not null;" json:"service_name"`
	PortNumber  uint               `gorm:"not null;" json:"port_number"`
	Public      bool               `gorm:"default:false;" json:"public"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
}
