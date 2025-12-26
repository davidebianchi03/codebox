package serializers

type AdminStatsSerializer struct {
	LoginCountsLast7Days []int64 `json:"login_counts_last_7_days"`
	TotalUsers           int64   `json:"total_users"`
	OnlineRunners        int64   `json:"online_runners"`
	OnlineWorkspaces     int64   `json:"online_workspaces"`
}

func LoadAdminStatsSerializer(
	loginCounts []int64,
	totalUsers,
	onlineRunners,
	onlineWorkspaces int64,
) *AdminStatsSerializer {
	return &AdminStatsSerializer{
		LoginCountsLast7Days: loginCounts,
		TotalUsers:           totalUsers,
		OnlineRunners:        onlineRunners,
		OnlineWorkspaces:     onlineWorkspaces,
	}
}
