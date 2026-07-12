package runnerinterface

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/codebox4073715/codebox/config"
)

/*
List all runner logs
*/
func (ri *RunnerInterface) ListRunnerLogs() (logs []RunnerLogRow, err error) {
	url := fmt.Sprintf("%s/api/v1/runner-logs/", ri.getRunnerBaseUrl())
	client := ri.getRequestsClient()

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return []RunnerLogRow{}, err
	}
	req.Header.Set(config.Environment.RunnerTokenHeader, ri.Runner.Token)

	res, err := client.Do(req)
	if err != nil {
		return []RunnerLogRow{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		if res.StatusCode == 404 {
			// old runners don't have this endpoint, so we just return an empty list
			return []RunnerLogRow{}, nil
		}
		return []RunnerLogRow{}, fmt.Errorf("receivedstatus %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []RunnerLogRow{}, err
	}

	var resp []RunnerLogRow
	if err = json.Unmarshal(body, &resp); err != nil {
		return []RunnerLogRow{}, err
	}

	return resp, nil
}
