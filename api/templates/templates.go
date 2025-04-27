package templates

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	var template *models.WorkspaceTemplate
	_ = template
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
// @Param request body templates.CreateTemplateRequestBody
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
	Type        string `json:"type" binding:"required"`
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
// @Param request body templates.UpdateTemplateRequestBody
// @Success 200 {object} []models.WorkspaceTemplate
// @Router /api/v1/templates/:templateId [put]
func HandleUpdateTemplate(c *gin.Context) {
	// templateId, err := c.Params.Get("templateId")

	// // get object
	// err, wt :=

	// var requestBody *UpdateTemplateRequestBody

	// if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"details": "missing or invalid argument",
	// 	})
	// 	return
	// }

	// // check if a row with the same name already exists
	// wt, err := models.RetrieveWorkspaceTemplateByName(requestBody.Name)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"details": "internal server error",
	// 	})
	// 	return
	// }

	// if wt != nil {
	// 	c.JSON(http.StatusConflict, gin.H{
	// 		"details": "another template with the same name already exists",
	// 	})
	// 	return
	// }

}
