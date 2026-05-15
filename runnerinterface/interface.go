package runnerinterface

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type RunnerInterface struct {
	Runner *models.Runner
}

func (ri *RunnerInterface) getRunnerBaseUrl() string {
	if ri.Runner.UsePublicUrl {
		return ri.Runner.PublicUrl
	}

	//TODO: raise exception if runner port is 0 (runner not connected)
	return fmt.Sprintf("http://127.0.0.1:%d", ri.Runner.Port)
}

func (ri *RunnerInterface) getRequestsClient() *http.Client {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return client
}
