package workspaces

// func HandleRetrieveWorkspaceLogs(ctx *gin.Context) {
// 	user, err := utils.GetUserFromContext(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"detail": "internal server error",
// 		})
// 		return
// 	}

// 	id, found := ctx.Params.Get("workspaceId")
// 	if !found {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"detail": "workspace not found",
// 		})
// 		return
// 	}

// 	var workspace models.Workspace
// 	result := db.DB.Where(map[string]interface{}{"ID": id, "owner_id": user.ID}).Find(&workspace)
// 	if result.Error != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"detail": "internal server error",
// 		})
// 		return
// 	}

// 	if result.RowsAffected == 0 {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"detail": "workspace not found",
// 		})
// 		return
// 	}

// 	logs, _ := workspace.RetrieveLogs()
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"logs": logs,
// 	})
// }
