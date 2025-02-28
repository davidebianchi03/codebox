package models

import (
	"time"

	"gorm.io/gorm"
)

type Runner struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"size:255;unique;not null;" json:"name"`
	Token         string         `gorm:"size:1024;unique;not null;" json:"-"`
	Type          string         `gorm:"size:255;" json:"type"`
	Restricted    bool           `gorm:"default:false;" json:"-"`
	AllowedGroups []Group        `gorm:"many2many:runner_allowed_groups;" json:"-"`
	UsePublicUrl  bool           `gorm:"default:false;" json:"-"`
	PublicUrl     string         `gorm:"size:1024;" json:"-"`
	LastContact   time.Time      `json:"last_contact"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
