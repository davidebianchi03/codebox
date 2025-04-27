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
	err := dbconn.DB.Model(WorkspaceTemplate{}).Find(&wt, map[string]interface{}{
		"id": id,
	}).Error
	return wt, err
}

// Retrieve workspace template by name
// return nil if object is not found
func RetrieveWorkspaceTemplateByName(name string) (*WorkspaceTemplate, error) {
	var wt *WorkspaceTemplate
	err := dbconn.DB.Model(WorkspaceTemplate{}).Find(&wt, map[string]interface{}{
		"name": name,
	}).Error
	return wt, err
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
