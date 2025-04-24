package runners

import (
	"net/http"

	dbconn "github.com/davidebianchi03/codebox/db/connection"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
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
