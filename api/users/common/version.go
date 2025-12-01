package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/users/serializers"
)

// HandleRetrieveServerVersion godoc
// @Summary Retrieve the version of the server
// @Schemes
// @Description Retrieve the version of the server
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.VersionSerializer
// @Router /api/v1/version [get]
func HandleRetrieveServerVersion(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		serializers.GetVersionSerializedResponse(),
	)
}
