package runnerinterface

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
)

func (ri *RunnerInterface) GetRunnerVersion() (version string, err error) {
	url := fmt.Sprintf("%s/api/v1/version/", ri.getRunnerBaseUrl())
	client := ri.getRequestsClient()

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("receivedstatus %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var jsonResponse map[string]interface{}
	if err = json.Unmarshal(body, &jsonResponse); err != nil {
		return "", err
	}

	version, ok := jsonResponse["version"].(string)
	if !ok {
		return "", errors.New("invalid response")
	}

	return version, nil
}

func (ri *RunnerInterface) StartWorkspace(workspace *models.Workspace) (err error) {
	url := fmt.Sprintf("%s/api/v1/workspace/", ri.getRunnerBaseUrl())

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("id", strconv.Itoa(int(workspace.ID)))

	configFilePath := ""
	if workspace.ConfigSource == models.WorkspaceConfigSourceGit && workspace.GitSource != nil {
		if workspace.GitSource.Sources == nil {
			return errors.New("source files do not exists")
		}
		configFilePath = workspace.GitSource.Sources.GetAbsolutePath()
	} else {
		if workspace.TemplateVersion.Sources == nil {
			return errors.New("source files do not exists")
		}
		configFilePath = workspace.TemplateVersion.Sources.GetAbsolutePath()
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	configFileFormPart, err := writer.CreateFormFile("config", "config.tar.gz")
	if err != nil {
		return err
	}

	_, err = io.Copy(configFileFormPart, configFile)
	if err != nil {
		return err
	}

	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		_ = writer.WriteField("git_repository_url", workspace.GitSource.RepositoryURL)
		_ = writer.WriteField("config_file_name", workspace.GitSource.ConfigFilePath)
	} else {
		_ = writer.WriteField("config_file_name", workspace.TemplateVersion.ConfigFilePath)
	}
	_ = writer.WriteField("type", workspace.Type)
	_ = writer.WriteField("git_user_name", fmt.Sprintf("%s %s", workspace.User.FirstName, workspace.User.LastName))
	_ = writer.WriteField("git_user_email", workspace.User.Email)

	// add default variables to environment
	_ = writer.WriteField(
		"environment",
		strings.Join(
			append(
				workspace.EnvironmentVariables,
				workspace.GetDefaultEnvironmentVariables()...,
			),
			";",
		),
	)

	err = writer.Close()
	if err != nil {
		return err
	}

	client := ri.getRequestsClient()
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return err
	}

	req.Header.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("receivedstatus code: %d", res.StatusCode)
	}

	return nil
}

func (ri *RunnerInterface) GetWorkspaceDetails(workspace *models.Workspace) (response RunnerWorkspaceStatusResponse, err error) {
	url := fmt.Sprintf("%s/api/v1/workspace/%d/?type=%s", ri.getRunnerBaseUrl(), workspace.ID, workspace.Type)

	client := ri.getRequestsClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return RunnerWorkspaceStatusResponse{}, err
	}

	req.Header.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)
	res, err := client.Do(req)
	if err != nil {
		return RunnerWorkspaceStatusResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return RunnerWorkspaceStatusResponse{}, fmt.Errorf("receivedstatus %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RunnerWorkspaceStatusResponse{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return RunnerWorkspaceStatusResponse{}, err
	}

	return response, nil
}

func (ri *RunnerInterface) GetWorkspaceLogs(workspace *models.Workspace) (logs string, err error) {
	url := fmt.Sprintf("%s/api/v1/workspace/%d/logs", ri.getRunnerBaseUrl(), workspace.ID)

	client := ri.getRequestsClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("receivedstatus %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var parsedBody map[string]string
	if err = json.Unmarshal(body, &parsedBody); err != nil {
		return "", err
	}

	return parsedBody["logs"], nil
}

func (ri *RunnerInterface) StopWorkspace(workspace *models.Workspace) error {
	url := fmt.Sprintf("%s/api/v1/workspace/%d/stop/", ri.getRunnerBaseUrl(), workspace.ID)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("type", workspace.Type)
	err := writer.Close()
	if err != nil {
		return err
	}

	client := ri.getRequestsClient()
	req, err := http.NewRequest(http.MethodPost, url, payload)

	if err != nil {
		return err
	}

	req.Header.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("receivedstatus %d", res.StatusCode)
	}

	return nil
}

func (ri *RunnerInterface) RemoveWorkspace(workspace *models.Workspace) error {
	url := fmt.Sprintf("%s/api/v1/workspace/%d/remove", ri.getRunnerBaseUrl(), workspace.ID)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("type", workspace.Type)
	err := writer.Close()
	if err != nil {
		return err
	}

	client := ri.getRequestsClient()
	req, err := http.NewRequest(http.MethodPost, url, payload)

	if err != nil {
		return err
	}

	req.Header.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("receivedstatus %d", res.StatusCode)
	}

	return nil
}
