package runners

import (
	"net/http"

	"github.com/davidebianchi03/codebox/config"
	"github.com/gin-gonic/gin"
)

func HandleListRunnerTypes(c *gin.Context) {
	c.JSON(http.StatusOK, config.ListAvailableRunnerTypes())
}
