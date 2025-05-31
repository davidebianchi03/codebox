package templates

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/utils/targz"
)

// TemplatesList godoc
// @Summary List templates
// @Schemes
// @Description List all templates
// @Tags Templates
// @Accept json
// @Produce json
// @Success 200 {object} []models.WorkspaceTemplate
// @Router /api/v1/templates [get]
func HandleListTemplates(c *gin.Context) {
	var templates *[]models.WorkspaceTemplate
	if err := dbconn.DB.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// TemplatesRetrieve godoc
// @Summary Retrieve template by id
// @Schemes
// @Description Retrieve a template by id
// @Tags Templates
// @Param id path string true "Template ID"
// @Accept json
// @Produce json
// @Success 200 {object} models.WorkspaceTemplate
// @Router /api/v1/templates/:id [get]
func HandleRetrieveTemplate(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")

	// get object
	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	template, err := models.RetrieveWorkspaceTemplateByID(uint(ti))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// TemplatesRetrieve godoc
// @Summary Retrieve template by name
// @Schemes
// @Description Retrieve a template by name
// @Tags Templates
// @Param name path string true "Template name"
// @Accept json
// @Produce json
// @Success 200 {object} models.WorkspaceTemplate
// @Router /api/v1/templates-by-name/:ma,e [get]
func HandleRetrieveTemplateByName(c *gin.Context) {
	templateName, _ := c.Params.Get("templateName")

	template, err := models.RetrieveWorkspaceTemplateByName(templateName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

type CreateTemplateRequestBody struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// TemplateCreate godoc
// @Summary Create template
// @Schemes
// @Description Create a template
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body CreateTemplateRequestBody true "Template data"
// @Success 201 {object} []models.WorkspaceTemplate
// @Router /api/v1/templates [post]
func HandleCreateTemplate(c *gin.Context) {
	var requestBody *CreateTemplateRequestBody

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "missing or invalid argument",
		})
		return
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	// check if a row with the same name already exists
	wt, err := models.RetrieveWorkspaceTemplateByName(requestBody.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if wt != nil {
		c.JSON(http.StatusConflict, gin.H{
			"details": "another template with the same name already exists",
		})
		return
	}

	// check if type is valid
	var workspaceType *config.WorkspaceType
	for _, wt := range config.ListWorkspaceTypes() {
		// check if the current workspace type supports
		// templates as config source
		templatesSupported := false
		for _, configSource := range wt.SupportedConfigSources {
			if configSource == "template" {
				templatesSupported = true
			}
		}

		if templatesSupported && wt.ID == requestBody.Type {
			workspaceType = &wt
		}
	}

	if workspaceType == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "'type' is not valid",
		})
		return
	}

	// add template
	wt, err = models.CreateWorkspaceTemplate(
		requestBody.Name,
		requestBody.Type,
		requestBody.Description,
		requestBody.Icon,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	// create the first version
	tv, err := models.CreateTemplateVersion(
		*wt,
		fmt.Sprintf("version at %s", time.Now().Format("2006-01-02 15:04:05")),
		user,
		workspaceType.ConfigFilesDefaultPath,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, *wt)

	// add README.md to sources
	tgm := targz.TarGZManager{
		Filepath: tv.Sources.GetAbsolutePath(),
	}

	if !tv.Sources.Exists() {
		if err := tgm.CreateArchive(); err != nil {
			return
		}
	}

	configFilePath := tv.ConfigFilePath
	if !strings.HasPrefix(configFilePath, "./") {
		configFilePath = "./" + configFilePath
	}

	tgm.WriteFile("./README.md", []byte(fmt.Sprintf("# %s", wt.Name)))
	tgm.WriteFile(configFilePath, []byte(""))
}

type UpdateTemplateRequestBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// TemplateUpdate godoc
// @Summary Update template
// @Schemes
// @Description Update a template
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body UpdateTemplateRequestBody true "Template data"
// @Success 200 {object} []models.WorkspaceTemplate
// @Router /api/v1/templates/:templateId [put]
func HandleUpdateTemplate(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")

	// get object
	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	wt, err := models.RetrieveWorkspaceTemplateByID(uint(ti))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if wt == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	var requestBody *UpdateTemplateRequestBody

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "missing or invalid argument",
		})
		return
	}

	// check if a row with the same name already exists
	wte, err := models.RetrieveWorkspaceTemplateByName(requestBody.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if wte != nil {
		if wt.ID != uint(ti) {
			c.JSON(http.StatusConflict, gin.H{
				"details": "another template with the same name already exists",
			})
			return
		}
	}

	wt.Name = requestBody.Name
	wt.Description = requestBody.Description
	wt.Icon = requestBody.Icon

	if err := models.UpdateWorkspaceTemplate(*wt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, wt)
}

func HandleDeleteWorkspace(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")

	// get object
	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	wt, err := models.RetrieveWorkspaceTemplateByID(uint(ti))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if wt == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	workspaces, err := models.ListWorkspacesByTemplate(*wt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if len(workspaces) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"details": "cannot delete template, there are some workspace that are using it",
		})
		return
	}

	// remove versions
	wtv, err := models.ListWorkspaceTemplateVersionsByTemplate(*wt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	for _, tv := range *wtv {
		if err := models.DeleteTemplateVersion(tv); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	// delete template
	if err := models.DeleteTemplate(*wt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

func HandleListWorkspacesByTemplate(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")

	// get object
	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	wt, err := models.RetrieveWorkspaceTemplateByID(uint(ti))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if wt == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	workspaces, err := models.ListWorkspacesByTemplate(*wt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}
