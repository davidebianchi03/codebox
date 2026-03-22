package analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
	"gitlab.com/codebox4073715/codebox/utils/analytics"
)

// HandleGetAnalyticsDataPreview godoc
// @Summary Get Analytics Data Preview
// @Schemes
// @Description Get Analytics Data Preview
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} analytics.AnalyticsRequestBody
// @Router /api/v1/admin/analytics-data-preview [get]
func HandleGetAnalyticsDataPreview(c *gin.Context) {
	data := analytics.GenerateAnalyticsData("preview")
	c.JSON(http.StatusOK, data)
}

// HandleGetAnalyticsConfig godoc
// @Summary Get Analytics Config
// @Schemes
// @Description Get Analytics Config
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AnalyticsConfigSerializer
// @Error 500
// @Router /api/v1/admin/analytics-config [get]
func HandleGetAnalyticsConfig(c *gin.Context) {
	conf, err := models.GetSingletonModelInstance[models.AnalyticsConfig]()
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to load analytics config",
		)
		return
	}

	c.JSON(http.StatusOK, serializers.LoadAnalyticsConfigSerializer(conf))
}

type UpdateAnalyticsConfigRequestBody struct {
	SendAnalyticsData bool `json:"send_analytics_data"`
}

// HandleUpdateAnalyticsConfig godoc
// @Summary Update Analytics Config
// @Schemes
// @Description Update Analytics Config
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body UpdateAnalyticsConfigRequestBody true "Request body"
// @Success 200 {object} serializers.AnalyticsConfigSerializer
// @Error 400
// @Error 500
// @Router /api/v1/admin/analytics-config [put]
func HandleUpdateAnalyticsConfig(c *gin.Context) {
	var req UpdateAnalyticsConfigRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	conf, err := models.GetSingletonModelInstance[models.AnalyticsConfig]()
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to load analytics config",
		)
		return
	}

	if conf == nil {
		conf = &models.AnalyticsConfig{
			SendAnalyticsData: req.SendAnalyticsData,
		}
	} else {
		conf.SendAnalyticsData = req.SendAnalyticsData
	}

	if err := models.SaveSingletonModel(conf); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to save analytics config",
		)
		return
	}

	c.JSON(http.StatusOK, serializers.LoadAnalyticsConfigSerializer(conf))
}

// HandleGetAnalyticsBannerSent godoc
// @Summary Get Analytics Banner Sent
// @Schemes
// @Description Get Analytics Banner Sent
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AnalyticsBannerSentSerializer
// @Error 500
// @Router /api/v1/admin/analytics-banner-sent [get]
func HandleGetAnalyticsBannerSent(c *gin.Context) {
	conf, err := models.GetSingletonModelInstance[models.AnalyticsConfig]()
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to load analytics config",
		)
		return
	}

	c.JSON(http.StatusOK, serializers.LoadAnalyticsBannerSentSerializer(conf))
}

// HandleUpdateAnalyticsBannerSent godoc
// @Summary Update Analytics Banner Sent
// @Schemes
// @Description Update Analytics Banner Sent
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AnalyticsBannerSentSerializer
// @Error 400
// @Error 500
// @Router /api/v1/admin/analytics-banner-sent [put]
func HandleUpdateAnalyticsBannerSent(c *gin.Context) {
	conf, err := models.GetSingletonModelInstance[models.AnalyticsConfig]()
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to load analytics config",
		)
		return
	}

	if conf == nil {
		conf = &models.AnalyticsConfig{
			AnalyticsBannerSent: true,
		}
	} else {
		conf.AnalyticsBannerSent = true
	}

	if err := models.SaveSingletonModel(conf); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to save analytics config",
		)
		return
	}

	c.JSON(http.StatusOK, serializers.LoadAnalyticsBannerSentSerializer(conf))
}
