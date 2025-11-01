package runners

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RunnerGitSSH godoc
// @Summary Handle ws connection to perform git pulls/pushs over ssh
// @Schemes
// @Description Handle ws connection to perform git pulls/pushs over ssh
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200
// @Error 400
func HandleRunnerGitSSH(c *gin.Context) {
	fmt.Print("DSfsdf")
}
