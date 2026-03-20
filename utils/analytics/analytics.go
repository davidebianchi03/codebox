package analytics

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
)

type AnalyticsData struct {
	DistinctId        string `json:"distinct_id"`
	LicenseType       string `json:"license_type"`
	ServerVersion     string `json:"server_version"`
	TotalUsers        int64  `json:"total_users"`
	ApprovedUsers     int64  `json:"approved_users"`
	TotalRunners      int64  `json:"total_runners"`
	OnlineRunners     int64  `json:"online_runners"`
	TotalWorkspaces   int64  `json:"total_workspaces"`
	RunningWorkspaces int64  `json:"running_workspaces"`
	TotalTemplates    int64  `json:"total_templates"`
}

type AnalyticsRequestBody struct {
	APIKey     string        `json:"api_key"`
	Event      string        `json:"event"`
	Properties AnalyticsData `json:"properties"`
	Timestamp  string        `json:"timestamp"`
}

func GenerateAnalyticsData(apiKey string) AnalyticsRequestBody {
	h := sha256.New()
	h.Write([]byte(config.Environment.ExternalUrl))
	instanceId := hex.EncodeToString(h.Sum(nil))

	totalUsers, _ := models.CountAllUsers()
	approvedUsers, _ := models.CountApprovedUsers()
	totalRunners, _ := models.CountAllRunners()
	onlineRunners, _ := models.CountOnlineRunners()
	totalWorkspaces, _ := models.CountAllWorkspaces()
	runningWorkspaces, _ := models.CountAllOnlineWorkspaces()
	totalTemplates, _ := models.CountAllTemplates()

	return AnalyticsRequestBody{
		APIKey: apiKey,
		Event:  "codebox-analytics",
		Properties: AnalyticsData{
			DistinctId:        instanceId,
			LicenseType:       "community",
			ServerVersion:     config.ServerVersion,
			TotalUsers:        totalUsers,
			ApprovedUsers:     approvedUsers,
			TotalRunners:      totalRunners,
			OnlineRunners:     onlineRunners,
			TotalWorkspaces:   totalWorkspaces,
			RunningWorkspaces: runningWorkspaces,
			TotalTemplates:    totalTemplates,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
