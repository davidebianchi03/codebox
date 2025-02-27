package models

import (
	"time"

	"gorm.io/gorm"
)

type Runner struct {
	gorm.Model
	ID            uint    `gorm:"primarykey"`
	Name          string  `gorm:"size:255;unique;not null;"`
	Token         string  `gorm:"size:1024;unique;not null;"`
	Type          string  `gorm:"size:255;"`
	Restricted    bool    `gorm:"default:false;"`
	AllowedGroups []Group `gorm:"many2many:runner_allowed_groups;"`
	UsePublicUrl  bool    `gorm:"default:false;"`
	PublicUrl     string  `gorm:"size:1024;"`
	LastContact   time.Time
}
