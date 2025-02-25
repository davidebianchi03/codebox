package models

import "gorm.io/gorm"

type WorkspaceTemplate struct {
	gorm.Model
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"size:255;unique;not null;"`
	Description string
	Icon        string `gorm:"size:1024;"`
}

type WorkspaceTemplateVersion struct {
	gorm.Model
	ID         uint `gorm:"primarykey"`
	TemplateID uint
	Template   WorkspaceTemplate `gorm:"constraint:OnDelete:CASCADE;"`
	Name       string            `gorm:"size:1024;not null;"`
	Files      string            `gorm:"size:1024;not null;"`
}
