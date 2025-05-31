package models

import (
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type WorkspaceTemplate struct {
	ID          uint           `gorm:"primarykey"  json:"id"`
	Name        string         `gorm:"size:255;unique;not null;"  json:"name"`
	Type        string         `gorm:"size:255;"  json:"type"`
	Description string         `json:"description"`
	Icon        string         `gorm:"type:text;" json:"icon"`
	CreatedAt   time.Time      `gorm:"index" json:"-"`
	UpdatedAt   time.Time      `gorm:"index" json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Retrieve workspace template by id
// return nil if object is not found
func RetrieveWorkspaceTemplateByID(id uint) (*WorkspaceTemplate, error) {
	var wt *WorkspaceTemplate
	r := dbconn.DB.Model(WorkspaceTemplate{}).Find(&wt, map[string]interface{}{
		"id": id,
	})

	if r.Error != nil {
		return nil, r.Error
	}

	if r.RowsAffected == 0 {
		return nil, nil
	}
	return wt, nil
}

// Retrieve workspace template by name
// return nil if object is not found
func RetrieveWorkspaceTemplateByName(name string) (*WorkspaceTemplate, error) {
	var wt *WorkspaceTemplate
	r := dbconn.DB.Model(WorkspaceTemplate{}).Find(&wt, map[string]interface{}{
		"name": name,
	})

	if r.Error != nil {
		return nil, r.Error
	}

	if r.RowsAffected == 0 {
		return nil, nil
	}
	return wt, nil
}

// create new template
func CreateWorkspaceTemplate(name string, templateType string, description string, icon string) (*WorkspaceTemplate, error) {
	wt := WorkspaceTemplate{
		Name:        name,
		Type:        templateType,
		Description: description,
		Icon:        icon,
	}

	err := dbconn.DB.Save(&wt).Error
	if err != nil {
		return nil, err
	}

	return &wt, nil
}

// update workspace template
func UpdateWorkspaceTemplate(wt WorkspaceTemplate) error {
	err := dbconn.DB.Save(&wt).Error
	if err != nil {
		return err
	}

	return nil
}

func ListWorkspacesByTemplate(wt WorkspaceTemplate) ([]Workspace, error) {
	var workspaces []Workspace
	err := dbconn.DB.
		Preload("User").
		Preload("TemplateVersion").
		Joins("JOIN workspace_template_versions ON workspace_template_versions.id = workspaces.template_version_id").
		Joins("JOIN workspace_templates ON workspace_templates.id = workspace_template_versions.template_id").
		Where("workspace_templates.id = ?", wt.ID).
		Find(&workspaces).Error

	if err != nil {
		return []Workspace{}, err
	}

	return workspaces, nil
}

func DeleteTemplate(wt WorkspaceTemplate) error {
	return dbconn.DB.Unscoped().Delete(&wt).Error
}
