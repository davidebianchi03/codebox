package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
)

func HandleRetrieveServerSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":         config.ServerVersion,
		"use_subdomains":  config.Environment.UseSubDomains,
		"external_url":    config.Environment.ExternalUrl,
		"wildcard_domain": config.Environment.WildcardDomain,
	})
}
