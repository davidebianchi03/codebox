package settings

import (
	"net/http"

	"github.com/davidebianchi03/codebox/config"
	"github.com/gin-gonic/gin"
)

func HandleRetrieveServerSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":         config.ServerVersion,
		"use_gravatar":    config.Environment.UseGravatar,
		"use_subdomains":  config.Environment.UseSubDomains,
		"external_url":    config.Environment.ExternalUrl,
		"wildcard_domain": config.Environment.WildcardDomain,
	})
}
