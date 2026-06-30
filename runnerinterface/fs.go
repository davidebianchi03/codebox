package runnerinterface

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
)

var (
	ErrorPathIsADir        = errors.New("selected path is a directory")
	ErrorPathIsNotADir     = errors.New("selected path is not a directory")
	ErrorPathAlreadyExists = errors.New("path already exists")
	ErrorPermissionDenied  = errors.New("permission denied")
	ErrorPathNotExist      = errors.New("path does not exist")
	ErrorInvalidFileMode   = errors.New("invalid permissions format")
	ErrorInvalidBase64     = errors.New("invalid base64 encoded string")
)

func IsPathIsNotADir(err error) bool {
	return err == ErrorPathIsNotADir
}

func IsPermissionDenied(err error) bool {
	return err == ErrorPermissionDenied
}

func IsPathNotExist(err error) bool {
	return err == ErrorPathNotExist
}

func IsErrorPathAlreadyExists(err error) bool {
	return err == ErrorPathAlreadyExists
}

func IsErrorPathIsADir(err error) bool {
	return err == ErrorPathIsADir
}

func IsErrorInvalidFileMode(err error) bool {
	return err == ErrorInvalidFileMode
}

func IsErrorInvalidBase64(err error) bool {
	return err == ErrorInvalidBase64
}

/*
List content of a directory
*/
func (ri *RunnerInterface) ContainerFsListDir(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
) ([]ContainerFileInfo, error) {
	if workspace.Runner.Port == 0 {
		return []ContainerFileInfo{}, errors.New("runner is not connected")
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/list-directory?path=%s",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
		url.QueryEscape(path),
	)

	client := ri.getRequestsClient()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []ContainerFileInfo{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return []ContainerFileInfo{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []ContainerFileInfo{}, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		switch res.StatusCode {
		case 400:
			return []ContainerFileInfo{}, ErrorPathIsNotADir
		case 403:
			return []ContainerFileInfo{}, ErrorPermissionDenied
		case 404:
			return []ContainerFileInfo{}, ErrorPathNotExist
		default:
			return []ContainerFileInfo{}, fmt.Errorf("server returned status %d", res.StatusCode)
		}
	}

	var data []ContainerFileInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return []ContainerFileInfo{}, errors.New("failed to parse runner response")
	}
	return data, nil
}

/*
Retrieve details of a file/folder
*/
func (ri *RunnerInterface) ContainerFsGetItemInfo(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
) (ContainerFileInfo, error) {
	if workspace.Runner.Port == 0 {
		return ContainerFileInfo{}, errors.New("runner is not connected")
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/get-item-info?path=%s",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
		url.QueryEscape(path),
	)

	client := ri.getRequestsClient()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ContainerFileInfo{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return ContainerFileInfo{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ContainerFileInfo{}, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		switch res.StatusCode {
		case 400:
			return ContainerFileInfo{}, ErrorPathIsNotADir
		case 403:
			return ContainerFileInfo{}, ErrorPermissionDenied
		case 404:
			return ContainerFileInfo{}, ErrorPathNotExist
		default:
			return ContainerFileInfo{}, fmt.Errorf("server returned status %d", res.StatusCode)
		}
	}

	var data ContainerFileInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return ContainerFileInfo{}, errors.New("failed to parse runner response")
	}
	return data, nil
}

/*
Create a directory
*/
func (ri *RunnerInterface) ContainerFsCreateDir(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
	permissions string,
) (ContainerFileInfo, error) {
	if !validatePermissionString(permissions) {
		return ContainerFileInfo{}, ErrorInvalidFileMode
	}

	if workspace.Runner.Port == 0 {
		return ContainerFileInfo{}, errors.New("runner is not connected")
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/create-directory",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
	)

	client := ri.getRequestsClient()

	requestBody, err := json.Marshal(map[string]interface{}{
		"path":        path,
		"permissions": permissions,
	})

	if err != nil {
		return ContainerFileInfo{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return ContainerFileInfo{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return ContainerFileInfo{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ContainerFileInfo{}, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		switch res.StatusCode {
		case 400:
			return ContainerFileInfo{}, ErrorPathAlreadyExists
		case 403:
			return ContainerFileInfo{}, ErrorPermissionDenied
		case 404:
			return ContainerFileInfo{}, ErrorPathNotExist
		default:
			return ContainerFileInfo{}, fmt.Errorf("server returned status %d", res.StatusCode)
		}
	}

	var data ContainerFileInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return ContainerFileInfo{}, errors.New("failed to parse runner response")
	}
	return data, nil
}

/*
Delete a file or directory (recursive for directories)
*/
func (ri *RunnerInterface) ContainerFsDeleteItem(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
) error {
	if workspace.Runner.Port == 0 {
		return errors.New("runner is not connected")
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/delete",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
	)

	client := ri.getRequestsClient()

	requestBody, err := json.Marshal(map[string]interface{}{
		"path": path,
	})

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		switch res.StatusCode {
		case 403:
			return ErrorPermissionDenied
		case 404:
			return ErrorPathNotExist
		default:
			return fmt.Errorf("server returned status %d", res.StatusCode)
		}
	}

	return nil
}

/*
Move/rename a file or directory
*/
func (ri *RunnerInterface) ContainerFsRenameItem(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
	newPath string,
) error {
	if workspace.Runner.Port == 0 {
		return errors.New("runner is not connected")
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/rename",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
	)

	client := ri.getRequestsClient()

	requestBody, err := json.Marshal(map[string]interface{}{
		"path":     path,
		"new_path": newPath,
	})

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		switch res.StatusCode {
		case 400:
			return ErrorPathAlreadyExists
		case 403:
			return ErrorPermissionDenied
		case 404:
			return ErrorPathNotExist
		default:
			return fmt.Errorf("server returned status %d", res.StatusCode)
		}
	}

	return nil
}

/*
Read a file
*/
func (ri *RunnerInterface) ContainerFsReadFile(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
) (ContainerReadFileResponse, error) {
	if workspace.Runner.Port == 0 {
		return ContainerReadFileResponse{}, errors.New("runner is not connected")
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/read-file?path=%s",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
		url.QueryEscape(path),
	)

	client := ri.getRequestsClient()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ContainerReadFileResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return ContainerReadFileResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ContainerReadFileResponse{}, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		switch res.StatusCode {
		case 403:
			return ContainerReadFileResponse{}, ErrorPermissionDenied
		case 404:
			return ContainerReadFileResponse{}, ErrorPathNotExist
		case 409:
			return ContainerReadFileResponse{}, ErrorPathIsADir
		default:
			return ContainerReadFileResponse{}, fmt.Errorf("server returned status %d", res.StatusCode)
		}
	}

	var data ContainerReadFileResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return ContainerReadFileResponse{}, errors.New("failed to parse runner response")
	}
	return data, nil
}

/*
Write a file
*/
func (ri *RunnerInterface) ContainerFsWriteFile(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
	content string, // base 64
	permissions string,
) error {
	if workspace.Runner.Port == 0 {
		return errors.New("runner is not connected")
	}

	if !validatePermissionString(permissions) {
		return ErrorInvalidFileMode
	}

	if _, err := base64.StdEncoding.DecodeString(content); err != nil {
		return ErrorInvalidBase64
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/write-file",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
	)

	client := ri.getRequestsClient()

	requestBody, err := json.Marshal(map[string]interface{}{
		"path":        path,
		"content":     content,
		"permissions": permissions,
	})

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		switch res.StatusCode {
		case 403:
			return ErrorPermissionDenied
		case 404:
			return ErrorPathNotExist
		case 409:
			return ErrorPathIsADir
		default:
			return fmt.Errorf("server returned status %d", res.StatusCode)
		}
	}

	return nil
}

/*
Get system info
*/
func (ri *RunnerInterface) ContainerFsGetSystemInfo(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
) (ContainerSystemInfo, error) {
	if workspace.Runner.Port == 0 {
		return ContainerSystemInfo{}, errors.New("runner is not connected")
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/system-info",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
	)

	client := ri.getRequestsClient()
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return ContainerSystemInfo{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return ContainerSystemInfo{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ContainerSystemInfo{}, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return ContainerSystemInfo{}, fmt.Errorf("server returned status %d", res.StatusCode)
	}

	var data ContainerSystemInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return ContainerSystemInfo{}, errors.New("failed to parse runner response")
	}
	return data, nil
}

/*
Execute command in container
*/
func (ri *RunnerInterface) ContainerFSExecuteCommand(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	command string,
	args []string,
	workDir string,
) (ExecuteCommandResponse, error) {
	if workspace.Runner.Port == 0 {
		return ExecuteCommandResponse{}, errors.New("runner is not connected")
	}

	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/fs/execute-command",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
	)

	requestBody, err := json.Marshal(map[string]interface{}{
		"command":  command,
		"args":     args,
		"work_dir": workDir,
	})

	client := ri.getRequestsClient()
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return ExecuteCommandResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return ExecuteCommandResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ExecuteCommandResponse{}, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return ExecuteCommandResponse{}, fmt.Errorf("server returned status %d", res.StatusCode)
	}

	var data ExecuteCommandResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return ExecuteCommandResponse{}, errors.New("failed to parse runner response")
	}
	return data, nil
}
