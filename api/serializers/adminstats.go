package serializers

type AdminStatsSerializer struct {
	LoginCountsLast7Days []int64 `json:"login_counts_last_7_days"`
	TotalUsers           int64   `json:"total_users"`
	TotalRunners         int64   `json:"total_runners"`
	OnlineRunners        int64   `json:"online_runners"`
}

func LoadAdminStatsSerializer(
	loginCounts []int64,
	totalUsers,
	totalRunners,
	onlineRunners int64,
) *AdminStatsSerializer {
	return &AdminStatsSerializer{
		LoginCountsLast7Days: loginCounts,
		TotalUsers:           totalUsers,
		TotalRunners:         totalRunners,
		OnlineRunners:        onlineRunners,
	}
}
