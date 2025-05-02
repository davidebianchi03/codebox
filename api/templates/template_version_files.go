package templates

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/utils/targz"
)

func HandleListTemplateVersionFiles(c *gin.Context) {
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
