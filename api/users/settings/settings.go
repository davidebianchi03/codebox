package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
)

// , this api is available only to administrator
func HandleRetrieveServerSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": config.ServerVersion,
	})
}
