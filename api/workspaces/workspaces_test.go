package workspaces_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/api"
	"gitlab.com/codebox4073715/codebox/api/serializers"
	"gitlab.com/codebox4073715/codebox/api/workspaces"
	"gitlab.com/codebox4073715/codebox/testutils"

	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
Try to create a workspace
*/
func TestCreateWorkspace_Debug(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := api.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Errorf("Failed to retrieve test user: '%s'\n", err)
			t.FailNow()
		}

		// list existing workspaces (should be none)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/workspace", nil)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		data, _ := serializers.MultipleWorkspaceSerializersFromJSON(w.Body.String())
		assert.Equal(t, 0, len(data))

		// create a new workspace
		reqBody := workspaces.CreateWorkspaceRequestBody{
			Name:                 "Test Workspace",
			Type:                 "docker_compose",
			RunnerID:             1,
			ConfigSource:         models.WorkspaceConfigSourceGit,
			TemplateVersionID:    0,
			GitRepoUrl:           "https://github.com/davidebianchi03/codebox.git",
			GitRefName:           "main",
			ConfigSourceFilePath: "/path/to/config",
			EnvironmentVariables: []string{"VAR1=value1", "VAR2=value2"},
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// check that the workspace was created
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/workspace", nil)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		data, _ = serializers.MultipleWorkspaceSerializersFromJSON(w.Body.String())
		assert.Equal(t, 1, len(data))
	})
}
