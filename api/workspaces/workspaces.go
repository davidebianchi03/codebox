package workspaces

import (
	"fmt"
	"net/http"
	"time"

	"codebox.com/api/utils"
	"codebox.com/bgtasks"
	"codebox.com/db"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
)

type WorkspaceWithoutDetailsResponse struct {
	Id             uint      `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	Type           string    `json:"type"`
	GitRepoUrl     string    `json:"git_repo_url"`
	CreatedAt      time.Time `json:"created_on"`
	LastActivityOn time.Time `json:"last_activity_on"`
	LastStartOn    time.Time `json:"last_start_on"`
}

func HandleListWorkspaces(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	var workspaces []db.Workspace
	result := db.DB.Where(map[string]interface{}{"owner_id": user.ID}).Find(&workspaces)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	responseObjs := make([]WorkspaceWithoutDetailsResponse, 0)
	for _, workspace := range workspaces {
		responseObjs = append(responseObjs, WorkspaceWithoutDetailsResponse{
			Id:             workspace.ID,
			Name:           workspace.Name,
			Status:         workspace.Status,
			Type:           workspace.Type,
			GitRepoUrl:     workspace.GitRepoUrl,
			CreatedAt:      workspace.CreatedAt,
			LastActivityOn: workspace.LastActivityOn,
			LastStartOn:    workspace.LastStartOn,
		})
	}

	ctx.JSON(http.StatusOK, responseObjs)
}

func HandleRetrieveWorkspace(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	id, found := ctx.Params.Get("id")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace db.Workspace
	result := db.DB.Where(map[string]interface{}{"ID": id, "owner_id": user.ID}).Find(&workspace)
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

	ctx.JSON(http.StatusOK, WorkspaceWithoutDetailsResponse{
		Id:             workspace.ID,
		Name:           workspace.Name,
		Status:         workspace.Status,
		Type:           workspace.Type,
		GitRepoUrl:     workspace.GitRepoUrl,
		CreatedAt:      workspace.CreatedAt,
		LastActivityOn: workspace.LastActivityOn,
		LastStartOn:    workspace.LastStartOn,
	})
}

func HandleCreateWorkspace(ctx *gin.Context) {
	type RequestBody struct {
		Name                       string `json:"name"`
		Type                       string `json:"type"`
		GitRepoUrl                 string `json:"git_repo_url"`
		GitRepoConfigurationFolder string `json:"git_repo_configuration_folder"`
	}

	var parsedBody RequestBody
	err := ctx.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	if parsedBody.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing field 'name'",
		})
		return
	}

	if parsedBody.Type == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing field 'type'",
		})
		return
	}

	if parsedBody.GitRepoUrl == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing field 'git_repo_url'",
		})
		return
	}

	if parsedBody.GitRepoConfigurationFolder == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing field 'git_repo_configuration_folder'",
		})
		return
	}

	if !db.IsItemInArray(parsedBody.Type, db.WorkspaceTypeChoices[:]) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": fmt.Sprintf("unsupported workspace type '%s'", parsedBody.Type),
		})
		return
	}

	owner, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// create workspace in db
	workspace := db.Workspace{
		Name:                       parsedBody.Name,
		Type:                       parsedBody.Type,
		Owner:                      owner,
		Status:                     db.WorkspaceStatusCreating,
		GitRepoUrl:                 parsedBody.GitRepoUrl,
		GitRepoConfigurationFolder: parsedBody.GitRepoConfigurationFolder,
		LastActivityOn:             time.Now(),
		LastStartOn:                time.Now(),
	}
	result := db.DB.Create(&workspace)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}
	workspace.AppendLogs("Creating workspace...")

	workspaceResponseObj := WorkspaceWithoutDetailsResponse{
		Id:             workspace.ID,
		Name:           workspace.Name,
		Status:         workspace.Status,
		Type:           workspace.Type,
		GitRepoUrl:     workspace.GitRepoUrl,
		CreatedAt:      workspace.CreatedAt,
		LastActivityOn: workspace.LastActivityOn,
		LastStartOn:    workspace.LastStartOn,
	}

	// start bg task
	bgtasks.BgTasksEnqueuer.Enqueue("start_workspace", work.Q{"workspace_id": workspace.ID})

	ctx.JSON(http.StatusOK, workspaceResponseObj)
}
