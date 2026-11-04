package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver"
	"gitlab.com/codebox4073715/codebox/testutils"
)

func TestAdminSetUserPassword(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		adminUser, err := models.RetrieveUserByEmail("admin@admin.com")
		require.NoError(t, err)
		require.NotNil(t, adminUser)

		adminPasswordHash := adminUser.Password
		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/admin/users/admin@admin.com/set-password",
			http.MethodPost,
			map[string]string{"password": "NewStrongPass!"},
		)
		testutils.AuthenticateHttpRequest(t, req, *adminUser)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)

		adminUser, err = models.RetrieveUserByEmail("admin@admin.com")
		require.NoError(t, err)
		require.NotNil(t, adminUser)
		require.Equal(t, adminPasswordHash, adminUser.Password)

		regularUser, err := models.RetrieveUserByEmail("user1@user.com")
		require.NoError(t, err)
		require.NotNil(t, regularUser)

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/admin/users/user1@user.com/set-password",
			http.MethodPost,
			map[string]string{"password": "NewStrongPass!"},
		)
		testutils.AuthenticateHttpRequest(t, req, *adminUser)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		regularUser, err = models.RetrieveUserByEmail("user1@user.com")
		require.NoError(t, err)
		require.NotNil(t, regularUser)
		require.True(t, regularUser.CheckPassword("NewStrongPass!"))
	})
}
