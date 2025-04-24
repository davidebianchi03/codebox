package admin

import (
	"net/http"

	dbconn "github.com/davidebianchi03/codebox/db/connection"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func AdminListWorkspaces(ctx *gin.Context) {
	var workspaces []models.Workspace
	if err := dbconn.DB.Preload("User").Preload("Runner").Find(&workspaces).Order("UpdatedAt ASC").Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, workspaces)
}
