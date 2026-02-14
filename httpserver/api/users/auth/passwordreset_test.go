package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/auth"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Try to reset password
*/
func TestPasswordReset(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		passwordResetReqBody := auth.RequestPasswordResetTokenBody{
			Email: user.Email,
		}

		// check that there are no token reset requests
		// in the database before the test
		count, err := models.CountPasswordResetTokensForUser(*user)
		if err != nil {
			t.Fatalf("Failed to count password reset tokens: '%s'", err)
		}
		assert.Equal(t, 0, int(count))

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/request-password-reset",
			"POST",
			passwordResetReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// now there should be one password reset token in the database
		count, err = models.CountPasswordResetTokensForUser(*user)
		if err != nil {
			t.Fatalf("Failed to count password reset tokens: '%s'", err)
		}
		assert.Equal(t, 1, int(count))

		tokens, err := models.GetPasswordResetTokensForUser(*user)
		if err != nil {
			t.Fatalf("Failed to retrieve password reset tokens: '%s'", err)
		}
		assert.Equal(t, 1, len(tokens))

		token := tokens[0]

		// try to reset the password with the token
		newPassword := "NewPassword.123"
		resetPasswordReqBody := auth.HandlePasswordResetFromTokenBody{
			Token:       token.Token,
			NewPassword: newPassword,
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/password-reset-from-token",
			"POST",
			resetPasswordReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// check that the password was actually reset
		user, err = models.RetrieveUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		assert.False(t, user.CheckPassword("password"))
		assert.True(t, user.CheckPassword("NewPassword.123"))
	})
}

/*
Try to reset password of a user that does not exist,
should return 200 but do nothing
*/
func TestRequestPasswordResetNonExistentUser(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		count, err := models.CountAllPasswordResetTokens()
		if err != nil {
			t.Fatalf("Failed to count all password reset tokens: '%s'", err)
		}
		assert.Equal(t, int64(0), count)

		passwordResetReqBody := auth.RequestPasswordResetTokenBody{
			Email: "nonexistent@user.com",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/request-password-reset",
			"POST",
			passwordResetReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		count, err = models.CountAllPasswordResetTokens()
		if err != nil {
			t.Fatalf("Failed to count all password reset tokens: '%s'", err)
		}
		assert.Equal(t, int64(0), count)
	})
}

/*
Try to reset password, but missing email field
*/
func TestRequestPasswordResetMissingEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		passwordResetReqBody := auth.RequestPasswordResetTokenBody{}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/request-password-reset",
			"POST",
			passwordResetReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing or invalid field")
	})
}

/*
Try to reset password, but invalid email
*/
func TestRequestPasswordResetInvalidEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		passwordResetReqBody := auth.RequestPasswordResetTokenBody{
			Email: "invalid-email",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/request-password-reset",
			"POST",
			passwordResetReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing or invalid field")
	})
}

/*
user logged in, should not be able to request password reset
*/
func TestRequestPasswordResetLoggedIn(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		passwordResetReqBody := auth.RequestPasswordResetTokenBody{
			Email: user.Email,
		}

		// check that there are no token reset requests
		// in the database before the test
		count, err := models.CountPasswordResetTokensForUser(*user)
		if err != nil {
			t.Fatalf("Failed to count password reset tokens: '%s'", err)
		}
		assert.Equal(t, 0, int(count))

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/request-password-reset",
			"POST",
			passwordResetReqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "already logged in")

		// now there should be one password reset token in the database
		count, err = models.CountPasswordResetTokensForUser(*user)
		if err != nil {
			t.Fatalf("Failed to count password reset tokens: '%s'", err)
		}
		assert.Equal(t, 0, int(count))
	})
}

/*
Try to reset password, but email service not configured
*/
func TestRequestPasswordResetEmailNotConfigured(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// remove configuration of email service (required to edit authentication settings)
		config.Environment.EmailSMTPHost = ""
		config.Environment.EmailSMTPPort = 0
		config.Environment.EmailSMTPUser = ""
		config.Environment.EmailSMTPPassword = ""

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		passwordResetReqBody := auth.RequestPasswordResetTokenBody{
			Email: user.Email,
		}

		// check that there are no token reset requests
		// in the database before the test
		count, err := models.CountPasswordResetTokensForUser(*user)
		if err != nil {
			t.Fatalf("Failed to count password reset tokens: '%s'", err)
		}
		assert.Equal(t, 0, int(count))

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/request-password-reset",
			"POST",
			passwordResetReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotAcceptable, w.Code)
		assert.Contains(t, w.Body.String(), "password reset is not available")

		// now there should be one password reset token in the database
		count, err = models.CountPasswordResetTokensForUser(*user)
		if err != nil {
			t.Fatalf("Failed to count password reset tokens: '%s'", err)
		}
		assert.Equal(t, 0, int(count))
	})
}

/*
Try to reset password
*/
func TestPasswordResetFromToken(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		// generate password reset token
		token, err := models.CreatePasswordResetToken(*user)
		if err != nil {
			t.Fatalf("Failed to create password reset token, %s", err)
		}

		// try to reset the password with the token
		newPassword := "NewPassword.123"
		resetPasswordReqBody := auth.HandlePasswordResetFromTokenBody{
			Token:       token.Token,
			NewPassword: newPassword,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/password-reset-from-token",
			"POST",
			resetPasswordReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// check that the password was actually reset
		user, err = models.RetrieveUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		assert.False(t, user.CheckPassword("password"))
		assert.True(t, user.CheckPassword("NewPassword.123"))
	})
}

/*
Try to reset password with missing token
*/
func TestPasswordResetFromTokenMissingToken(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		// try to reset the password with the token
		newPassword := "NewPassword.123"
		resetPasswordReqBody := auth.HandlePasswordResetFromTokenBody{
			NewPassword: newPassword,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/password-reset-from-token",
			"POST",
			resetPasswordReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing or invalid field")

		// check that the password was actually reset
		user, err = models.RetrieveUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		assert.True(t, user.CheckPassword("password"))
	})
}

/*
Try to reset password with missing token
*/
func TestPasswordResetFromTokenMissingPassword(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		// generate password reset token
		token, err := models.CreatePasswordResetToken(*user)
		if err != nil {
			t.Fatalf("Failed to create password reset token, %s", err)
		}

		// try to reset the password with the token
		resetPasswordReqBody := auth.HandlePasswordResetFromTokenBody{
			Token: token.Token,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/password-reset-from-token",
			"POST",
			resetPasswordReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing or invalid field")

		// check that the password was actually reset
		user, err = models.RetrieveUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		assert.True(t, user.CheckPassword("password"))
	})
}

/*
Try to reset password with an invalid token
*/
func TestPasswordResetFromTokenInvalidToken(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		// try to reset the password with the token
		newPassword := "NewPassword.123"
		resetPasswordReqBody := auth.HandlePasswordResetFromTokenBody{
			Token:       "invalid-token",
			NewPassword: newPassword,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/password-reset-from-token",
			"POST",
			resetPasswordReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "invalid or expired token")

		// check that the password was actually reset
		user, err = models.RetrieveUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		assert.True(t, user.CheckPassword("password"))
	})
}

/*
Try to reset password with an expired token
*/
func TestPasswordResetFromTokenWithAnExpiredToken(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		// generate an expired password reset token
		token, err := models.CreatePasswordResetToken(*user)
		if err != nil {
			t.Fatalf("Failed to create password reset token, %s", err)
		}
		token.Expiration = time.Now()
		if err := models.UpdatePasswordResetToken(token); err != nil {
			t.Fatalf("Failed to update password reset token, %s", err)
		}

		// try to reset the password with the token
		newPassword := "NewPassword.123"
		resetPasswordReqBody := auth.HandlePasswordResetFromTokenBody{
			Token:       token.Token,
			NewPassword: newPassword,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/password-reset-from-token",
			"POST",
			resetPasswordReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "invalid or expired token")

		// check that the password was actually reset
		user, err = models.RetrieveUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		assert.True(t, user.CheckPassword("password"))
	})
}

/*
Try to reset password with an invalid password
*/
func TestPasswordResetFromTokenWithAnInvalidPassword(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		// generate an expired password reset token
		token, err := models.CreatePasswordResetToken(*user)
		if err != nil {
			t.Fatalf("Failed to create password reset token, %s", err)
		}

		// try to reset the password with the token
		newPassword := "invalidpassword"
		resetPasswordReqBody := auth.HandlePasswordResetFromTokenBody{
			Token:       token.Token,
			NewPassword: newPassword,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/password-reset-from-token",
			"POST",
			resetPasswordReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid password")

		// check that the password was actually reset
		user, err = models.RetrieveUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		assert.True(t, user.CheckPassword("password"))
	})
}

/*
Try to reset password but user was logged in
*/
func TestPasswordResetFromTokenuserLoggedIn(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		// check password
		assert.True(t, user.CheckPassword("password"))

		// generate password reset token
		token, err := models.CreatePasswordResetToken(*user)
		if err != nil {
			t.Fatalf("Failed to create password reset token, %s", err)
		}

		// try to reset the password with the token
		newPassword := "NewPassword.123"
		resetPasswordReqBody := auth.HandlePasswordResetFromTokenBody{
			Token:       token.Token,
			NewPassword: newPassword,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/password-reset-from-token",
			"POST",
			resetPasswordReqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "already logged in")

		// check that the password was actually reset
		user, err = models.RetrieveUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to retrieve user: '%s'", err)
		}

		if user == nil {
			t.Fatal("User not found")
		}

		assert.True(t, user.CheckPassword("password"))
	})
}

/*
test api to check if password reset is available
*/
func TestPasswordResetAvailable(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// configure email sender
		config.Environment.EmailSMTPHost = "smtp.example.com"
		config.Environment.EmailSMTPPort = 25
		config.Environment.EmailSMTPUser = "user@example.com"
		config.Environment.EmailSMTPPassword = "password"

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/can-reset-password",
			"GET",
			nil,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseBody map[string]bool
		if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
			t.Fatalf("Failed to parse json, %s", err)
		}

		assert.True(t, responseBody["can_reset_password"])

		// disable password reset and test again
		config.Environment.EmailSMTPHost = ""
		config.Environment.EmailSMTPPort = 0
		config.Environment.EmailSMTPUser = ""
		config.Environment.EmailSMTPPassword = ""

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/can-reset-password",
			"GET",
			nil,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
			t.Fatalf("Failed to parse json, %s", err)
		}

		assert.False(t, responseBody["can_reset_password"])
	})
}
