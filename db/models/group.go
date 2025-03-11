package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"size:255;unique"`
}
