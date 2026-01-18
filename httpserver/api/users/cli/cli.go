package cli

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

// HandleRetrieveCLIVersion godoc
// @Summary Retrieve recommended cli version
// @Schemes
// @Description Retrieve recommended cli version
// @Tags CLI
// @Accept json
// @Produce json
// @Success 200 {object} serializers.CLIVersionSerializer
// @Router /api/v1/cli-version [get]
func HandleRetrieveCLIVersion(c *gin.Context) {
	c.JSON(http.StatusOK, serializers.CLIVersionSerializer{
		Version: config.RecommendedCLIVersion,
	})
}

// List cli godoc
// @Summary List cli
// @Schemes
// @Description List cli
// @Tags CLI
// @Accept json
// @Produce json
// @Success 200 {object} []cli.CLIBuild
// @Router /api/v1/cli [get]
func HandleListCLI(c *gin.Context) {
	c.JSON(http.StatusOK, CliBuilds)
}

// Retrieve cli godoc
// @Summary Retrieve cli by its id
// @Schemes
// @Description Retrieve cli by its id
// @Tags CLI
// @Accept json
// @Produce json
// @Success 200 {object} cli.CLIBuild
// @Router /api/v1/cli/:id [get]
func HandleRetrieveCLI(c *gin.Context) {
	buildId, _ := c.Params.Get("id")

	for _, build := range CliBuilds {
		if build.Id == buildId {
			c.JSON(http.StatusOK, build)
			return
		}
	}

	utils.ErrorResponse(c, http.StatusNotFound, "cli build not found")
}

// Download cli godoc
// @Summary Download cli
// @Schemes
// @Description Download cli
// @Tags CLI
// @Accept json
// @Produce json
// @Success 200 {file} binary
// @Router /api/v1/cli/:id/download [get]
func HandleDownloadCLI(c *gin.Context) {
	buildId, _ := c.Params.Get("id")

	var build *CLIBuild

	for _, b := range CliBuilds {
		if b.Id == buildId {
			build = &b
			break
		}
	}

	if build == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "cli build not found")
		return
	}

	c.Status(http.StatusOK)
	c.FileAttachment(
		path.Join(config.Environment.CliBinariesPath, build.File),
		build.File,
	)
}
