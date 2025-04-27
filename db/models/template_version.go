package models

import (
	"errors"
	"path"
	"strconv"

	"gitlab.com/codebox4073715/codebox/config"
	"gorm.io/gorm"
)

type WorkspaceTemplateVersion struct {
	gorm.Model
	ID             uint `gorm:"primarykey"`
	TemplateID     uint
	Template       *WorkspaceTemplate `gorm:"constraint:OnDelete:CASCADE;not null;"`
	Name           string             `gorm:"size:255;not null;"`
	ConfigFilePath string             `gorm:"type:text;" json:"config_file_relative_path"`
	Files          string             `gorm:"type:text;not null;"`
}

func (wtv *WorkspaceTemplateVersion) GetConfigFileAbsPath() (p string, err error) {
	if wtv.ID <= 0 || wtv.Files == "" {
		return "", errors.New("object does not exist")
	}
	p = path.Join(config.Environment.UploadsPath, "git-sources", strconv.Itoa(int(wtv.ID)))
	return p, nil
}
