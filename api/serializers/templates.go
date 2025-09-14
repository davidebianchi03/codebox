package serializers

import "gitlab.com/codebox4073715/codebox/db/models"

type WorkspaceTemplateVersionSerializer struct {
	ID         uint   `json:"id"`
	TemplateID uint   `json:"template_id"`
	Name       string `json:"name"`
	Published  bool   `json:"published"`
}

func LoadWorkspaceTemplateVersionSerializer(templateVersion *models.WorkspaceTemplateVersion) *WorkspaceTemplateVersionSerializer {
	if templateVersion == nil {
		return nil
	}
	return &WorkspaceTemplateVersionSerializer{
		ID:         templateVersion.ID,
		TemplateID: templateVersion.TemplateID,
		Name:       templateVersion.Name,
		Published:  templateVersion.Published,
	}
}

func LoadMultipleWorkspaceTemplateVersionSerializer(templateVersions []models.WorkspaceTemplateVersion) []WorkspaceTemplateVersionSerializer {
	serializers := make([]WorkspaceTemplateVersionSerializer, len(templateVersions))
	for i, templateVersion := range templateVersions {
		serializers[i] = *LoadWorkspaceTemplateVersionSerializer(&templateVersion)
	}
	return serializers
}
