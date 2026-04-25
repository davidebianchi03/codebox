package runnerinterface

import (
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/proxy"
)

func (ri *RunnerInterface) PingAgent(container *models.WorkspaceContainer) bool {
	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/agent-version",
		ri.getRunnerBaseUrl(),
		container.WorkspaceID,
		container.ContainerName,
	)

	client := ri.getRequestsClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return false
	}

	req.Header.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)

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

func (ri *RunnerInterface) AgentForwardHttp(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	port *models.WorkspaceContainerPort,
	path string,
	rw http.ResponseWriter,
	req *http.Request,
) error {
	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/http-reverse-proxy?request_protocol=http&target_port=%d&target_path=%s&%s=%s",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
		port.PortNumber,
		url.QueryEscape(path),
		config.Environment.RunnerTokenQueryParam,
		url.QueryEscape(ri.Runner.Token),
	)

	proxyHeaders := http.Header{}
	proxy, err := proxy.CreateReverseProxy(url, 30, 30, true, proxyHeaders)
	if err != nil {
		return err
	}

	proxy.ServeHTTP(rw, req)
	return nil
}

func (ri *RunnerInterface) AgentForwardSsh(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	rw http.ResponseWriter,
	req *http.Request,
) error {
	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/ssh-proxy",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
	)

	proxyHeaders := http.Header{}
	proxyHeaders.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	proxy, err := proxy.CreateReverseProxy(url, 30, 30, true, proxyHeaders)
	if err != nil {
		return err
	}

	proxy.ServeHTTP(rw, req)
	return nil
}

func (ri *RunnerInterface) AgentForwardTerminal(
	workspace *models.Workspace,
	container *models.WorkspaceContainer,
	rw http.ResponseWriter,
	req *http.Request,
) error {
	url := fmt.Sprintf(
		"%s/api/v1/workspace/%d/container/%s/terminal?username=%s",
		ri.getRunnerBaseUrl(),
		workspace.ID,
		container.ContainerName,
		container.ContainerUserName,
	)

	proxyHeaders := http.Header{}
	proxyHeaders.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	proxy, err := proxy.CreateReverseProxy(url, 30, 30, true, proxyHeaders)
	if err != nil {
		return err
	}

	proxy.ServeHTTP(rw, req)
	return nil
}
