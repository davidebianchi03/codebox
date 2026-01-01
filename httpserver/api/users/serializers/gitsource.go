package serializers

import "gitlab.com/codebox4073715/codebox/db/models"

type GitWorkspaceSourceSerializer struct {
	RepositoryURL  string `json:"repository_url"`
	RefName        string `json:"ref_name"`
	ConfigFilePath string `json:"config_file_relative_path"`
}

func LoadGitWorkspaceSourceSerializer(gitSource *models.GitWorkspaceSource) *GitWorkspaceSourceSerializer {
	if gitSource == nil {
		return nil
	}

	return &GitWorkspaceSourceSerializer{
		RepositoryURL:  gitSource.RepositoryURL,
		RefName:        gitSource.RefName,
		ConfigFilePath: gitSource.ConfigFilePath,
	}
}

func LoadMultipleGitWorkspaceSourceSerializer(gitSources []models.GitWorkspaceSource) []GitWorkspaceSourceSerializer {
	serializers := make([]GitWorkspaceSourceSerializer, len(gitSources))
	for i, gitSource := range gitSources {
		serializers[i] = *LoadGitWorkspaceSourceSerializer(&gitSource)
	}
	return serializers
}
