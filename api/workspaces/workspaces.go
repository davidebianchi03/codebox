package workspaces

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/api/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// HandleListWorkspaces godoc
// @Summary List workspaces
// @Schemes
// @Description List workspaces created by the current user
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} []serializers.WorkspaceSerializer
// @Router /api/v1/workspace [get]
func HandleListWorkspaces(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspaces, err := models.ListUserWorkspaces(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, serializers.LoadMultipleWorkspaceSerializer(workspaces))
}

// HandleRetrieveWorkspace godoc
// @Summary Retrieve workspace by id
// @Schemes
// @Description Retrieve a workspace by id
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} serializers.WorkspaceSerializer
// @Router /api/v1/workspace/:id [get]
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

	workspaceId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspce not found",
		})
		return
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, uint(workspaceId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if workspace == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, *serializers.LoadWorkspaceSerializer(workspace))
}

type CreateWorkspaceRequestBody struct {
	Name                 string   `json:"name" binding:"required"`
	Type                 string   `json:"type" binding:"required"`
	RunnerID             uint     `json:"runner_id" binding:"required"`
	ConfigSource         string   `json:"config_source" binding:"required"`
	TemplateVersionID    uint     `json:"template_version_id"`
	GitRepoUrl           string   `json:"git_repo_url"`
	GitRefName           string   `json:"git_ref_name"`
	ConfigSourceFilePath string   `json:"config_source_path"`
	EnvironmentVariables []string `json:"environment_variables" binding:"required"`
}

// HandleRetrieveWorkspace godoc
// @Summary Create a workspace
// @Schemes
// @Description Create a new workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param request body CreateWorkspaceRequestBody true "Data for creating a workspace"
// @Success 201 {object} serializers.WorkspaceSerializer
// @Router /api/v1/workspace [post]
func HandleCreateWorkspace(c *gin.Context) {
	var requestBody *CreateWorkspaceRequestBody
	err := c.ShouldBindBodyWithJSON(&requestBody)
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

	wt := config.RetrieveWorkspaceType(requestBody.Type)
	if wt == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid workspace type",
		})
		return
	}

	// validate runner
	runner, err := models.RetrieveRunnerByID(requestBody.RunnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if runner == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "runner matching runner_id and type not found",
		})
		return
	}

	rt := config.RetrieveRunnerTypeByID(runner.Type)
	if rt == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "runner matching runner_id and type not found",
		})
		return
	}

	// check if the runner supports the requested workspace type
	supported := false
	for _, supportedType := range rt.SupportedTypes {
		if supportedType.ID == wt.ID {
			supported = true
			break
		}
	}

	if !supported {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "runner does not support the requested workspace type",
		})
		return
	}

	// TODO: check if user is allowed to use requested runner
	// validate workspace configuration source
	var gitSource *models.GitWorkspaceSource
	var templateVersion *models.WorkspaceTemplateVersion
	if requestBody.ConfigSource == models.WorkspaceConfigSourceGit {
		if requestBody.GitRepoUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "missing param 'git_repo_url",
			})
			return
		}
		if requestBody.ConfigSourceFilePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "missing param 'config_source_path",
			})
			return
		}

		gitSource, err = models.CreateGitWorkspaceSource(
			requestBody.GitRepoUrl,
			requestBody.GitRefName,
			requestBody.ConfigSourceFilePath,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}
	} else if requestBody.ConfigSource == models.WorkspaceConfigSourceTemplate {
		templateVersion, err = models.RetrieveWorkspaceTemplateVersionsById(requestBody.TemplateVersionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}

		if templateVersion == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "requested template version does not exist",
			})
			return
		}

		if templateVersion.Template.Type != requestBody.Type {
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

	workspace, err := models.CreateWorkspace(
		requestBody.Name,
		&currentUser,
		requestBody.Type,
		runner,
		requestBody.ConfigSource,
		templateVersion,
		gitSource,
		requestBody.EnvironmentVariables,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspace.AppendLogs("Creating workspace...")
	bgtasks.BgTasksEnqueuer.Enqueue("start_workspace", work.Q{"workspace_id": workspace.ID})

	c.JSON(http.StatusCreated, workspace)
}

// HandleStopWorkspace godoc
// @Summary Stop a workspace
// @Schemes
// @Description Stop a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} serializers.WorkspaceSerializer
// @Router /api/v1/workspace/:id/stop [post]
func HandleStopWorkspace(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	id, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if workspace == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	if workspace.Status == models.WorkspaceStatusStopping || workspace.Status == models.WorkspaceStatusStopped {
		ctx.JSON(http.StatusConflict, gin.H{
			"detail": "workspace is already stopped",
		})
		return
	}

	workspace, err = models.UpdateWorkspace(
		workspace,
		workspace.Name,
		models.WorkspaceStatusStopping,
		workspace.Runner,
		workspace.ConfigSource,
		workspace.TemplateVersion,
		workspace.GitSource,
		workspace.EnvironmentVariables,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspace.ClearLogs()

	// start bg task
	bgtasks.BgTasksEnqueuer.Enqueue("stop_workspace", work.Q{"workspace_id": workspace.ID})

	ctx.JSON(http.StatusOK, serializers.LoadWorkspaceSerializer(workspace))
}

// HandleStartWorkspace godoc
// @Summary Start a workspace
// @Schemes
// @Description Start a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} serializers.WorkspaceSerializer
// @Router /api/v1/workspace/:id/start [post]
func HandleStartWorkspace(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	id, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if workspace == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	if workspace.Status == models.WorkspaceStatusStarting || workspace.Status == models.WorkspaceStatusRunning {
		ctx.JSON(http.StatusConflict, gin.H{
			"detail": "workspace is already running",
		})
		return
	}

	workspace, err = models.UpdateWorkspace(
		workspace,
		workspace.Name,
		models.WorkspaceStatusStarting,
		workspace.Runner,
		workspace.ConfigSource,
		workspace.TemplateVersion,
		workspace.GitSource,
		workspace.EnvironmentVariables,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspace.ClearLogs()

	// start bg task
	bgtasks.BgTasksEnqueuer.Enqueue("start_workspace", work.Q{"workspace_id": workspace.ID})

	ctx.JSON(http.StatusOK, serializers.LoadWorkspaceSerializer(workspace))
}

type UpdateWorkspaceRequestBody struct {
	GitRepoUrl           *string   `json:"git_repo_url"`
	GitRefName           *string   `json:"git_ref_name"`
	ConfigSourcePath     *string   `json:"config_source_path"`
	EnvironmentVariables *[]string `json:"environment_variables"`
}

// HandleUpdateWorkspace godoc
// @Summary Update a workspace
// @Schemes
// @Description Update a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param request body CreateWorkspaceRequestBody true "Data to update a workspace"
// @Success 200 {object} serializers.WorkspaceSerializer
// @Router /api/v1/workspace/:id [put]
func HandleUpdateWorkspace(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}
	id, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if workspace == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	if workspace.Status != models.WorkspaceStatusStopped {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"detail": "cannot update, workspace is running",
		})
		return
	}

	var reqBody UpdateWorkspaceRequestBody
	if err := ctx.ShouldBindBodyWithJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"detail": "missing or invalid request argument",
		})
		return
	}

	if workspace.ConfigSource == models.WorkspaceConfigSourceGit {
		gitSource, err := models.UpdateGitWorkspaceSource(
			workspace.GitSource,
			*reqBody.GitRepoUrl,
			*reqBody.GitRefName,
			*reqBody.ConfigSourcePath,
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}

		workspace.GitSource = gitSource
	}

	if reqBody.EnvironmentVariables != nil {
		workspace, err = models.UpdateWorkspace(
			workspace,
			workspace.Name,
			workspace.Status,
			workspace.Runner,
			workspace.ConfigSource,
			workspace.TemplateVersion,
			workspace.GitSource,
			*reqBody.EnvironmentVariables,
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, workspace)
}

// HandleDeleteWorkspace godoc
// @Summary Delete a workspace
// @Schemes
// @Description Delete a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 204
// @Router /api/v1/workspace/:id [delete]
func HandleDeleteWorkspace(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	id, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if workspace == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	skipErrors := false
	if strings.ToLower(ctx.Request.URL.Query().Get("skip_errors")) == "true" {
		skipErrors = true
	}

	workspace, err = models.UpdateWorkspace(
		workspace,
		workspace.Name,
		models.WorkspaceStatusDeleting,
		workspace.Runner,
		workspace.ConfigSource,
		workspace.TemplateVersion,
		workspace.GitSource,
		workspace.EnvironmentVariables,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspace.ClearLogs()

	// start bg task
	bgtasks.BgTasksEnqueuer.Enqueue(
		"delete_workspace",
		work.Q{"workspace_id": workspace.ID, "skip_errors": skipErrors},
	)

	ctx.JSON(http.StatusOK, gin.H{
		"detail": "deleting workspace...",
	})
}

// HandleDeleteWorkspace godoc
// @Summary Update workspace configuration
// @Schemes
// @Description Update workspace configuration, retrieving the configuration files from the git repository or template
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/workspace/:workspaceId/update-config [post]
func HandleUpdateWorkspaceConfiguration(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	id, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if workspace == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	if workspace.Status != models.WorkspaceStatusStopped {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"detail": "cannot update, workspace is running",
		})
		return
	}

	workspace, err = models.UpdateWorkspace(
		workspace,
		workspace.Name,
		models.WorkspaceStatusStarting,
		workspace.Runner,
		workspace.ConfigSource,
		workspace.TemplateVersion,
		workspace.GitSource,
		workspace.EnvironmentVariables,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspace.ClearLogs()

	workspace.AppendLogs("Updating workspace configuration sources...")
	bgtasks.BgTasksEnqueuer.Enqueue("update_workspace_config", work.Q{"workspace_id": workspace.ID})

	ctx.JSON(http.StatusOK, gin.H{
		"details": "starting workspace",
	})
}
