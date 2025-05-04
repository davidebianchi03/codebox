package templates

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/utils/targz"
)

func HandleListTemplateVersionEntries(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")
	templateVersionId, _ := c.Params.Get("versionId")

	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	tvi, err := strconv.Atoi(templateVersionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
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
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template version not found",
		})
		return
	}

	if !tv.Sources.Exists() {
		c.JSON(http.StatusOK, []string{})
		return
	}

	tgm := targz.TarGZManager{
		Filepath: tv.Sources.GetAbsolutePath(),
	}

	files, err := tgm.EntriesTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, files)
}

func HandleRetrieveTemplateVersionFile(c *gin.Context) {

}

type CreateTemplateVersionEntryRequestBody struct {
	Path    string `json:"path" binding:"required"` // must start with a .
	Type    string `json:"type" binding:"required"`
	Content string `json:"content"`
}

func HandleCreateTemplateVersionEntry(c *gin.Context) {
	templateId, _ := c.Params.Get("templateId")
	templateVersionId, _ := c.Params.Get("versionId")

	requestBody := CreateTemplateVersionEntryRequestBody{}
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "missing or invalid request parameter",
		})
		return
	}

	if !strings.HasPrefix(requestBody.Path, ".") {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "path must start with a .",
		})
		return
	}

	if requestBody.Type != "dir" && requestBody.Type != "file" {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "invalid field type, must be 'dir' or 'file'",
		})
		return
	}

	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return
	}

	tvi, err := strconv.Atoi(templateVersionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
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
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template version not found",
		})
		return
	}

	if !tv.Sources.Exists() {
		c.JSON(http.StatusOK, []string{})
		return
	}

	tgm := targz.TarGZManager{
		Filepath: tv.Sources.GetAbsolutePath(),
	}

	// check if file already exists
	entry, err := tgm.RetrieveEntry(requestBody.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if entry != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "entry already exists",
		})
		return
	}

	if requestBody.Type == "dir" {
		if err := tgm.Mkdir(requestBody.Path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	} else {
		if err := tgm.WriteFile(requestBody.Path, []byte(requestBody.Content)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	entry, err = tgm.RetrieveEntry(requestBody.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
	}

	c.JSON(http.StatusCreated, entry)
}
