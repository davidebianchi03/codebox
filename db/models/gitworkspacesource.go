package models

import (
	"time"

	"gorm.io/gorm"
)

type GitWorkspaceSource struct {
	ID             uint   `gorm:"primarykey" json:"id"`
	RepositoryURL  string `gorm:"type:text;not null;" json:"repository_url"`
	RefName        string `gorm:"size:255; not null;" json:"ref_name"`
	ConfigFilePath string `gorm:"type:text;" json:"config_file_relative_path"` // path of the configuration files relative to the root of the repo
	SourcesID      uint
	Sources        *File          `json:"-"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
