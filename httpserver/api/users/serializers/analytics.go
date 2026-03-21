package serializers

import "gitlab.com/codebox4073715/codebox/db/models"

type AnalyticsConfigSerializer struct {
	SendAnalyticsData bool `json:"send_analytics_data"`
}

func LoadAnalyticsConfigSerializer(c *models.AnalyticsConfig) *AnalyticsConfigSerializer {
	return &AnalyticsConfigSerializer{
		SendAnalyticsData: c.SendAnalyticsData,
	}
}

type AnalyticsBannerSentSerializer struct {
	AnalyticsBannerSent bool `json:"analytics_banner_sent"`
}

func LoadAnalyticsBannerSentSerializer(c *models.AnalyticsConfig) *AnalyticsBannerSentSerializer {
	return &AnalyticsBannerSentSerializer{
		AnalyticsBannerSent: c.AnalyticsBannerSent,
	}
}
