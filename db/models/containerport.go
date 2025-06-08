package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkspaceContainerPort struct {
	ID          uint               `gorm:"primarykey" json:"-"`
	ContainerID uint               `gorm:"column:container_id;" json:"-"`
	Container   WorkspaceContainer `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ServiceName string             `gorm:"column:service_name; size:255; not null;" json:"service_name"`
	PortNumber  uint               `gorm:"column:port_number; not null;" json:"port_number"`
	Public      bool               `gorm:"column:public; default:false;" json:"public"`
	CreatedAt   time.Time          `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
}
