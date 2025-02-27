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

	"github.com/davidebianchi03/codebox/db/models"
)

type RunnerInterface struct {
	Runner *models.Runner
}

func (ri *RunnerInterface) getRunnerBaseUrl() string {
	if ri.Runner.UsePublicUrl {
		return ri.Runner.PublicUrl
	}

	// TODO: support for wsockporter
	panic("not implemented")
}

func (ri *RunnerInterface) GetRunnerVersion() (version string, err error) {
	url := fmt.Sprintf("%s/api/v1/version/", ri.getRunnerBaseUrl())
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", ri.Runner.Token)

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
	_ = writer.WriteField("config_file_name", "docker-compose.yml")

	_ = writer.WriteField("type", workspace.Type)
	_ = writer.WriteField("environment", strings.Join(workspace.EnvironmentVariables, ";"))

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

func (ri *RunnerInterface) GetDetails() {

}

func (ri *RunnerInterface) GetLogs() {

}

func (ri *RunnerInterface) StopWorkpace() {

}

func (ri *RunnerInterface) RemoveWorkspace() {

}
