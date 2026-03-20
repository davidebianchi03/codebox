package bgtasks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/utils/analytics"
)

/*
Background task that sends analytics data to the analytics service,
this task is scheduled to run every 24 hours,
it checks if the user has allowed sending analytics data,
if allowed, it sends the data to the analytics service.
*/
func (jobContext *Context) SendAnalyticsData(job *work.Job) error {
	analyticsConfig, err := models.GetSingletonModelInstance[models.AnalyticsConfig]()
	if err != nil {
		// log the error and return nil to prevent retrying the task
		return nil
	}

	if analyticsConfig != nil {
		lastSuccessfullAttempt := time.UnixMilli(0)
		if analyticsConfig.LastSuccessfullAttempt != nil {
			lastSuccessfullAttempt = *analyticsConfig.LastSuccessfullAttempt
		}

		// if the last successfull attempt was less than 23 hours ago, skip sending analytics data
		if time.Since(lastSuccessfullAttempt) < 23*time.Hour {
			return nil
		}

		if analyticsConfig.SendAnalyticsData {
			analyticsData, err := json.Marshal(analytics.GenerateAnalyticsData(config.AnalyticsApiKey))
			if err != nil {
				// TODO: log the error
				return nil
			}

			client := &http.Client{}
			req, err := http.NewRequest(
				http.MethodPost,
				config.AnalyticsEndpoint,
				io.NopCloser(bytes.NewBuffer(analyticsData)),
			)
			if err != nil {
				// TODO: log the error
				return nil
			}

			res, err := client.Do(req)
			if err != nil {
				// TODO: log the error
				return nil
			}
			defer res.Body.Close()

			now := time.Now()
			if res.StatusCode < 200 || res.StatusCode > 299 {
				analyticsConfig.LastAttempt = &now
				models.SaveSingletonModel(analyticsConfig)
				// TODO: log the error
				return nil
			} else {
				analyticsConfig.LastSuccessfullAttempt = &now
				analyticsConfig.LastAttempt = &now
				models.SaveSingletonModel(analyticsConfig)
			}
		}
	}

	return nil
}
