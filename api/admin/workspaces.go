package admin

import (
	"net/http"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func AdminListWorkspaces(ctx *gin.Context) {
	var workspaces []models.Workspace
	if err := db.DB.Preload("User").Preload("Runner").Find(&workspaces).Order("UpdatedAt ASC").Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, workspaces)
}
