package settings

import (
	"net/http"

	"github.com/davidebianchi03/codebox/config"
	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func HandleRetrieveServerSettings(c *gin.Context) {

	var usersCount int64
	if err := db.DB.Model(models.User{}).Count(&usersCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"use_gravatar":        config.Environment.UseGravatar,
		"use_subdomains":      config.Environment.UseSubDomains,
		"server_hostname":     c.Request.Host,
		"initial_user_exists": usersCount > 0,
	})
}
