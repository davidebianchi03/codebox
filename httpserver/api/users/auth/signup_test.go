package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/auth"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Try to signup
*/
func TestSignup(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		usersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}

		// enable signup in config
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		s.IsSignUpOpen = true
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		loginReqBody := auth.SignUpRequestBody{
			Email:     "pippo@pluto.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// verify that a new user was created
		newUsersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, usersCount+1, newUsersCount)
	})
}

/*
Try to signup existing email
*/
func TestSignupExistingEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		usersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}

		// enable signup in config
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		s.IsSignUpOpen = true
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		loginReqBody := auth.SignUpRequestBody{
			Email:     "admin@admin.com",
			FirstName: "user1",
			LastName:  "user",
			Password:  "Password.123455",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		// email already exists, send status created to avoid user enumeration
		assert.Equal(t, http.StatusCreated, w.Code)

		newUsersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, usersCount, newUsersCount)
	})
}

/*
Test signup first user
*/
func TestSignupFirstUser(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// signup is initially closed so don't need to change settings,
		// but first user can signup anyway

		// delete all users
		users, err := models.ListUsers(-1)
		if err != nil {
			t.Fatalf("Failed to list users: '%s'", err)
		}

		for _, user := range *users {
			if err := models.DeleteUser(&user); err != nil {
				t.Fatalf("Failed to delete user: '%s'", err)
			}
		}

		usersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, 0, int(usersCount))

		loginReqBody := auth.SignUpRequestBody{
			Email:     "first@user.com",
			FirstName: "user1",
			LastName:  "user",
			Password:  "Password.123455",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		// email already exists, send status created to avoid user enumeration
		assert.Equal(t, http.StatusCreated, w.Code)

		newUsersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, 1, int(newUsersCount))

		user, err := models.RetrieveUserByEmail("first@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		assert.True(t, user.IsSuperuser, "first signed up user should be a superuser")
		assert.True(t, user.EmailVerified, "first signed up user's email should be verified")
		assert.True(t, user.Approved, "first signed up user should be approved")
	})
}

/*
Try to signup, but registration is closed
*/
func TestSignupRegistrationClosed(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		usersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}

		// enable signup in config
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		s.IsSignUpOpen = false
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		loginReqBody := auth.SignUpRequestBody{
			Email:     "pippo@pluto.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotAcceptable, w.Code)

		// verify that a new user was created
		newUsersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, usersCount, newUsersCount)
	})
}

/*
Try to signup, registration is restricted
*/
func TestSignupRegistrationRestricted(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		usersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}

		// enable signup in config
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		s.IsSignUpOpen = true
		s.IsSignUpRestricted = true
		s.AllowedEmailRegex = `^.*@allowed-domain\.com$`
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		// try to signup with an email that matches a regex
		loginReqBody := auth.SignUpRequestBody{
			Email:     "user@allowed-domain.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// verify that a new user was created
		newUsersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, usersCount+1, newUsersCount)

		// try to signup with an email that does not match a regex
		loginReqBody = auth.SignUpRequestBody{
			Email:     "utente@domain.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotAcceptable, w.Code)

		// verify that a new user was not created
		newUsersCount, err = models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, usersCount+1, newUsersCount)
	})
}

/*
Try to signup with a blacklisted email
*/
func TestSignupBlacklistedEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		usersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}

		// enable signup in config
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		s.IsSignUpOpen = true
		s.BlockedEmailRegex = `^.*@blacklisted-domain\.com$`
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		// try to signup with an email that matches a regex
		loginReqBody := auth.SignUpRequestBody{
			Email:     "user@blacklisted-domain.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotAcceptable, w.Code)

		// verify that a new user was created
		newUsersCount, err := models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, usersCount, newUsersCount)

		// try to signup with an email that does not match a regex
		loginReqBody = auth.SignUpRequestBody{
			Email:     "utente@domain.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// verify that a new user was not created
		newUsersCount, err = models.CountAllUsers()
		if err != nil {
			t.Fatalf("Failed to count users: '%s'", err)
		}
		assert.Equal(t, usersCount+1, newUsersCount)
	})
}

/*
Try to signup with invalid email
*/
func TestSignupInvalidEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// enable signup in config
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		s.IsSignUpOpen = true
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		loginReqBody := auth.SignUpRequestBody{
			Email:     "invalid-email-format",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

/*
Try to signup with invalid email
*/
func TestSignupInvalidPassword(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// enable signup in config
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		s.IsSignUpOpen = true
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		loginReqBody := auth.SignUpRequestBody{
			Email:     "test@user.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "invalidpassword",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

/*
Try to signup with missing fields
*/
func TestSignupMissingFields(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		// enable signup in config
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		s.IsSignUpOpen = true
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		testCases := []auth.SignUpRequestBody{
			{
				FirstName: "pippo",
				LastName:  "pluto",
				Password:  "Password.12345",
			},
			{
				Email:     "",
				FirstName: "pippo",
				LastName:  "pluto",
				Password:  "Password.12345",
			},
			{
				Email:     "user@test.com",
				FirstName: "",
				LastName:  "pluto",
				Password:  "Password.12345",
			},
			{
				Email:     "user@test.com",
				FirstName: "pippo",
				LastName:  "",
				Password:  "Password.12345",
			},
			{
				Email:     "user@test.com",
				FirstName: "pippo",
				LastName:  "pluto",
				Password:  "",
			},
		}

		for _, loginReqBody := range testCases {
			w := httptest.NewRecorder()
			req := testutils.CreateRequestWithJSONBody(
				t,
				"/api/v1/auth/signup",
				"POST",
				loginReqBody,
			)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
	})
}

/*
Test the auto approve behavior
*/
func TestSignupUserMatchingAutoApprovedRegex(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Fatalf("Failed to retrieve instance settings: '%s'", err)
			return
		}
		// enable signup in config
		s.IsSignUpOpen = true
		// add an entry to auto approved users regex
		s.ApprovedByDefaultEmailRegex = `^.*@autoapproved\.com$`
		if err := models.SaveSingletonModel(s); err != nil {
			t.Fatalf("Failed to update instance settings: '%s'", err)
			return
		}

		loginReqBody := auth.SignUpRequestBody{
			Email:     "pippo@pluto.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		user, err := models.RetrieveUserByEmail("pippo@pluto.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}
		// this user should not be approved
		assert.False(t, user.Approved)

		// now signup with an email that ends with autoapproved.com
		loginReqBody = auth.SignUpRequestBody{
			Email:     "pippo@autoapproved.com",
			FirstName: "pippo",
			LastName:  "pluto",
			Password:  "Password.123455",
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/auth/signup",
			"POST",
			loginReqBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		user, err = models.RetrieveUserByEmail("pippo@autoapproved.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}
		// this user should be approved
		assert.True(t, user.Approved)
	})
}
