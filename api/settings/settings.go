package settings

import (
	"net/http"

	"codebox.com/env"
	"github.com/gin-gonic/gin"
)

func HandleRetrieveServerSettings(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"use_gravatar": env.CodeBoxEnv.UseGravatar,
	})
}
