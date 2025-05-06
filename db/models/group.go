package models

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"size:255;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
