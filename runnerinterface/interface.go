package runnerinterface

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/davidebianchi03/codebox/db/models"
	"github.com/davidebianchi03/codebox/proxy"
)

type RunnerInterface struct {
	Runner *models.Runner
}

func (ri *RunnerInterface) getRunnerBaseUrl() string {
	if ri.Runner.UsePublicUrl {
		return ri.Runner.PublicUrl
	}

	return fmt.Sprintf("http://127.0.0.1:%d", ri.Runner.Port)
}

func (ri *RunnerInterface) GetRunnerVersion() (version string, err error) {
	url := fmt.Sprintf("%s/api/v1/version/", ri.getRunnerBaseUrl())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("X-Codebox-Runner-Token", ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("recived status %d", res.StatusCode)
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
	_ = writer.WriteField("guid", strconv.Itoa(int(workspace.ID)))

	configFilePath := ""
	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		configFilePath, err = workspace.GitSource.GetConfigFileAbsPath()
		if err != nil {
			return err
		}
	} else {
		configFilePath, err = workspace.TemplateVersion.GetConfigFileAbsPath()
		if err != nil {
			return err
		}
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
	}
	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
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

	sshPrivateKeyFormPart, err := writer.CreateFormFile("ssh_private_key", "ssh_key")
	if err != nil {
		return err
	}

	sshPrivateKeyReader := bytes.NewReader([]byte(workspace.User.SshPrivateKey))
	_, err = io.Copy(sshPrivateKeyFormPart, sshPrivateKeyReader)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("X-Codebox-Runner-Token", ri.Runner.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("recived status code: %d", res.StatusCode)
	}

	return nil
}

func (ri *RunnerInterface) GetDetails(workspace *models.Workspace) (response RunnerWorkspaceStatusResponse, err error) {
	url := fmt.Sprintf("%s/api/v1/workspace/%d/?type=%s", ri.getRunnerBaseUrl(), workspace.ID, workspace.Type)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return RunnerWorkspaceStatusResponse{}, err
	}

	req.Header.Add("X-Codebox-Runner-Token", ri.Runner.Token)
	res, err := client.Do(req)
	if err != nil {
		return RunnerWorkspaceStatusResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return RunnerWorkspaceStatusResponse{}, fmt.Errorf("recived status %d", res.StatusCode)
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

func (ri *RunnerInterface) GetLogs(workspace *models.Workspace) (logs string, err error) {
	url := fmt.Sprintf("%s/api/v1/workspace/%d/logs", ri.getRunnerBaseUrl(), workspace.ID)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("X-Codebox-Runner-Token", ri.Runner.Token)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("recived status %d", res.StatusCode)
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

func (ri *RunnerInterface) StopWorkpace(workspace *models.Workspace) error {
	url := fmt.Sprintf("%s/api/v1/workspace/%d/stop", ri.getRunnerBaseUrl(), workspace.ID)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("type", workspace.Type)
	err := writer.Close()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, payload)

	if err != nil {
		return err
	}

	req.Header.Add("X-Codebox-Runner-Token", ri.Runner.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("recived status %d", res.StatusCode)
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

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, payload)

	if err != nil {
		return err
	}

	req.Header.Add("X-Codebox-Runner-Token", ri.Runner.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("recived status %d", res.StatusCode)
	}

	return nil
}

func (ri *RunnerInterface) PingAgent(container *models.WorkspaceContainer) bool {
	url := fmt.Sprintf(
		"%s/api/v1/agent-forward/?path=%s&workspace_guid=%s&container_id=%s",
		ri.getRunnerBaseUrl(),
		url.QueryEscape("/"),
		strconv.Itoa(int(container.WorkspaceID)),
		container.ContainerID,
	)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return false
	}

	req.Header.Add("X-CodeBox-Forward-Scheme", "ping")
	req.Header.Add("X-CodeBox-Forward-Host", "127.0.0.1")
	req.Header.Add("X-CodeBox-Forward-Port", "2222")
	req.Header.Add("X-CodeBox-Forward-Protocol", "http")
	req.Header.Add("X-Codebox-Runner-Token", ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return false
	}

	return true
}

func (ri *RunnerInterface) ForwardHttp(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	port *models.WorkspaceContainerPort,
	path string,
	rw http.ResponseWriter,
	req *http.Request,
) error {
	url := fmt.Sprintf(
		"%s/api/v1/agent-forward/?path=%s&workspace_guid=%s&container_id=%s",
		ri.getRunnerBaseUrl(),
		url.QueryEscape(path),
		strconv.Itoa(int(workspace.ID)),
		container.ContainerID,
	)

	proxyHeaders := http.Header{}
	proxyHeaders.Set("X-CodeBox-Forward-Host", "127.0.0.1")
	proxyHeaders.Set("X-CodeBox-Forward-Port", strconv.Itoa(int(port.PortNumber)))
	proxyHeaders.Set("X-CodeBox-Forward-Domain", "localhost")
	proxyHeaders.Set("X-CodeBox-Forward-Scheme", "http")
	proxyHeaders.Set("X-Codebox-Runner-Token", ri.Runner.Token)

	proxy, err := proxy.CreateReverseProxy(url, 30, 30, true, proxyHeaders)
	if err != nil {
		return err
	}

	proxy.ServeHTTP(rw, req)
	return nil
}

func (ri *RunnerInterface) ForwardSsh(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	rw http.ResponseWriter,
	req *http.Request,
) error {
	url := fmt.Sprintf(
		"%s/api/v1/agent-forward/?path=%s&workspace_guid=%s&container_id=%s",
		ri.getRunnerBaseUrl(),
		url.QueryEscape("/"),
		strconv.Itoa(int(workspace.ID)),
		container.ContainerID,
	)

	proxyHeaders := http.Header{}
	proxyHeaders.Set("X-CodeBox-Forward-Host", "127.0.0.1")
	proxyHeaders.Set("X-CodeBox-Forward-Port", "2222")
	proxyHeaders.Set("X-CodeBox-Forward-Domain", "localhost")
	proxyHeaders.Set("X-CodeBox-Forward-Scheme", "tcp_over_websocket")
	proxyHeaders.Set("X-Codebox-Runner-Token", ri.Runner.Token)

	proxy, err := proxy.CreateReverseProxy(url, 30, 30, true, proxyHeaders)
	if err != nil {
		return err
	}

	proxy.ServeHTTP(rw, req)
	return nil
}
