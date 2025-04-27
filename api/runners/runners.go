package runners

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

func HandleListRunners(c *gin.Context) {
	var runners []models.Runner
	r := dbconn.DB.Find(&runners)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, runners)
}
