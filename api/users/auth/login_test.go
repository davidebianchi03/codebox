package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/api/users/auth"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/router"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Try to login with valid credentials
Expect 200 OK and a token in the response body
*/
func TestLogin(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		loginReqBody := auth.LoginRequestBody{
			Email:    "user1@user.com",
			Password: "password",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/login",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})
}

/*
Try to login with wrong email
*/
func TestLoginInvalidEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		loginReqBody := auth.LoginRequestBody{
			Email:    "doesnotexists@user.com",
			Password: "password",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/login",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})
}

/*
Try to login with wrong password
*/
func TestLoginInvalidPassword(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		loginReqBody := auth.LoginRequestBody{
			Email:    "user1@user.com",
			Password: "wrongpassword",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/login",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})
}

/*
Try to login without email
*/
func TestLoginNoEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		loginReqBody := auth.LoginRequestBody{
			Password: "wrongpassword",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/login",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing or invalid field")
	})
}

/*
Try to login with invalid email format
*/
func TestLoginInvalidEmailFormat(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		loginReqBody := auth.LoginRequestBody{
			Email: "invalid email format",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/login",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing or invalid field")
	})
}

/*
Try to login without email
*/
func TestLoginNoPassword(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		loginReqBody := auth.LoginRequestBody{
			Email: "user1@user.com",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/login",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing or invalid field")
	})
}

/*
With a user already logged in, try to login again
Expect 400 Bad Request
*/
func TestLoginAlreadyLoggedIn(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		loginReqBody := auth.LoginRequestBody{
			Email:    "user1@user.com",
			Password: "password",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/login",
			"POST",
			loginReqBody,
		)

		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "already logged in")
	})
}

/*
Try to login with unverified email
*/
func TestLoginUnverifiedEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		user.EmailVerified = false
		if err := models.UpdateUser(user); err != nil {
			t.Fatalf("Failed to update test user: '%s'", err)
			return
		}

		loginReqBody := auth.LoginRequestBody{
			Email:    "user1@user.com",
			Password: "password",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/login",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusPreconditionFailed, w.Code)
		assert.Contains(t, w.Body.String(), "the email address has not yet been verified")
	})
}
