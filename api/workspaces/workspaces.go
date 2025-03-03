package workspaces

import (
	"net/http"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/bgtasks"
	"github.com/davidebianchi03/codebox/config"
	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
)

/*
GET api/v1/workspace
*/
func HandleListWorkspaces(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	var workspaces []models.Workspace
	result := db.DB.Find(&workspaces, map[string]interface{}{"user_id": user.ID})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, workspaces)
}

/*
GET api/v1/workspace/:id
*/
func HandleRetrieveWorkspace(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	id, found := ctx.Params.Get("workspaceId")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace models.Workspace
	result := db.DB.Where(map[string]interface{}{"ID": id, "user_id": user.ID}).Find(&workspace)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, workspace)
}

/*
POST api/v1/workspace
*/
func HandleCreateWorkspace(c *gin.Context) {
	type RequestBody struct {
		Name                 string   `json:"name" binding:"required"`
		Type                 string   `json:"type" binding:"required"`
		RunnerID             uint     `json:"runner_id" binding:"required"`
		ConfigSource         string   `json:"config_source" binding:"required"`
		TemplateVersionID    uint     `json:"template_version_id"`
		GitRepoUrl           string   `json:"git_repo_url"`
		ConfigSourceFilePath string   `json:"config_source_path"`
		EnvironmentVariables []string `json:"environment_variables" binding:"required"`
	}

	var parsedBody RequestBody
	err := c.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	currentUser, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspaceTypeFound := false
	for _, t := range config.ListWorkspaceTypes() {
		if t.ID == parsedBody.Type {
			workspaceTypeFound = true
			break
		}
	}

	if !workspaceTypeFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid workspace type",
		})
		return
	}

	// validate runner
	runner := &models.Runner{}
	db.DB.First(&runner, map[string]interface{}{
		"ID":   parsedBody.RunnerID,
		"Type": parsedBody.Type,
	})

	if runner.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "runner matching runner_id and type not found",
		})
		return
	}
	// TODO: check if user is allowed to use requested runner

	// validate workspace configuration source
	var gitSource *models.GitWorkspaceSource
	var templateVersion *models.WorkspaceTemplateVersion
	if parsedBody.ConfigSource == models.WorkspaceConfigSourceGit {
		if parsedBody.GitRepoUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "missing param 'git_repo_url",
			})
			return
		}
		if parsedBody.ConfigSourceFilePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "missing param 'config_source_path",
			})
			return
		}

		gitSource = &models.GitWorkspaceSource{
			RepositoryURL: parsedBody.GitRepoUrl,
		}

		r := db.DB.Create(gitSource)
		if r.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}
	} else if parsedBody.ConfigSource == models.WorkspaceConfigSourceTemplate {
		db.DB.First(templateVersion, map[string]interface{}{
			"ID": parsedBody.TemplateVersionID,
		})

		if templateVersion == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "requested template version does not exist",
			})
			return
		}

		if templateVersion.Template.Type != parsedBody.Type {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "requested template version does not exist",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid value for 'config_source'",
		})
		return
	}

	// create new database row
	workspace := models.Workspace{
		Name:                 parsedBody.Name,
		User:                 currentUser,
		Status:               models.WorkspaceStatusStarting,
		Type:                 parsedBody.Type,
		Runner:               runner,
		ConfigSource:         parsedBody.ConfigSource,
		TemplateVersion:      templateVersion,
		GitSource:            gitSource,
		ConfigSourceFilePath: parsedBody.ConfigSourceFilePath,
		EnvironmentVariables: parsedBody.EnvironmentVariables,
	}

	r := db.DB.Create(&workspace)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspace.AppendLogs("Creating workspace...")
	bgtasks.BgTasksEnqueuer.Enqueue("start_workspace", work.Q{"workspace_id": workspace.ID})

	c.JSON(http.StatusCreated, workspace)
}

// /*
// POST api/v1/workspace/:id/stop
// */
// func HandleStopWorkspace(ctx *gin.Context) {
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

// 	var workspace db.Workspace
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

// 	if workspace.Status == db.WorkspaceStatusStopping || workspace.Status == db.WorkspaceStatusStopped {
// 		ctx.JSON(http.StatusConflict, gin.H{
// 			"detail": "workspace is already stopped",
// 		})
// 		return
// 	}

// 	workspace.Status = db.WorkspaceStatusStopping
// 	db.DB.Save(&workspace)

// 	// start bg task
// 	bgtasks.BgTasksEnqueuer.Enqueue("stop_workspace", work.Q{"workspace_id": workspace.ID})

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"detail": "stopping workspace...",
// 	})
// }

// /*
// POST api/v1/workspace/:id/start
// */
// func HandleStartWorkspace(ctx *gin.Context) {
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

// 	var workspace db.Workspace
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

// 	if workspace.Status == db.WorkspaceStatusCreating || workspace.Status == db.WorkspaceStatusStarting || workspace.Status == db.WorkspaceStatusRunning {
// 		ctx.JSON(http.StatusConflict, gin.H{
// 			"detail": "workspace is already running",
// 		})
// 		return
// 	}

// 	workspace.Status = db.WorkspaceStatusStarting
// 	db.DB.Save(&workspace)

// 	// start bg task
// 	bgtasks.BgTasksEnqueuer.Enqueue("start_workspace", work.Q{"workspace_id": workspace.ID})

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"detail": "starting workspace...",
// 	})
// }

// /*
// DELETE api/v1/workspace/:id
// */
// func HandleDeleteWorkspace(ctx *gin.Context) {
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

// 	var workspace db.Workspace
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

// 	workspace.Status = db.WorkspaceStatusDeleting
// 	db.DB.Save(&workspace)

// 	// start bg task
// 	bgtasks.BgTasksEnqueuer.Enqueue("delete_workspace", work.Q{"workspace_id": workspace.ID})

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"detail": "deleting workspace...",
// 	})
// }
