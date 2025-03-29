package admin

import (
	"net/http"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func HandleAdminListUsers(c *gin.Context) {
	var users *[]models.User

	if db.DB.Find(&users).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
