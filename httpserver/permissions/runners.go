package permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

/*
Wrap a Gin handler to require that the request is
authenticated using a valid runner token.
The runner urls must contain the :runnerId parameter.
If the token is missing or invalid, returns 401 Unauthorized.
Otherwise, calls the original handler.
*/
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

		runnerToken := c.Request.Header.Get(config.Environment.RunnerTokenHeader)

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
