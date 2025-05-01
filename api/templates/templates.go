package templates

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
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
	valid := false
	for _, workspaceType := range config.ListWorkspaceTypes() {
		// check if the current workspace type supports
		// templates as config source
		templatesSupported := false
		for _, configSource := range workspaceType.SupportedConfigSources {
			if configSource == "template" {
				templatesSupported = true
			}
		}

		if templatesSupported && workspaceType.ID == requestBody.Type {
			valid = true
		}
	}

	if !valid {
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

	c.JSON(http.StatusCreated, *wt)
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
