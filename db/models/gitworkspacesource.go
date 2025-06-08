package models

import (
	"time"

	"gorm.io/gorm"
)

type GitWorkspaceSource struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	RepositoryURL  string         `gorm:"column:repository_url; type:text;not null;" json:"repository_url"`
	RefName        string         `gorm:"column:ref_name; size:255; not null;" json:"ref_name"`
	ConfigFilePath string         `gorm:"column:config_file_path; type:text;" json:"config_file_relative_path"` // path of the configuration files relative to the root of the repo
	SourcesID      uint           `gorm:"column:sources_id;"`
	Sources        *File          `json:"-"`
	CreatedAt      time.Time      `gorm:"column:created_at;" json:"-"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;" json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
