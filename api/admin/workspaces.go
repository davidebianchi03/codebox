package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
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
