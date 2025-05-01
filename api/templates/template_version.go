package templates

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// TemplateVersionByTemplateList godoc
// @Summary List template versions by template
// @Schemes
// @Description List all template versions by template
// @Tags Templates
// @Accept json
// @Produce json
// @Success 200 {object} []models.WorkspaceTemplateVersion
// @Router /api/v1/templates/:templateId/versions [get]
func HandleListTemplateVersionsByTemplate(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")

	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
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

	tv, err := models.ListWorkspaceTemplateVersionsByTemplate(*wt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, tv)
}

// RetrieveTemplateVersionByTemplate godoc
// @Summary Retrieve template version by id
// @Schemes
// @Description Retrieve template version by id
// @Tags Templates
// @Accept json
// @Produce json
// @Success 200 {object} models.WorkspaceTemplateVersion
// @Router /api/v1/templates/:templateId/versions/:versionId [get]
func HandleRetrieveTemplateVersionByTemplate(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")
	templateVersionId, _ := c.Params.Get("versionId")

	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"details": "template not found",
		})
		return
	}

	tvi, err := strconv.Atoi(templateVersionId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"details": "template version not found",
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

	tv, err := models.RetrieveWorkspaceTemplateVersionsByIdByTemplate(*wt, uint(tvi))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if tv == nil {
		c.JSON(http.StatusOK, gin.H{
			"details": "template version not found",
		})
		return
	}

	c.JSON(http.StatusOK, tv)
}

type CreateTemplateVersionRequestBody struct {
	Name string `json:"name" binding:"required,min=1"`
}

// CreateTemplateVersionByTemplate godoc
// @Summary Create new template version
// @Schemes
// @Description Create new template version
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body CreateTemplateVersionRequestBody true "Template version data"
// @Success 201 {object} models.WorkspaceTemplateVersion
// @Router /api/v1/templates/:templateId/versions [post]
func HandleCreateTemplateVersionByTemplate(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")

	var requestBody *CreateTemplateVersionRequestBody
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "missing or invalid request param",
		})
		return
	}

	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
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
		c.JSON(http.StatusOK, gin.H{
			"details": "template not found",
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

	wtb, err := models.CreateTemplateVersion(*wt, "Dfdf", user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, wtb)
}

type UpdateTemplateVersionRequestBody struct {
	Name      string `json:"name" binding:"required,min=1"`
	Published bool   `json:"published"`
}

// UpdateTemplateversionByTemplate godoc
// @Summary Update a template version
// @Schemes
// @Description Update a template version
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body UpdateTemplateVersionRequestBody true "Template version data"
// @Success 200 {object} models.WorkspaceTemplateVersion
// @Router /api/v1/templates/:templateId/versions/:versionId [put]
func HandleUpdateTemplateversionByTemplate(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")
	templateVersionId, _ := c.Params.Get("versionId")

	var requestBody *UpdateTemplateVersionRequestBody
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "missing or invalid request param",
		})
		return
	}

	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"details": "template not found",
		})
		return
	}

	tvi, err := strconv.Atoi(templateVersionId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"details": "template version not found",
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

	tv, err := models.RetrieveWorkspaceTemplateVersionsByIdByTemplate(*wt, uint(tvi))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if tv == nil {
		c.JSON(http.StatusOK, gin.H{
			"details": "template version not found",
		})
		return
	}

	if tv.Published {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "cannot edit a published version",
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

	tv, err = models.UpdateTemplateVersion(
		*wt,
		*tv,
		requestBody.Name,
		requestBody.Published,
		user,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, tv)
}
