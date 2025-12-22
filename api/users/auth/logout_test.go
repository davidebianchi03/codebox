package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/router"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Test the logout endpoint
*/
func TestLogout(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/logout",
			"POST",
			nil,
		)

		// create the token and set it in the request header
		// use this instead of AuthenticateHttpRequest to test logout
		// for more control over the token
		token, err := models.CreateToken(*user, 24*time.Hour)
		if err != nil {
			t.Fatalf("Failed to create token: '%s'", err)
		}
		req.Header.Set("Authorization", "Bearer "+token.Token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// try to retrieve the details of the user using the same token
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/user-details",
			"GET",
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token.Token)

		router.ServeHTTP(w, req)

		// should return 401 Unauthorized since the token was deleted on logout
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

/*
TODO:
- test logout
- test logout when not logged in
- test verify email
- test verify email with invalid token
- test verify email with expired token
- test change password
- test update profile
- test subdomains authentication
- test impersonation
*/
