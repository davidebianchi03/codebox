package runnerinterface

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
)

func (ri *RunnerInterface) ContainerListDir(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
) ([]ContainerFileInfo, error) {
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
		return []ContainerFileInfo{}, fmt.Errorf("error from agent: %s", string(body))
	}

	var data []ContainerFileInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return []ContainerFileInfo{}, errors.New("failed to parse runner response")
	}

	return data, nil
}
