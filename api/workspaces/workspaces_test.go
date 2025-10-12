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
Try to create a workspace from git source
*/
func TestCreateWorkspaceFromGitSource(t *testing.T) {
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

		// retrieve the id of the first runner
		runners, err := models.ListRunners(1, 0)
		if err != nil || len(runners) == 0 {
			t.Errorf("Failed to retrieve test runner: '%s'\n", err)
			t.FailNow()
		}

		// create a new workspace
		reqBody := workspaces.CreateWorkspaceRequestBody{
			Name:                 "Test Workspace",
			Type:                 "docker_compose",
			RunnerID:             runners[0].ID,
			ConfigSource:         models.WorkspaceConfigSourceGit,
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

/*
Try to create a workspace from git source
*/
func TestCreateWorkspaceFromTemplateSource(t *testing.T) {
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

		// retrieve the id of the first runner
		runners, err := models.ListRunners(1, 0)
		if err != nil || len(runners) == 0 {
			t.Errorf("Failed to retrieve test runner: '%s'\n", err)
			t.FailNow()
		}

		// create a template and a template version
		template, err := models.CreateWorkspaceTemplate("Test Template", "docker_compose", "", "")
		if err != nil {
			t.Errorf("Failed to create template: '%s'\n", err)
			t.FailNow()
		}

		templateVersion, err := models.CreateTemplateVersion(*template, "v1.0.0", *user, "docker-compose.yml")
		if err != nil {
			t.Errorf("Failed to create template version: '%s'\n", err)
			t.FailNow()
		}

		// create a new workspace
		reqBody := workspaces.CreateWorkspaceRequestBody{
			Name:                 "Test Workspace",
			Type:                 "docker_compose",
			RunnerID:             runners[0].ID,
			ConfigSource:         models.WorkspaceConfigSourceTemplate,
			TemplateVersionID:    templateVersion.ID,
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

/*
Try to create a workspace with invalid parameters
*/
func TestCreateWorkspaceErrors(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := api.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		runners, err := models.ListRunners(1, 0)
		if err != nil || len(runners) == 0 {
			t.Fatalf("Failed to retrieve test runner: '%s'", err)
		}

		baseReq := workspaces.CreateWorkspaceRequestBody{
			Name:                 "Test Workspace",
			Type:                 "docker_compose",
			RunnerID:             runners[0].ID,
			ConfigSource:         models.WorkspaceConfigSourceGit,
			TemplateVersionID:    0,
			GitRepoUrl:           "https://github.com/davidebianchi03/codebox.git",
			GitRefName:           "main",
			ConfigSourceFilePath: "/path/to/config",
			EnvironmentVariables: []string{"VAR1=value1", "VAR2=value2"},
		}

		tests := []struct {
			name       string
			modifyBody func(b *workspaces.CreateWorkspaceRequestBody)
			wantCode   int
		}{
			{
				name: "missing name",
				modifyBody: func(b *workspaces.CreateWorkspaceRequestBody) {
					b.Name = ""
				},
				wantCode: http.StatusBadRequest,
			},
			{
				name: "invalid type",
				modifyBody: func(b *workspaces.CreateWorkspaceRequestBody) {
					b.Type = "unknown_type"
				},
				wantCode: http.StatusBadRequest,
			},
			{
				name: "invalid runner id",
				modifyBody: func(b *workspaces.CreateWorkspaceRequestBody) {
					b.RunnerID = 999999
				},
				wantCode: http.StatusBadRequest,
			},
			{
				name: "invalid config source",
				modifyBody: func(b *workspaces.CreateWorkspaceRequestBody) {
					b.ConfigSource = "invalid_source"
				},
				wantCode: http.StatusBadRequest,
			},
			{
				name: "create workspace with invalid template version id",
				modifyBody: func(b *workspaces.CreateWorkspaceRequestBody) {
					b.ConfigSource = models.WorkspaceConfigSourceTemplate
					b.TemplateVersionID = 999999
				},
				wantCode: http.StatusBadRequest,
			},
			{
				name: "create workspace with git source but missing git repo url",
				modifyBody: func(b *workspaces.CreateWorkspaceRequestBody) {
					b.ConfigSource = models.WorkspaceConfigSourceGit
					b.GitRepoUrl = ""
				},
				wantCode: http.StatusBadRequest,
			},
			{
				name: "create workspace with git source but missing config source file path",
				modifyBody: func(b *workspaces.CreateWorkspaceRequestBody) {
					b.ConfigSource = models.WorkspaceConfigSourceGit
					b.ConfigSourceFilePath = ""
				},
				wantCode: http.StatusBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body := baseReq
				tt.modifyBody(&body)

				w := httptest.NewRecorder()
				req := testutils.CreateRequestWithJSONBody(
					t,
					"/api/v1/workspace",
					"POST",
					body,
				)
				testutils.AuthenticateHttpRequest(t, req, *user)
				router.ServeHTTP(w, req)

				if w.Code != tt.wantCode {
					t.Errorf("[%s] expected status %d, got %d\nBody: %s",
						tt.name, tt.wantCode, w.Code, w.Body.String())
				}
			})
		}
	})
}
