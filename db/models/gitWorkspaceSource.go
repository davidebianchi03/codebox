package models

import "gorm.io/gorm"

type GitWorkspaceSource struct {
	gorm.Model
	ID            uint   `gorm:"primarykey"`
	RepositoryURL string `gorm:"size:1024;unique;not null;"`
	Files         string `gorm:"size:1024;unique;not null;"`
}
