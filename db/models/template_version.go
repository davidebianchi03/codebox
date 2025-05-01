package models

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type WorkspaceTemplateVersion struct {
	gorm.Model
	ID             uint `gorm:"primarykey"`
	TemplateID     uint
	Template       *WorkspaceTemplate `gorm:"constraint:OnDelete:CASCADE;not null;"`
	Name           string             `gorm:"size:255;not null;"`
	ConfigFilePath string             `gorm:"type:text;" json:"config_file_relative_path"`
	SourcesID      uint
	Sources        *File `json:"-"`
	Published      bool  `gorm:"default:false"`
	EditedByID     uint
	EditedBy       *User `gorm:"constraint:OnDelete:SET NULL;"`
	EditedOn       time.Time
}

func ListWorkspaceTemplateVersionsByTemplate(template WorkspaceTemplate) (*[]WorkspaceTemplateVersion, error) {
	var tv *[]WorkspaceTemplateVersion
	if err := dbconn.DB.Find(
		&tv,
		map[string]interface{}{
			"template_id": template.ID,
		},
	).Error; err != nil {
		return nil, err
	}

	return tv, nil
}

func RetrieveWorkspaceTemplateVersionsByIdByTemplate(template WorkspaceTemplate, versionId uint) (*WorkspaceTemplateVersion, error) {
	var tv *WorkspaceTemplateVersion
	r := dbconn.DB.Find(
		&tv,
		map[string]interface{}{
			"id":          versionId,
			"template_id": template.ID,
		},
	)

	if r.Error != nil {
		return nil, r.Error
	}

	if r.RowsAffected == 0 {
		return nil, nil
	}

	return tv, nil
}

func CreateTemplateVersion(template WorkspaceTemplate, name string, user User) (*WorkspaceTemplateVersion, error) {
	// retrieve the latest version for the template
	var lastTemplateVersion *WorkspaceTemplateVersion
	r := dbconn.DB.Last(
		&lastTemplateVersion,
		map[string]interface{}{
			"template_id": template.ID,
		},
	)

	if r.Error != nil {
		return nil, r.Error
	}

	if r.RowsAffected == 0 {
		lastTemplateVersion = nil
	}

	// create source file item
	sourceFile := File{
		Filepath: path.Join("templates", fmt.Sprintf("%s.tar.gz", uuid.New().String())),
	}

	if sourceFile.Exists() {
		os.RemoveAll(sourceFile.GetAbsolutePath())
	}

	if err := dbconn.DB.Save(&sourceFile).Error; err != nil {
		return nil, err
	}

	// create object
	templateVersion := WorkspaceTemplateVersion{
		TemplateID:     template.ID,
		Template:       &template,
		Name:           name,
		ConfigFilePath: "",
		Sources:        &sourceFile,
		Published:      false,
		EditedByID:     user.ID,
		EditedBy:       &user,
		EditedOn:       time.Now(),
	}

	if err := dbconn.DB.Save(&templateVersion).Error; err != nil {
		return nil, err
	}

	// create or copy sources

	return nil, nil
}
