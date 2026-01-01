package settings_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/admin/settings"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Try to retrieve authentication settings,
try first with an admin, then with a user that is not an admin
and at the end without authentication
*/
func TestRetrieveAuthenticationSettings(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		adminUser, err := models.RetrieveUserByEmail("admin@admin.com")
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/admin/authentication-settings",
			"GET",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *adminUser)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		as, err := serializers.AuthenticationSettingsSerializerFromJSON(
			w.Body.String(),
		)
		if err != nil {
			t.Error(err)
		}

		assert.False(t, as.IsSignUpOpen)
		assert.False(t, as.IsSignUpRestricted)
		assert.False(t, as.UsersMustBeApproved)
		assert.Equal(t, "", as.AllowedEmailRegex)
		assert.Equal(t, "", as.BlockedEmailRegex)
		assert.Equal(t, "", as.ApprovedByDefaultEmailRegex)

		// try with a user that is not an admin
		commonUser, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Error(err)
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/admin/authentication-settings",
			"GET",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *commonUser)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		// try without authentication
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/admin/authentication-settings",
			"GET",
			nil,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

/*
Test endpoint to update authentication settings
*/
func TestUpdateAuthenticationSettings(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		adminUser, err := models.RetrieveUserByEmail("admin@admin.com")
		if err != nil {
			t.Error(err)
		}

		signupOpen := true
		signupRestricted := true
		usersMustBeApproved := true
		allowedEmailRegex := `^.*@example\.com$`
		blockedEmailRegex := `^.*@blocked\.com$`
		approvedByDefaultEmailRegex := `^.*@example\.com$`

		requestBody := settings.HandleUpdateServerSettingsRequestBody{
			IsSignUpOpen:                &signupOpen,
			IsSignUpRestricted:          &signupRestricted,
			UsersMustBeApproved:         &usersMustBeApproved,
			AllowedEmailRegex:           &allowedEmailRegex,
			BlockedEmailRegex:           &blockedEmailRegex,
			ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/admin/authentication-settings",
			"PUT",
			requestBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *adminUser)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
		if err != nil {
			t.Error(err)
		}
		assert.True(t, s.IsSignUpOpen)
		assert.True(t, s.IsSignUpRestricted)
		assert.True(t, s.UsersMustBeApproved)
		assert.Equal(t, s.AllowedEmailRegex, `^.*@example\.com$`)
		assert.Equal(t, s.BlockedEmailRegex, `^.*@blocked\.com$`)
		assert.Equal(t, s.ApprovedByDefaultEmailRegex, `^.*@example\.com$`)
	})
}

type TestUpdateAuthSettingsTestCase struct {
	RequestBody        settings.HandleUpdateServerSettingsRequestBody
	ExpectedStatusCode int
}

/*
Test endpoint to update authentication settings with missing arguments
*/
func TestUpdateAuthenticationSettingsMissingArg(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		adminUser, err := models.RetrieveUserByEmail("admin@admin.com")
		if err != nil {
			t.Error(err)
		}

		signupOpen := true
		signupRestricted := true
		usersMustBeApproved := true
		allowedEmailRegex := `^.*@example\.com$`
		blockedEmailRegex := `^.*@blocked\.com$`
		approvedByDefaultEmailRegex := `^.*@example\.com$`

		testCases := []TestUpdateAuthSettingsTestCase{
			{
				RequestBody: settings.HandleUpdateServerSettingsRequestBody{
					IsSignUpOpen:        &signupOpen,
					IsSignUpRestricted:  &signupRestricted,
					UsersMustBeApproved: &usersMustBeApproved,
					AllowedEmailRegex:   &allowedEmailRegex,
					BlockedEmailRegex:   &blockedEmailRegex,
				},
				ExpectedStatusCode: http.StatusBadRequest,
			},
			{
				RequestBody: settings.HandleUpdateServerSettingsRequestBody{
					IsSignUpOpen:                &signupOpen,
					IsSignUpRestricted:          &signupRestricted,
					UsersMustBeApproved:         &usersMustBeApproved,
					AllowedEmailRegex:           &allowedEmailRegex,
					ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
				},
				ExpectedStatusCode: http.StatusBadRequest,
			},
			{
				RequestBody: settings.HandleUpdateServerSettingsRequestBody{
					IsSignUpOpen:                &signupOpen,
					IsSignUpRestricted:          &signupRestricted,
					UsersMustBeApproved:         &usersMustBeApproved,
					BlockedEmailRegex:           &blockedEmailRegex,
					ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
				},
				ExpectedStatusCode: http.StatusBadRequest,
			},
			{
				RequestBody: settings.HandleUpdateServerSettingsRequestBody{
					IsSignUpOpen:                &signupOpen,
					IsSignUpRestricted:          &signupRestricted,
					AllowedEmailRegex:           &allowedEmailRegex,
					BlockedEmailRegex:           &blockedEmailRegex,
					ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
				},
				ExpectedStatusCode: http.StatusBadRequest,
			},
			{
				RequestBody: settings.HandleUpdateServerSettingsRequestBody{
					IsSignUpOpen:                &signupOpen,
					UsersMustBeApproved:         &usersMustBeApproved,
					AllowedEmailRegex:           &allowedEmailRegex,
					BlockedEmailRegex:           &blockedEmailRegex,
					ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
				},
				ExpectedStatusCode: http.StatusBadRequest,
			},
			{
				RequestBody: settings.HandleUpdateServerSettingsRequestBody{
					IsSignUpRestricted:          &signupRestricted,
					UsersMustBeApproved:         &usersMustBeApproved,
					AllowedEmailRegex:           &allowedEmailRegex,
					BlockedEmailRegex:           &blockedEmailRegex,
					ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
				},
				ExpectedStatusCode: http.StatusBadRequest,
			},
			{
				RequestBody: settings.HandleUpdateServerSettingsRequestBody{
					IsSignUpOpen:                &signupOpen,
					IsSignUpRestricted:          &signupRestricted,
					UsersMustBeApproved:         &usersMustBeApproved,
					AllowedEmailRegex:           &allowedEmailRegex,
					BlockedEmailRegex:           &blockedEmailRegex,
					ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
				},
				ExpectedStatusCode: http.StatusOK,
			},
		}

		for _, tc := range testCases {
			w := httptest.NewRecorder()
			req := testutils.CreateRequestWithJSONBody(
				t,
				"/api/v1/admin/authentication-settings",
				"PUT",
				tc.RequestBody,
			)
			testutils.AuthenticateHttpRequest(t, req, *adminUser)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.ExpectedStatusCode, w.Code)
		}
	})
}

/*
Test endpoint to update authentication settings with
a user that is not a superuser
*/
func TestUpdateAuthenticationSettingsCommonUser(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		commonUser, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil {
			t.Error(err)
		}

		signupOpen := true
		signupRestricted := true
		usersMustBeApproved := true
		allowedEmailRegex := `^.*@example\.com$`
		blockedEmailRegex := `^.*@blocked\.com$`
		approvedByDefaultEmailRegex := `^.*@example\.com$`

		requestBody := settings.HandleUpdateServerSettingsRequestBody{
			IsSignUpOpen:                &signupOpen,
			IsSignUpRestricted:          &signupRestricted,
			UsersMustBeApproved:         &usersMustBeApproved,
			AllowedEmailRegex:           &allowedEmailRegex,
			BlockedEmailRegex:           &blockedEmailRegex,
			ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/admin/authentication-settings",
			"PUT",
			requestBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *commonUser)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

/*
Test endpoint to update authentication settings without
authentication
*/
func TestUpdateAuthenticationSettingsWithoutAuthentication(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		signupOpen := true
		signupRestricted := true
		usersMustBeApproved := true
		allowedEmailRegex := `^.*@example\.com$`
		blockedEmailRegex := `^.*@blocked\.com$`
		approvedByDefaultEmailRegex := `^.*@example\.com$`

		requestBody := settings.HandleUpdateServerSettingsRequestBody{
			IsSignUpOpen:                &signupOpen,
			IsSignUpRestricted:          &signupRestricted,
			UsersMustBeApproved:         &usersMustBeApproved,
			AllowedEmailRegex:           &allowedEmailRegex,
			BlockedEmailRegex:           &blockedEmailRegex,
			ApprovedByDefaultEmailRegex: &approvedByDefaultEmailRegex,
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/admin/authentication-settings",
			"PUT",
			requestBody,
		)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
