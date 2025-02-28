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

// type WorkspaceWithoutDetailsResponse struct {
// 	Id             uint      `json:"id"`
// 	Name           string    `json:"name"`
// 	Status         string    `json:"status"`
// 	Type           string    `json:"type"`
// 	GitRepoUrl     string    `json:"git_repo_url"`
// 	CreatedAt      time.Time `json:"created_on"`
// 	LastActivityOn time.Time `json:"last_activity_on"`
// 	LastStartOn    time.Time `json:"last_start_on"`
// }

// type ForwardedPortDetails struct {
// 	PortNumber     uint   `json:"port_number"`
// 	Active         bool   `json:"active"`
// 	ConnectionType string `json:"connection_type"`
// 	Public         bool   `json:"public"`
// }

// type WorkspaceContainerDetails struct {
// 	Id                         int                    `json:"id"`
// 	Type                       string                 `json:"type"`
// 	Name                       string                 `json:"name"`
// 	ContainerUser              string                 `json:"container_user"`
// 	ContainerStatus            string                 `json:"container_status"`
// 	AgentStatus                string                 `json:"agent_status"`
// 	CanConnectRemoteDeveloping bool                   `json:"can_connect_remote_developing"`
// 	WorkspacePathInContainer   string                 `json:"workspace_path_in_container"`
// 	ForwardedPorts             []ForwardedPortDetails `json:"forwarded_ports"`
// }

// type WorkspaceWithDetailsResponse struct {
// 	Id             uint                        `json:"id"`
// 	Name           string                      `json:"name"`
// 	Status         string                      `json:"status"`
// 	Type           string                      `json:"type"`
// 	GitRepoUrl     string                      `json:"git_repo_url"`
// 	Containers     []WorkspaceContainerDetails `json:"containers"`
// 	CreatedAt      time.Time                   `json:"created_on"`
// 	LastActivityOn time.Time                   `json:"last_activity_on"`
// 	LastStartOn    time.Time                   `json:"last_start_on"`
// }

// /*
// GET api/v1/workspace
// */
// func HandleListWorkspaces(ctx *gin.Context) {
// 	user, err := utils.GetUserFromContext(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"detail": "internal server error",
// 		})
// 		return
// 	}

// 	var workspaces []models.Workspace
// 	result := db.DB.Where(map[string]interface{}{"owner_id": user.ID}).Find(&workspaces)
// 	if result.Error != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"detail": "internal server error",
// 		})
// 		return
// 	}

// 	responseObjs := make([]WorkspaceWithoutDetailsResponse, 0)
// 	for _, workspace := range workspaces {
// 		responseObjs = append(responseObjs, WorkspaceWithoutDetailsResponse{
// 			Id:             workspace.ID,
// 			Name:           workspace.Name,
// 			Status:         workspace.Status,
// 			Type:           workspace.Type,
// 			GitRepoUrl:     workspace.GitRepoUrl,
// 			CreatedAt:      workspace.CreatedAt,
// 			LastActivityOn: workspace.LastActivityOn,
// 			LastStartOn:    workspace.LastStartOn,
// 		})
// 	}

// 	ctx.JSON(http.StatusOK, responseObjs)
// }

// /*
// GET api/v1/workspace/:id
// */
// func HandleRetrieveWorkspace(ctx *gin.Context) {
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

// 	// retrieve workspace containers
// 	workspaceContainers := []models.WorkspaceContainer{}
// 	result = db.DB.Where(map[string]interface{}{"workspace_id": workspace.ID}).Preload("ForwardedPorts").Find(&workspaceContainers)
// 	if result.Error != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"detail": "internal server error",
// 		})
// 		return
// 	}

// 	responseObj := WorkspaceWithDetailsResponse{
// 		Id:             workspace.ID,
// 		Name:           workspace.Name,
// 		Status:         workspace.Status,
// 		Type:           workspace.Type,
// 		GitRepoUrl:     workspace.GitRepoUrl,
// 		CreatedAt:      workspace.CreatedAt,
// 		LastActivityOn: workspace.LastActivityOn,
// 		LastStartOn:    workspace.LastStartOn,
// 	}

// 	for _, container := range workspaceContainers {
// 		containerResponseObj := WorkspaceContainerDetails{
// 			Id:                         int(container.ID),
// 			Type:                       container.Type,
// 			Name:                       container.Name,
// 			ContainerUser:              container.ContainerUser,
// 			ContainerStatus:            container.ContainerStatus,
// 			AgentStatus:                container.AgentStatus,
// 			CanConnectRemoteDeveloping: container.CanConnectRemoteDeveloping,
// 			WorkspacePathInContainer:   container.WorkspacePathInContainer,
// 		}

// 		for _, forwardedPort := range container.ForwardedPorts {
// 			containerResponseObj.ForwardedPorts = append(containerResponseObj.ForwardedPorts, ForwardedPortDetails{
// 				PortNumber:     forwardedPort.PortNumber,
// 				Active:         forwardedPort.Active,
// 				ConnectionType: forwardedPort.ConnectionType,
// 				Public:         forwardedPort.Public,
// 			})
// 		}

// 		responseObj.Containers = append(responseObj.Containers, containerResponseObj)
// 	}

// 	ctx.JSON(http.StatusOK, responseObj)
// }

/*
POST api/v1/workspace
*/
func HandleCreateWorkspace(c *gin.Context) {
	type RequestBody struct {
		Name                       string   `json:"name" binding:"required"`
		Type                       string   `json:"type" binding:"required"`
		RunnerID                   uint     `json:"runner_id" binding:"required"`
		ConfigSource               string   `json:"config_source" binding:"required"`
		TemplateVersionID          uint     `json:"template_version_id"`
		GitRepoUrl                 string   `json:"git_repo_url"`
		GitRepoConfigurationFolder string   `json:"git_repo_configuration_folder"`
		EnvironmentVariables       []string `json:"environment_variables" binding:"required"`
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
		if parsedBody.GitRepoConfigurationFolder == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "missing param 'git_repo_configuration_folder",
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
		Status:               models.WorkspaceStatusCreating,
		Type:                 parsedBody.Type,
		Runner:               *runner,
		ConfigSource:         parsedBody.ConfigSource,
		TemplateVersion:      templateVersion,
		GitSource:            gitSource,
		ConfigSourceFilePath: parsedBody.GitRepoConfigurationFolder,
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

	// check if type is a valid option

	// workspace := models.Workspace{
	// 	Name: parsedBody.Name,
	// }

	// var parsedBody RequestBody
	// err := ctx.ShouldBindBodyWithJSON(&parsedBody)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"detail": err.Error(),
	// 	})
	// 	return
	// }

	// if parsedBody.Name == "" {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"detail": "missing field 'name'",
	// 	})
	// 	return
	// }

	// if parsedBody.Type == "" {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"detail": "missing field 'type'",
	// 	})
	// 	return
	// }

	// if parsedBody.GitRepoUrl == "" {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"detail": "missing field 'git_repo_url'",
	// 	})
	// 	return
	// }

	// if parsedBody.GitRepoConfigurationFolder == "" {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"detail": "missing field 'git_repo_configuration_folder'",
	// 	})
	// 	return
	// }

	// if !db.IsItemInArray(parsedBody.Type, db.WorkspaceTypeChoices[:]) {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"detail": fmt.Sprintf("unsupported workspace type '%s'", parsedBody.Type),
	// 	})
	// 	return
	// }

	// owner, err := utils.GetUserFromContext(ctx)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"detail": "internal server error",
	// 	})
	// 	return
	// }

	// // create workspace in db
	// workspace := models.Workspace{
	// 	Name:                       parsedBody.Name,
	// 	Type:                       parsedBody.Type,
	// 	User:                       owner,
	// 	Status:                     models.WorkspaceStatusCreating,
	// 	GitRepoUrl:                 parsedBody.GitRepoUrl,
	// 	GitRepoConfigurationFolder: parsedBody.GitRepoConfigurationFolder,
	// 	LastActivityOn:             time.Now(),
	// 	LastStartOn:                time.Now(),
	// }
	// result := db.DB.Create(&workspace)
	// if result.Error != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"detail": "internal server error",
	// 	})
	// 	return
	// }
	// workspace.AppendLogs("Creating workspace...")

	// // workspaceResponseObj := WorkspaceWithoutDetailsResponse{
	// // 	Id:             workspace.ID,
	// // 	Name:           workspace.Name,
	// // 	Status:         workspace.Status,
	// // 	Type:           workspace.Type,
	// // 	GitRepoUrl:     workspace.GitRepoUrl,
	// // 	CreatedAt:      workspace.CreatedAt,
	// // 	LastActivityOn: workspace.LastActivityOn,
	// // 	LastStartOn:    workspace.LastStartOn,
	// // }

	// // start bg task
	// bgtasks.BgTasksEnqueuer.Enqueue("start_workspace", work.Q{"workspace_id": workspace.ID})

	// ctx.JSON(http.StatusCreated, workspaceResponseObj)
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
