package workspaces_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/httpserver"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/workspaces"
	"gitlab.com/codebox4073715/codebox/testutils"

	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
Try to create a workspace from git source
*/
func TestCreateWorkspaceFromGitSource(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

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
		assert.Equal(t, data[0].Name, "Test Workspace")
		assert.Equal(t, data[0].Type, "docker_compose")
		assert.Equal(t, data[0].Status, models.WorkspaceStatusStarting)
		assert.Equal(t, data[0].Runner.ID, runners[0].ID)
		assert.Equal(t, data[0].ConfigSource, models.WorkspaceConfigSourceGit)
		assert.Equal(t, data[0].GitSource.RepositoryURL, "https://github.com/davidebianchi03/codebox.git")
	})
}

/*
Try to create a workspace from git source
*/
func TestCreateWorkspaceFromTemplateSource(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

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
		assert.Equal(t, 1, len(data))
		assert.Equal(t, data[0].Name, "Test Workspace")
		assert.Equal(t, data[0].Type, "docker_compose")
		assert.Equal(t, data[0].Status, models.WorkspaceStatusStarting)
		assert.Equal(t, data[0].Runner.ID, runners[0].ID)
		assert.Equal(t, data[0].ConfigSource, models.WorkspaceConfigSourceTemplate)
	})
}

/*
Try to create a workspace with invalid parameters
*/
func TestCreateWorkspaceErrors(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

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

/*
Try to create a workspace without authentication
*/
func TestCreateWorkspaceWithoutAuthentication(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		runners, err := models.ListRunners(1, 0)
		if err != nil || len(runners) == 0 {
			t.Errorf("Failed to retrieve test runner: '%s'\n", err)
			t.FailNow()
		}

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

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		// Note: no authentication
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

/*
Try to update a workspace, both when is is running and when it is stopped
*/
func TestUpdateWorkspace(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

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
			EnvironmentVariables: []string{"VAR1=value1"},
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// parse the created workspace
		createdWorkspace, err := serializers.WorkspaceSerializerFromJSON(w.Body.String())
		if err != nil {
			t.Fatalf("Failed to parse created workspace: '%s'", err)
		}

		// update the workspace while it is running
		updateReqBody := workspaces.UpdateWorkspaceRequestBody{
			EnvironmentVariables: []string{"VAR1=newvalue", "VAR3=value3"},
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d", createdWorkspace.ID),
			"PUT",
			updateReqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotAcceptable, w.Code)

		// mark the workspace as stopped
		workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
		if err != nil || workspace == nil {
			t.Fatalf("Failed to retrieve workspace: '%s'", err)
		}

		if _, err := models.UpdateWorkspace(
			workspace,
			workspace.Name,
			models.WorkspaceStatusStopped,
			workspace.Runner,
			workspace.ConfigSource,
			workspace.TemplateVersion,
			workspace.GitSource,
			workspace.EnvironmentVariables,
		); err != nil {
			t.Fatalf("Failed to stop workspace: '%s'", err)
		}

		// try to update the workspace again
		updateReqBody = workspaces.UpdateWorkspaceRequestBody{
			GitRepoUrl:           "https://git.example.com/new/repo.git",
			EnvironmentVariables: []string{"VAR1=newvalue", "VAR3=value3"},
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d", createdWorkspace.ID),
			"PUT",
			updateReqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// check that the workspace was updated
		workspace, err = models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
		if err != nil || workspace == nil {
			t.Fatalf("Failed to retrieve updated workspace: '%s'", err)
		}

		assert.Equal(t, workspace.GitSource.RepositoryURL, "https://git.example.com/new/repo.git")
		assert.Equal(t, len(workspace.EnvironmentVariables), 2)
	})
}

// Helper struct for start/stop workspace tests
type WorkspaceStatusTestCase struct {
	WorkspaceStatus string
	ExpectedCode    int
}

/*
Try to start a workspace
*/
func TestStartWorkspace(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

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
			EnvironmentVariables: []string{"VAR1=value1"},
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// parse the created workspace
		createdWorkspace, err := serializers.WorkspaceSerializerFromJSON(w.Body.String())
		if err != nil {
			t.Fatalf("Failed to parse created workspace: '%s'", err)
		}

		testCases := []WorkspaceStatusTestCase{
			{models.WorkspaceStatusStarting, http.StatusConflict},
			{models.WorkspaceStatusRunning, http.StatusConflict},
			{models.WorkspaceStatusStopping, http.StatusConflict},
			{models.WorkspaceStatusStopped, http.StatusOK},
			{models.WorkspaceStatusError, http.StatusConflict},
			{models.WorkspaceStatusDeleting, http.StatusConflict},
		}

		for _, tc := range testCases {
			t.Run(tc.WorkspaceStatus, func(t *testing.T) {
				// set the status of the workspace
				workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
				if err != nil || workspace == nil {
					t.Fatalf("Failed to retrieve workspace: '%s'", err)
				}

				if _, err := models.UpdateWorkspace(
					workspace,
					workspace.Name,
					tc.WorkspaceStatus,
					workspace.Runner,
					workspace.ConfigSource,
					workspace.TemplateVersion,
					workspace.GitSource,
					workspace.EnvironmentVariables,
				); err != nil {
					t.Fatalf("Failed to stop workspace: '%s'", err)
				}

				// try to start the workspace
				w = httptest.NewRecorder()
				req = testutils.CreateRequestWithJSONBody(
					t,
					fmt.Sprintf("/api/v1/workspace/%d/start", createdWorkspace.ID),
					"POST",
					nil,
				)
				testutils.AuthenticateHttpRequest(t, req, *user)
				router.ServeHTTP(w, req)
				assert.Equal(t, tc.ExpectedCode, w.Code)
			})
		}
	})
}

/*
Try to start a workspace that has no runner assigned
*/
func TestStartWorkspaceNoRunner(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

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
			EnvironmentVariables: []string{"VAR1=value1"},
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// parse the created workspace
		createdWorkspace, err := serializers.WorkspaceSerializerFromJSON(w.Body.String())
		if err != nil {
			t.Fatalf("Failed to parse created workspace: '%s'", err)
		}

		// set the status of the workspace and remove the runner
		workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
		if err != nil || workspace == nil {
			t.Fatalf("Failed to retrieve workspace: '%s'", err)
		}

		if _, err := models.UpdateWorkspace(
			workspace,
			workspace.Name,
			models.WorkspaceStatusStopped,
			nil,
			workspace.ConfigSource,
			workspace.TemplateVersion,
			workspace.GitSource,
			workspace.EnvironmentVariables,
		); err != nil {
			t.Fatalf("Failed to stop workspace: '%s'", err)
		}

		// try to start the workspace with no runner assigned
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/start", createdWorkspace.ID),
			"POST",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusFailedDependency, w.Code)

		// try to start the workspace with another user
		user2, err := models.RetrieveUserByEmail("user2@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/start", createdWorkspace.ID),
			"POST",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user2)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// try to start the workspace without authentication
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/start", createdWorkspace.ID),
			"POST",
			nil,
		)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// try to start a workspace that does not exist
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace/104/start",
			"POST",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

/*
Try to stop a workspace
*/
func TestStopWorkspace(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

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
			EnvironmentVariables: []string{"VAR1=value1"},
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// parse the created workspace
		createdWorkspace, err := serializers.WorkspaceSerializerFromJSON(w.Body.String())
		if err != nil {
			t.Fatalf("Failed to parse created workspace: '%s'", err)
		}

		testCases := []WorkspaceStatusTestCase{
			{models.WorkspaceStatusStarting, http.StatusConflict},
			{models.WorkspaceStatusRunning, http.StatusOK},
			{models.WorkspaceStatusStopping, http.StatusConflict},
			{models.WorkspaceStatusStopped, http.StatusConflict},
			{models.WorkspaceStatusError, http.StatusOK},
			{models.WorkspaceStatusDeleting, http.StatusConflict},
		}

		for _, tc := range testCases {
			t.Run(tc.WorkspaceStatus, func(t *testing.T) {
				// set the status of the workspace
				workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
				if err != nil || workspace == nil {
					t.Fatalf("Failed to retrieve workspace: '%s'", err)
				}

				if _, err := models.UpdateWorkspace(
					workspace,
					workspace.Name,
					tc.WorkspaceStatus,
					workspace.Runner,
					workspace.ConfigSource,
					workspace.TemplateVersion,
					workspace.GitSource,
					workspace.EnvironmentVariables,
				); err != nil {
					t.Fatalf("Failed to stop workspace: '%s'", err)
				}

				// try to start the workspace
				w = httptest.NewRecorder()
				req = testutils.CreateRequestWithJSONBody(
					t,
					fmt.Sprintf("/api/v1/workspace/%d/stop", createdWorkspace.ID),
					"POST",
					nil,
				)
				testutils.AuthenticateHttpRequest(t, req, *user)
				router.ServeHTTP(w, req)
				assert.Equal(t, tc.ExpectedCode, w.Code)
			})
		}

		// try to stop the workspace with another user
		user2, err := models.RetrieveUserByEmail("user2@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/stop", createdWorkspace.ID),
			"POST",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user2)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// try to stop the workspace without authentication
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/stop", createdWorkspace.ID),
			"POST",
			nil,
		)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// try to stop a workspace that does not exist
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace/104/stop",
			"POST",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdateConfigOfAWorkspace(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

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
			EnvironmentVariables: []string{"VAR1=value1"},
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// parse the created workspace
		createdWorkspace, err := serializers.WorkspaceSerializerFromJSON(w.Body.String())
		if err != nil {
			t.Fatalf("Failed to parse created workspace: '%s'", err)
		}

		testCases := []WorkspaceStatusTestCase{
			{models.WorkspaceStatusStarting, http.StatusNotAcceptable},
			{models.WorkspaceStatusRunning, http.StatusNotAcceptable},
			{models.WorkspaceStatusStopping, http.StatusNotAcceptable},
			{models.WorkspaceStatusStopped, http.StatusOK},
			{models.WorkspaceStatusError, http.StatusNotAcceptable},
			{models.WorkspaceStatusDeleting, http.StatusNotAcceptable},
		}

		for _, tc := range testCases {
			t.Run(tc.WorkspaceStatus, func(t *testing.T) {
				// set the status of the workspace
				workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
				if err != nil || workspace == nil {
					t.Fatalf("Failed to retrieve workspace: '%s'", err)
				}

				if _, err := models.UpdateWorkspace(
					workspace,
					workspace.Name,
					tc.WorkspaceStatus,
					workspace.Runner,
					workspace.ConfigSource,
					workspace.TemplateVersion,
					workspace.GitSource,
					workspace.EnvironmentVariables,
				); err != nil {
					t.Fatalf("Failed to update workspace: '%s'", err)
				}

				// try to start the workspace
				w = httptest.NewRecorder()
				req = testutils.CreateRequestWithJSONBody(
					t,
					fmt.Sprintf("/api/v1/workspace/%d/update-config", createdWorkspace.ID),
					"POST",
					nil,
				)
				testutils.AuthenticateHttpRequest(t, req, *user)
				router.ServeHTTP(w, req)
				assert.Equal(t, tc.ExpectedCode, w.Code)
			})
		}

		// try to update the config of the workspace with another user
		user2, err := models.RetrieveUserByEmail("user2@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/update-config", createdWorkspace.ID),
			"POST",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user2)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// try to update the config of the workspace without authentication
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/update-config", createdWorkspace.ID),
			"POST",
			nil,
		)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// try to update the config of a workspace that has no runner assigned
		workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
		if err != nil || workspace == nil {
			t.Fatalf("Failed to retrieve workspace: '%s'", err)
		}

		if _, err := models.UpdateWorkspace(
			workspace,
			workspace.Name,
			workspace.Status,
			nil,
			workspace.ConfigSource,
			workspace.TemplateVersion,
			workspace.GitSource,
			workspace.EnvironmentVariables,
		); err != nil {
			t.Fatalf("Failed to update workspace: '%s'", err)
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/update-config", createdWorkspace.ID),
			"POST",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotAcceptable, w.Code)

		// try to update the config of a workspace that does not exist
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace/104/update-config",
			"POST",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

type DeleteWorkspaceTestCase struct {
	WorkspaceStatus string
	ExpectedCode    int
	SkipErrors      bool
}

func TestDeleteWorkspace(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		runners, err := models.ListRunners(1, 0)
		if err != nil || len(runners) == 0 {
			t.Errorf("Failed to retrieve test runner: '%s'\n", err)
			t.FailNow()
		}

		testCases := []DeleteWorkspaceTestCase{
			{models.WorkspaceStatusStarting, http.StatusNotAcceptable, false},
			{models.WorkspaceStatusRunning, http.StatusOK, false},
			{models.WorkspaceStatusStopping, http.StatusNotAcceptable, false},
			{models.WorkspaceStatusStopped, http.StatusOK, false},
			{models.WorkspaceStatusError, http.StatusOK, false},
			{models.WorkspaceStatusError, http.StatusOK, true},
			{models.WorkspaceStatusDeleting, http.StatusNotAcceptable, false},
		}

		for _, tc := range testCases {
			t.Run(tc.WorkspaceStatus, func(t *testing.T) {
				// create a new workspace
				reqBody := workspaces.CreateWorkspaceRequestBody{
					Name:                 "Test Workspace",
					Type:                 "docker_compose",
					RunnerID:             runners[0].ID,
					ConfigSource:         models.WorkspaceConfigSourceGit,
					GitRepoUrl:           "https://github.com/davidebianchi03/codebox.git",
					GitRefName:           "main",
					ConfigSourceFilePath: "/path/to/config",
					EnvironmentVariables: []string{"VAR1=value1"},
				}

				w := httptest.NewRecorder()
				req := testutils.CreateRequestWithJSONBody(
					t,
					"/api/v1/workspace",
					"POST",
					reqBody,
				)
				testutils.AuthenticateHttpRequest(t, req, *user)
				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusCreated, w.Code)

				// parse the created workspace
				createdWorkspace, err := serializers.WorkspaceSerializerFromJSON(w.Body.String())
				if err != nil {
					t.Fatalf("Failed to parse created workspace: '%s'", err)
				}

				// set the status of the workspace
				workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
				if err != nil || workspace == nil {
					t.Fatalf("Failed to retrieve workspace: '%s'", err)
				}

				if _, err := models.UpdateWorkspace(
					workspace,
					workspace.Name,
					tc.WorkspaceStatus,
					workspace.Runner,
					workspace.ConfigSource,
					workspace.TemplateVersion,
					workspace.GitSource,
					workspace.EnvironmentVariables,
				); err != nil {
					t.Fatalf("Failed to stop workspace: '%s'", err)
				}

				// try to start the workspace
				w = httptest.NewRecorder()
				req = testutils.CreateRequestWithJSONBody(
					t,
					fmt.Sprintf("/api/v1/workspace/%d?skip_errors=%t", createdWorkspace.ID, tc.SkipErrors),
					"DELETE",
					nil,
				)
				testutils.AuthenticateHttpRequest(t, req, *user)
				router.ServeHTTP(w, req)
				assert.Equal(t, tc.ExpectedCode, w.Code)
			})
		}

		reqBody := workspaces.CreateWorkspaceRequestBody{
			Name:                 "Test Workspace",
			Type:                 "docker_compose",
			RunnerID:             runners[0].ID,
			ConfigSource:         models.WorkspaceConfigSourceGit,
			GitRepoUrl:           "https://github.com/davidebianchi03/codebox.git",
			GitRefName:           "main",
			ConfigSourceFilePath: "/path/to/config",
			EnvironmentVariables: []string{"VAR1=value1"},
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// parse the created workspace
		createdWorkspace, err := serializers.WorkspaceSerializerFromJSON(w.Body.String())
		if err != nil {
			t.Fatalf("Failed to parse created workspace: '%s'", err)
		}

		// try to delete the workspace with another user
		user2, err := models.RetrieveUserByEmail("user2@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d", createdWorkspace.ID),
			"DELETE",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user2)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// try to delete the workspace without authentication
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d", createdWorkspace.ID),
			"DELETE",
			nil,
		)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// try to delete a workspace that does not exist
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace/104",
			"DELETE",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// try to delete a workspace that has no runner assigned
		workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
		if err != nil {
			t.Fatalf("failed to retrieve workspace %s", err)
			return
		}

		if workspace == nil {
			t.Fatal("workspace not found")
			return
		}

		if _, err := models.UpdateWorkspace(
			workspace,
			workspace.Name,
			models.WorkspaceStatusRunning,
			nil,
			workspace.ConfigSource,
			workspace.TemplateVersion,
			workspace.GitSource,
			workspace.EnvironmentVariables,
		); err != nil {
			t.Fatalf("failed to unassign workspace runner %s", err)
			return
		}

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d", workspace.ID),
			"DELETE",
			nil,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

/*
Test the API that sets the runner of a workspace
*/
func TestSetWorkspaceRunner(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}

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
			EnvironmentVariables: []string{"VAR1=value1"},
		}

		w := httptest.NewRecorder()
		req := testutils.CreateRequestWithJSONBody(
			t,
			"/api/v1/workspace",
			"POST",
			reqBody,
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// parse the created workspace
		createdWorkspace, err := serializers.WorkspaceSerializerFromJSON(w.Body.String())
		if err != nil {
			t.Fatalf("Failed to parse created workspace: '%s'", err)
		}

		// set the status of the workspace
		workspace, err := models.RetrieveWorkspaceByUserAndId(*user, createdWorkspace.ID)
		if err != nil || workspace == nil {
			t.Fatalf("Failed to retrieve workspace: '%s'", err)
		}

		// mark the workspace as stopped and remove the runner
		if _, err := models.UpdateWorkspace(
			workspace,
			workspace.Name,
			models.WorkspaceStatusStopped,
			nil,
			workspace.ConfigSource,
			workspace.TemplateVersion,
			workspace.GitSource,
			workspace.EnvironmentVariables,
		); err != nil {
			t.Fatalf("Failed to update workspace: '%s'", err)
		}

		// try to set the runner of the workspace (should succeed)
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/set-runner", createdWorkspace.ID),
			"POST",
			workspaces.SetRunnerForWorkspaceBody{
				RunnerId: runners[0].ID,
			},
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// try to set again the runner of the workspace (use the same runner, should succeed)
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/set-runner", createdWorkspace.ID),
			"POST",
			workspaces.SetRunnerForWorkspaceBody{
				RunnerId: runners[0].ID,
			},
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// create a second runner
		runner2, err := models.CreateRunner(
			"docker-runner-2",
			"docker",
			false,
			"",
		)
		if err != nil {
			t.Fatalf("Failed to create second runner: '%s'", err)
		}

		// try to set again the runner of the workspace (use another runner, should fail)
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/set-runner", createdWorkspace.ID),
			"POST",
			workspaces.SetRunnerForWorkspaceBody{
				RunnerId: runner2.ID,
			},
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// try to set the runner of the workspace, use id of a
		// runner that does not exist (try twice)
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/set-runner", createdWorkspace.ID),
			"POST",
			workspaces.SetRunnerForWorkspaceBody{
				RunnerId: 0,
			},
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/set-runner", createdWorkspace.ID),
			"POST",
			workspaces.SetRunnerForWorkspaceBody{
				RunnerId: 104,
			},
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// set the runner of a workspace that does not exist
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/set-runner", 104),
			"POST",
			workspaces.SetRunnerForWorkspaceBody{
				RunnerId: runners[0].ID,
			},
		)
		testutils.AuthenticateHttpRequest(t, req, *user)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// set the runner of a workspace with another user
		user2, err := models.RetrieveUserByEmail("user2@user.com")
		if err != nil || user == nil {
			t.Fatalf("Failed to retrieve test user: '%s'", err)
		}
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/set-runner", createdWorkspace.ID),
			"POST",
			workspaces.SetRunnerForWorkspaceBody{
				RunnerId: runners[0].ID,
			},
		)
		testutils.AuthenticateHttpRequest(t, req, *user2)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// set the runner of a workspace without authentication
		w = httptest.NewRecorder()
		req = testutils.CreateRequestWithJSONBody(
			t,
			fmt.Sprintf("/api/v1/workspace/%d/set-runner", createdWorkspace.ID),
			"POST",
			workspaces.SetRunnerForWorkspaceBody{
				RunnerId: runners[0].ID,
			},
		)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
