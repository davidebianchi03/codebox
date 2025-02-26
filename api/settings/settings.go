package settings

import (
	"net/http"

	"github.com/davidebianchi03/codebox/config"
	"github.com/gin-gonic/gin"
)

func HandleRetrieveServerSettings(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"use_gravatar":    config.Environment.UseGravatar,
		"use_subdomains":  config.Environment.UseSubDomains,
		"server_hostname": ctx.Request.Host,
	})
}
