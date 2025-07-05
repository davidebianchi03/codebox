package admin

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

func RandStringBytesRmndr(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-#!_=+"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func HandleAdminListRunners(c *gin.Context) {
	var runners []models.Runner
	r := dbconn.DB.Find(&runners)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, runners)
}

func HandleAdminRetrieveRunners(c *gin.Context) {
	runnerId, _ := c.Params.Get("runnerId")

	var runner models.Runner
	r := dbconn.DB.Find(&runner, map[string]interface{}{
		"id": runnerId,
	})
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if runner.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "runner not found",
		})
		return
	}

	c.JSON(http.StatusOK, runner)
}

func HandleAdminCreateRunner(c *gin.Context) {
	type RequestBody struct {
		Name         string `json:"name" binding:"required"`
		Type         string `json:"type" binding:"required"`
		UsePublicUrl bool   `json:"use_public_url"`
		PublicUrl    string `json:"public_url"`
	}

	var parsedBody RequestBody
	err := c.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	var exists bool
	err = dbconn.DB.Model(models.Runner{}).Select("count(*) > 0").Where("name = ?", parsedBody.Name).Find(&exists).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"detail": "another runner with the same name already exists",
		})
		return
	}

	runnerTypeFound := false
	for _, rt := range config.ListAvailableRunnerTypes() {
		if rt.ID == parsedBody.Type {
			runnerTypeFound = true
		}
	}

	if !runnerTypeFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "runner type not found",
		})
		return
	}

	if parsedBody.UsePublicUrl {
		if parsedBody.PublicUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "'public_url' is required",
			})
			return
		}

		err = dbconn.DB.Model(models.Runner{}).Select("count(*) > 0").Where("public_url = ?", parsedBody.PublicUrl).Find(&exists).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}

		if exists {
			c.JSON(http.StatusConflict, gin.H{
				"detail": "another runner with the same public url already exists",
			})
			return
		}
	}

	token := fmt.Sprintf("cbrt-%s", RandStringBytesRmndr(30))

	runner := models.Runner{
		Name:         parsedBody.Name,
		Type:         parsedBody.Type,
		Token:        token,
		UsePublicUrl: parsedBody.UsePublicUrl,
		PublicUrl:    parsedBody.PublicUrl,
	}

	if err = dbconn.DB.Create(&runner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	bgtasks.BgTasksEnqueuer.Enqueue("ping_runners", work.Q{})

	c.JSON(http.StatusCreated, gin.H{
		"id":    runner.ID,
		"token": token,
	})
}

func HandleAdminUpdateRunner(c *gin.Context) {
	runnerId, _ := c.Params.Get("runnerId")

	var runner models.Runner
	r := dbconn.DB.Find(&runner, map[string]interface{}{
		"id": runnerId,
	})
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if runner.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "runner not found",
		})
		return
	}

	type RequestBody struct {
		Name         string `json:"name" binding:"required"`
		Type         string `json:"type" binding:"required"`
		UsePublicUrl *bool  `json:"use_public_url" binding:"required"`
		PublicUrl    string `json:"public_url" binding:"required"`
	}

	var reqBody RequestBody
	err := c.ShouldBindBodyWithJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid field",
		})
		return
	}

	runner.Name = reqBody.Name
	runner.Type = reqBody.Type
	runner.UsePublicUrl = *reqBody.UsePublicUrl
	runner.PublicUrl = reqBody.PublicUrl
	dbconn.DB.Save(&runner)

	c.JSON(http.StatusOK, runner)
}
