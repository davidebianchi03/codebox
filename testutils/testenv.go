package testutils

import (
	"testing"

	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
)

// Helper function that wraps a test function with
// setup and teardown of the test environment
// The setup is done before calling the test function, it initializes the db and
// clears its contents
// The teardown is done after the test function returns, it clears the db contents
// and closes the db connection
// If setup or teardown fail, the test will fail immediately
func WithSetupAndTearDownTestEnvironment(t *testing.T, testFunc func(t *testing.T)) {
	// load config
	if err := config.InitCodeBoxEnv(); err != nil {
		t.Errorf("Failed to load server configuration from environment: '%s'\n", err)
		t.FailNow()
		return
	}

	// clear cache
	if err := ClearCache(); err != nil {
		t.Errorf("Cannot clear cache for tests: '%s'\n", err)
		t.FailNow()
		return
	}

	// setup db
	if err := dbconn.ConnectDB(); err != nil {
		t.Errorf("Cannot init connection with DB: '%s'\n", err)
		t.FailNow()
		return
	}

	if err := ClearDB(dbconn.DB); err != nil {
		t.Errorf("Cannot clear DB: '%s'\n", err)
		t.FailNow()
		return
	}

	if err := SetupDBForTests(); err != nil {
		t.Errorf("Cannot setup DB for tests: '%s'\n", err)
		t.FailNow()
		return
	}

	// mock bg tasks
	mock := &MockEnqueuer{}
	bgtasks.BgTasksEnqueuer = mock

	testFunc(t)
}
