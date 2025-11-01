package templates

import (
	"encoding/base64"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/utils/targz"
)

// retrieve workspace template version from context, return nil if not found
// this function writes http responses
func getTemplateVersionFromContext(c *gin.Context) *models.WorkspaceTemplateVersion {
	templateId, _ := c.Params.Get("templateId")
	templateVersionId, _ := c.Params.Get("versionId")

	ti, err := strconv.Atoi(templateId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template not found",
		})
		return nil
	}

	tvi, err := strconv.Atoi(templateVersionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template version not found",
		})
		return nil
	}

	wt, err := models.RetrieveWorkspaceTemplateByID(uint(ti))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return nil
	}

	tv, err := models.RetrieveWorkspaceTemplateVersionsByIdByTemplate(*wt, uint(tvi))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return nil
	}

	if tv == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "template version not found",
		})
		return nil
	}

	return tv
}

// ListTemplateVersionEntries godoc
// @Summary List template version entries
// @Schemes
// @Description List template version entries
// @Tags Templates
// @Accept json
// @Produce json
// @Success 201 {object} []targz.TarTreeItem
// @Router /api/v1/templates/:templateId/versions/:versionId/entries [get]
func HandleListTemplateVersionEntries(c *gin.Context) {
	tv := getTemplateVersionFromContext(c)
	if tv == nil {
		return
	}

	tgm := targz.TarGZManager{
		Filepath: tv.Sources.GetAbsolutePath(),
	}

	if !tv.Sources.Exists() {
		if err := tgm.CreateArchive(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	files, err := tgm.EntriesTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if len(files) > 0 {
		if files[0].Name == "." {
			files = files[0].Children
		}
	}

	c.JSON(http.StatusOK, files)
}

// RetrieveTemplateVersionFile godoc
// @Summary Retrieve template version entry
// @Schemes
// @Description Retrieve template version entry
// @Tags Templates
// @Accept json
// @Produce json
// @Success 201 {object} targz.TarEntry
// @Router /api/v1/templates/:templateId/versions/:versionId/entries/:path [get]
func HandleRetrieveTemplateVersionFile(c *gin.Context) {
	path, _ := c.Params.Get("path")
	path = "./" + strings.TrimPrefix(strings.TrimPrefix(path, "/"), "./")

	tv := getTemplateVersionFromContext(c)
	if tv == nil {
		return
	}

	tgm := targz.TarGZManager{
		Filepath: tv.Sources.GetAbsolutePath(),
	}

	if !tv.Sources.Exists() {
		if err := tgm.CreateArchive(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	entry, err := tgm.RetrieveEntry(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if entry == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "entry not found",
		})
		return
	}

	entry.Path = strings.TrimPrefix(entry.Path, "./")

	c.JSON(http.StatusOK, entry)
}

type CreateTemplateVersionEntryRequestBody struct {
	Path    string `json:"path" binding:"required"` // must start with a .
	Type    string `json:"type" binding:"required"`
	Content string `json:"content"`
}

// CreateTemplateVersionEntry godoc
// @Summary Create new template version entry
// @Schemes
// @Description Create new template version entry
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body CreateTemplateVersionEntryRequestBody true "Template version entry data"
// @Success 201 {object} targz.TarEntry
// @Router /api/v1/templates/:templateId/versions/:templateId/entries [post]
func HandleCreateTemplateVersionEntry(c *gin.Context) {
	requestBody := CreateTemplateVersionEntryRequestBody{}
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "missing or invalid request parameter",
		})
		return
	}

	if requestBody.Type != "dir" && requestBody.Type != "file" {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "invalid field type, must be 'dir' or 'file'",
		})
		return
	}

	if strings.HasPrefix(requestBody.Path, "./") || strings.HasPrefix(requestBody.Path, ".") {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "path cannot start with ./",
		})
		return
	}

	if strings.HasPrefix(requestBody.Path, "/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "path cannot start with /",
		})
		return
	}

	path := "./" + requestBody.Path

	tv := getTemplateVersionFromContext(c)
	if tv == nil {
		return
	}

	if tv.Published {
		c.JSON(http.StatusLocked, gin.H{
			"details": "cannot edit a template version that has already been released",
		})
		return
	}

	tgm := targz.TarGZManager{
		Filepath: tv.Sources.GetAbsolutePath(),
	}

	if !tv.Sources.Exists() {
		if err := tgm.CreateArchive(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	// check if parent element exists and is a folder
	parentEntryPath := filepath.Dir(strings.TrimSuffix(path, "/"))
	if parentEntryPath != "." {
		parentEntryPath = "./" + parentEntryPath
		parentEntry, err := tgm.RetrieveEntry(parentEntryPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}

		if parentEntry == nil {
			parts := strings.Split(path, "/")

			for i := 0; i < len(parts); i++ {
				p := strings.Join(parts[:i+1], "/")
				entry, err := tgm.RetrieveEntry(p)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"details": "internal server error",
					})
					return
				}

				if entry != nil {
					if entry.Type != "dir" {
						c.JSON(http.StatusBadRequest, gin.H{
							"details": "parent entry is not a directory",
						})
						return
					}
				}
			}

			if err := tgm.MkDirAll(parentEntryPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"details": "internal server error",
				})
				return
			}
		} else {
			if parentEntry.Type != "dir" {
				c.JSON(http.StatusBadRequest, gin.H{
					"details": "parent entry is not a directory",
				})
				return
			}
		}
	}

	// check if file already exists
	entry, err := tgm.RetrieveEntry(path)
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
		if err := tgm.Mkdir(path); err != nil {
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	} else {
		if strings.HasSuffix(path, "/") {
			c.JSON(http.StatusBadRequest, gin.H{
				"details": "filename may not have trailing slash",
			})
			return
		}

		content, err := base64.StdEncoding.DecodeString(requestBody.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"details": "invalid content, it must be a base64 string",
			})
			return
		}

		if err := tgm.WriteFile(path, content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	entry, err = tgm.RetrieveEntry(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
	}

	c.JSON(http.StatusCreated, entry)
}

type UpdateTemplateVersionEntryRequestBody struct {
	Path    string  `json:"path" binding:"required"`
	Content *string `json:"content" binding:"required"`
}

// CreateTemplateVersionEntry godoc
// @Summary Updates a template version entry
// @Schemes
// @Description Updates a template version entry
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body UpdateTemplateVersionEntryRequestBody true "Template version entry data"
// @Success 200 {object} targz.TarEntry
// @Router /api/v1/templates/:templateId/versions/:templateId/entries/:path [put]
func HandleUpdateTemplateVersionEntry(c *gin.Context) {
	path, _ := c.Params.Get("path")
	path = strings.TrimPrefix(strings.TrimPrefix(path, "/"), "./")
	path = "./" + strings.TrimPrefix(strings.TrimPrefix(path, "/"), "./")

	requestBody := UpdateTemplateVersionEntryRequestBody{}
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "missing or invalid request parameter",
		})
		return
	}

	if strings.HasPrefix(requestBody.Path, "/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "path cannot start with /",
		})
		return
	}

	newPath := "./" + strings.TrimPrefix(strings.TrimPrefix(requestBody.Path, "/"), "./")
	if newPath == "." || newPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "invalid path",
		})
		return
	}

	tv := getTemplateVersionFromContext(c)
	if tv == nil {
		return
	}

	if tv.Published {
		c.JSON(http.StatusLocked, gin.H{
			"details": "cannot edit a template version that has already been released",
		})
		return
	}

	tgm := targz.TarGZManager{
		Filepath: tv.Sources.GetAbsolutePath(),
	}

	if !tv.Sources.Exists() {
		if err := tgm.CreateArchive(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	entry, err := tgm.RetrieveEntry(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if entry == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "entry not found",
		})
		return
	}

	// if path has been changed
	if path != newPath {
		destinationEntry, err := tgm.RetrieveEntry(newPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}

		if destinationEntry != nil {
			c.JSON(http.StatusConflict, gin.H{
				"details": "destination already exists",
			})
			return
		}

		// check if parent element exists and is a folder
		parentEntryPath := filepath.Dir(strings.TrimSuffix(newPath, "/"))
		if parentEntryPath != "." {
			parentEntryPath = "./" + parentEntryPath
			parentEntry, err := tgm.RetrieveEntry(parentEntryPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"details": "internal server error",
				})
				return
			}

			if parentEntry == nil {
				parts := strings.Split(newPath, "/")

				for i := 0; i < len(parts); i++ {
					p := strings.Join(parts[:i+1], "/")
					entry, err := tgm.RetrieveEntry(p)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"details": "internal server error",
						})
						return
					}

					if entry != nil {
						if entry.Type != "dir" {
							c.JSON(http.StatusBadRequest, gin.H{
								"details": "parent entry is not a directory",
							})
							return
						}
					}
				}

				if err := tgm.MkDirAll(parentEntryPath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"details": "internal server error",
					})
					return
				}
			} else {
				if parentEntry.Type != "dir" {
					c.JSON(http.StatusBadRequest, gin.H{
						"details": "parent entry is not a directory",
					})
					return
				}
			}
		}

		if err := tgm.Move(path, newPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	// update the content
	if *requestBody.Content != string(entry.Content) {
		content, err := base64.StdEncoding.DecodeString(*requestBody.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"details": "invalid content, it must be a base64 string",
			})
			return
		}

		if err := tgm.WriteFile(newPath, content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	entry, err = tgm.RetrieveEntry(newPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// DeleteTemplateVersionEntry godoc
// @Summary Delete template version entry
// @Schemes
// @Description Delete template version entry
// @Tags Templates
// @Accept json
// @Produce json
// @Success 204
// @Router /api/v1/templates/:templateId/versions/:templateId/entries/:path [delete]
func HandleDeleteTemplateVersionEntry(c *gin.Context) {
	path, _ := c.Params.Get("path")
	path = strings.TrimPrefix(strings.TrimPrefix(path, "/"), "./")
	path = "./" + strings.TrimPrefix(strings.TrimPrefix(path, "/"), "./")

	tv := getTemplateVersionFromContext(c)
	if tv == nil {
		return
	}

	if tv.Published {
		c.JSON(http.StatusLocked, gin.H{
			"details": "cannot edit a template version that has already been released",
		})
		return
	}

	tgm := targz.TarGZManager{
		Filepath: tv.Sources.GetAbsolutePath(),
	}

	if !tv.Sources.Exists() {
		if err := tgm.CreateArchive(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"details": "internal server error",
			})
			return
		}
	}

	entry, err := tgm.RetrieveEntry(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	if entry == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"details": "entry not found",
		})
		return
	}

	if err := tgm.Delete(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "internal server error",
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"details": "entry has been deleted",
	})
}
