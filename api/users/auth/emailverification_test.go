package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/api/users/auth"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/router"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Try to verify email with a valid code and check
that the user's email is marked as verified.
*/
func TestVerifyEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		user.EmailVerified = false
		if err := models.UpdateUser(user); err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// create verification code
		expiration := time.Now().Add(15 * time.Minute)
		code, err := models.CreateEmailVerificationCode(&expiration, *user)
		if err != nil {
			t.Fatalf("Failed to create verification code: %v", err)
		}

		loginReqBody := auth.VerifyEmailAddressRequestBody{
			Code: code.Code,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/verify-email-address",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "email has been verified")

		// verify that user's email is marked as verified
		updatedUser, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve updated user: %v", err)
		}
		assert.True(t, updatedUser.EmailVerified)
	})
}

/*
Try to verify email with an expired code
*/
func TestVerifyEmailExpiredCode(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		user.EmailVerified = false
		if err := models.UpdateUser(user); err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// create an expiredverification code
		expiration := time.Now().Add(-15 * time.Minute)
		code, err := models.CreateEmailVerificationCode(&expiration, *user)
		if err != nil {
			t.Fatalf("Failed to create verification code: %v", err)
		}

		loginReqBody := auth.VerifyEmailAddressRequestBody{
			Code: code.Code,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/verify-email-address",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotAcceptable, w.Code)
		assert.Contains(t, w.Body.String(), "verification code is not valid or is expired")

		// verify that user's email is marked as verified
		updatedUser, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve updated user: %v", err)
		}
		assert.False(t, updatedUser.EmailVerified)
	})
}

/*
Try to verify email with an invalid code
*/
func TestVerifyInvalidCode(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		user.EmailVerified = false
		if err := models.UpdateUser(user); err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		loginReqBody := auth.VerifyEmailAddressRequestBody{
			Code: "invalid-code",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/verify-email-address",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotAcceptable, w.Code)
		assert.Contains(t, w.Body.String(), "verification code is not valid or is expired")

		// verify that user's email is marked as verified
		updatedUser, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve updated user: %v", err)
		}
		assert.False(t, updatedUser.EmailVerified)
	})
}

/*
Try to verify email while user is logged in
*/
func TestVerifyEmailUserLoggedIn(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		user.EmailVerified = false
		if err := models.UpdateUser(user); err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// create verification code
		expiration := time.Now().Add(15 * time.Minute)
		code, err := models.CreateEmailVerificationCode(&expiration, *user)
		if err != nil {
			t.Fatalf("Failed to create verification code: %v", err)
		}

		loginReqBody := auth.VerifyEmailAddressRequestBody{
			Code: code.Code,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/verify-email-address",
			"POST",
			loginReqBody,
		)

		// authenticate the request
		admin, err := models.RetrieveUserByEmail("admin@admin.com")
		if err != nil {
			t.Fatalf("Failed to retrieve admin user: %v", err)
		}
		testutils.AuthenticateHttpRequest(t, req, *admin)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusPreconditionFailed, w.Code)
		assert.Contains(t, w.Body.String(), "logged in users cannot verify email")

		// verify that user's email is marked as verified
		updatedUser, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve updated user: %v", err)
		}
		assert.False(t, updatedUser.EmailVerified)
	})
}

/*
Try to verify email of a user whose email is already verified
*/
func TestVerifyEmailAddressAlreadyVerified(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		user.EmailVerified = true
		if err := models.UpdateUser(user); err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// create verification code
		expiration := time.Now().Add(15 * time.Minute)
		code, err := models.CreateEmailVerificationCode(&expiration, *user)
		if err != nil {
			t.Fatalf("Failed to create verification code: %v", err)
		}

		loginReqBody := auth.VerifyEmailAddressRequestBody{
			Code: code.Code,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/verify-email-address",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "email has already been verified")
	})
}

/*
Try to verify email with a missing code field
*/
func TestVerifyEmailMissingCode(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		user.EmailVerified = false
		if err := models.UpdateUser(user); err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		loginReqBody := auth.VerifyEmailAddressRequestBody{}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/verify-email-address",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

/*
Try to verify email with an empty code field
*/
func TestVerifyEmailEmptyCode(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := router.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Fatalf("Failed to retrieve user: %v", err)
		}

		user.EmailVerified = false
		if err := models.UpdateUser(user); err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		loginReqBody := auth.VerifyEmailAddressRequestBody{
			Code: "",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/verify-email-address",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
