package models

import (
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type GitWorkspaceSource struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	RepositoryURL  string         `gorm:"column:repository_url; type:text;not null;" json:"repository_url"`
	RefName        string         `gorm:"column:ref_name; size:255; not null;" json:"ref_name"`
	ConfigFilePath string         `gorm:"column:config_file_path; type:text; not null;" json:"config_file_relative_path"` // path of the configuration files relative to the root of the repo
	SourcesID      *uint          `gorm:"column:sources_id; default: null;" json:"-"`
	Sources        *File          `json:"-"`
	CreatedAt      time.Time      `gorm:"column:created_at;" json:"-"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;" json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func CreateGitWorkspaceSource(gitRepoUrl, gitRefName, configSourceFilePath string) (*GitWorkspaceSource, error) {
	gitSource := &GitWorkspaceSource{
		RepositoryURL:  gitRepoUrl,
		RefName:        gitRefName,
		ConfigFilePath: configSourceFilePath,
	}

	r := dbconn.DB.Create(gitSource)
	if r.Error != nil {
		return nil, r.Error
	}
	return gitSource, nil
}
