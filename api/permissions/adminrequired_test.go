package permissions_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidebianchi03/codebox/api/permissions"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	api := permissions.AdminRequiredRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"detail": "ok",
		})
	})

	api(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code) // Expect 401 Unauthorized
	assert.JSONEq(t, `{"error": "Unauthorized"}`, w.Body.String())
}
