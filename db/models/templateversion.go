package models

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/utils/targz"
	"gorm.io/gorm"
)

type WorkspaceTemplateVersion struct {
	ID             uint               `gorm:"primarykey" json:"id"`
	TemplateID     uint               `gorm:"column:template_id;" json:"template_id"`
	Template       *WorkspaceTemplate `gorm:"constraint:OnDelete:CASCADE;not null;" json:"-"`
	Name           string             `gorm:"column:name; size:255;not null;" json:"name"`
	ConfigFilePath string             `gorm:"column:config_file_path; type:text;" json:"config_file_relative_path"`
	SourcesID      uint               `gorm:"column:sources_id;" json:"-"`
	Sources        *File              `json:"-"`
	Published      bool               `gorm:"column:published; default:false" json:"published"`
	PublishedOn    *time.Time         `gorm:"column:published_on; default:null" json:"published_on"`
	EditedByID     uint               `gorm:"column:edited_by_id;" json:"-"`
	EditedBy       *User              `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
	EditedOn       time.Time          `gorm:"column:edited_on;" json:"edited_on"`
	CreatedAt      time.Time          `json:"-"`
	UpdatedAt      time.Time          `json:"-"`
	DeletedAt      gorm.DeletedAt     `gorm:"index" json:"-"`
}

func ListWorkspaceTemplateVersionsByTemplate(template WorkspaceTemplate) (*[]WorkspaceTemplateVersion, error) {
	var tv *[]WorkspaceTemplateVersion
	if err := dbconn.DB.
		Preload("Template").
		Preload("Sources").
		Find(
			&tv,
			map[string]interface{}{
				"template_id": template.ID,
			},
		).Error; err != nil {
		return nil, err
	}

	return tv, nil
}

func CountWorkspaceTemplateVersionsByTemplate(template WorkspaceTemplate) (int64, error) {
	var count int64
	if err := dbconn.DB.
		Model(&WorkspaceTemplateVersion{}).
		Where(
			map[string]interface{}{
				"template_id": template.ID,
			},
		).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func CountPublishedWorkspaceTemplateVersionsByTemplate(template WorkspaceTemplate) (int64, error) {
	var count int64
	if err := dbconn.DB.
		Model(&WorkspaceTemplateVersion{}).
		Where(
			map[string]interface{}{
				"template_id": template.ID,
				"published":   true,
			},
		).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func RetrieveWorkspaceTemplateVersionsById(versionId uint) (*WorkspaceTemplateVersion, error) {
	var tv *WorkspaceTemplateVersion
	r := dbconn.DB.
		Preload("Template").
		Preload("Sources").
		Find(
			&tv,
			map[string]interface{}{
				"id": versionId,
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

func RetrieveWorkspaceTemplateVersionsByIdByTemplate(template WorkspaceTemplate, versionId uint) (*WorkspaceTemplateVersion, error) {
	var tv *WorkspaceTemplateVersion
	r := dbconn.DB.
		Preload("Template").
		Preload("Sources").
		Find(
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

func RetrieveLatestTemplateVersionByTemplate(template WorkspaceTemplate) (*WorkspaceTemplateVersion, error) {
	count, err := CountPublishedWorkspaceTemplateVersionsByTemplate(template)
	if err != nil {
		return nil, err
	}

	var lastTemplateVersion *WorkspaceTemplateVersion

	if count > 0 {
		r := dbconn.DB.
			Preload("Sources").
			Last(
				&lastTemplateVersion,
				map[string]interface{}{
					"template_id": template.ID,
					"published":   true,
				},
			)

		if r.Error != nil {
			return nil, r.Error
		}

		return lastTemplateVersion, nil
	}

	return nil, nil
}

func CreateTemplateVersion(template WorkspaceTemplate, name string, user User, configFilePath string) (*WorkspaceTemplateVersion, error) {
	// retrieve the latest version for the template
	lastTemplateVersion, err := RetrieveLatestTemplateVersionByTemplate(template)
	if err != nil {
		return nil, err
	}

	// create source file item
	sourceFile := File{
		Filepath: path.Join("templates", fmt.Sprintf("%s.tar.gz", uuid.New().String())),
	}

	if sourceFile.Exists() {
		os.RemoveAll(sourceFile.GetAbsolutePath())
	}

	copySourcesFromLastVersion := false
	if lastTemplateVersion != nil {
		if lastTemplateVersion.Sources != nil {
			copySourcesFromLastVersion = lastTemplateVersion.Sources.Exists()
		}
	}

	if copySourcesFromLastVersion {
		// copy previous version
		src, err := os.Open(lastTemplateVersion.Sources.GetAbsolutePath())
		if err != nil {
			return nil, err
		}
		defer src.Close()

		destinationFile, err := os.Create(sourceFile.GetAbsolutePath())
		if err != nil {
			return nil, err
		}
		defer destinationFile.Close()

		_, err = io.Copy(destinationFile, src)
		if err != nil {
			return nil, err
		}
	} else {
		// create an empty tar.gz archive
		tgm := targz.TarGZManager{
			Filepath: sourceFile.GetAbsolutePath(),
		}

		if err := tgm.CreateArchive(); err != nil {
			return nil, err
		}
	}

	if err := dbconn.DB.Save(&sourceFile).Error; err != nil {
		return nil, err
	}

	// create object
	templateVersion := WorkspaceTemplateVersion{
		TemplateID:     template.ID,
		Template:       &template,
		Name:           name,
		ConfigFilePath: configFilePath,
		Sources:        &sourceFile,
		Published:      false,
		PublishedOn:    nil,
		EditedByID:     user.ID,
		EditedBy:       &user,
		EditedOn:       time.Now(),
	}

	if err := dbconn.DB.Save(&templateVersion).Error; err != nil {
		return nil, err
	}

	return &templateVersion, nil
}

func UpdateTemplateVersion(
	template WorkspaceTemplate,
	tv WorkspaceTemplateVersion,
	name string,
	published bool,
	user User,
	configFilePath string,
) (*WorkspaceTemplateVersion, error) {
	// check if template version exists
	templateVersion, err := RetrieveWorkspaceTemplateVersionsByIdByTemplate(template, tv.ID)
	if err != nil {
		return nil, err
	}

	if templateVersion == nil {
		return nil, errors.New("template version does not exists")
	}

	if templateVersion.Published && !published {
		return nil, errors.New("is not possible to unpublish a version")
	}

	if published && !templateVersion.Published {
		now := time.Now()
		templateVersion.PublishedOn = &now
	}

	templateVersion.Name = name
	templateVersion.Published = published
	templateVersion.EditedByID = user.ID
	templateVersion.EditedBy = &user
	templateVersion.EditedOn = time.Now()

	if err := dbconn.DB.Save(&templateVersion).Error; err != nil {
		return nil, err
	}

	return templateVersion, nil
}

func DeleteTemplateVersion(tv WorkspaceTemplateVersion) error {
	os.RemoveAll(tv.Sources.GetAbsolutePath())
	return dbconn.DB.Unscoped().Delete(&tv).Error
}
