package models

import (
	"errors"
	"path"
	"strconv"

	"github.com/davidebianchi03/codebox/config"
	"gorm.io/gorm"
)

type WorkspaceTemplate struct {
	gorm.Model
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"size:255;unique;not null;"`
	Type        string `gorm:"size:255;"`
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

func (wtv *WorkspaceTemplateVersion) GetConfigFileAbsPath() (p string, err error) {
	if wtv.ID <= 0 || wtv.Files == "" {
		return "", errors.New("object does not exist")
	}
	p = path.Join(config.Environment.UploadsPath, "git-sources", strconv.Itoa(int(wtv.ID)))
	return p, nil
}
