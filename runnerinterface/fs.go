package runnerinterface

import (
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
)

func (ri *RunnerInterface) ContainerListDir(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	path string,
) error {
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
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
