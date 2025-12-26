package runners

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
)

func HandleListRunnerTypes(c *gin.Context) {
	c.JSON(http.StatusOK, config.ListAvailableRunnerTypes())
}
