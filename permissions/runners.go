package permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

func RunnerTokenAuthenticationRequired(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// the runner urls must contain arg :runnerId
		runnerId, err := utils.GetUIntParamFromContext(c, "runnerId")
		if err != nil {
			utils.ErrorResponse(
				c,
				http.StatusUnauthorized,
				"missing or invalid token",
			)
			return
		}

		runnerToken := c.Request.Header.Get("X-Codebox-Runner-Token")

		runner, err := models.RetrieveRunnerByID(uint(runnerId))
		if err != nil {
			utils.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"internal server error",
			)
			return
		}

		if runner.Token != runnerToken {
			utils.ErrorResponse(
				c,
				http.StatusUnauthorized,
				"missing or invalid token",
			)
			return
		} else {
			handler(c)
		}
	}
}
