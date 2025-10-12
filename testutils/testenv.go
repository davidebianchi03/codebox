package testutils

import (
	"testing"

	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
)

func SetupTestEnvironment(t *testing.T) error {
	err := config.InitCodeBoxEnv()
	if err != nil {
		t.Errorf("Failed to load server configuration from environment: '%s'\n", err)
		return err
	}

	// test della connessione con il database
	err = dbconn.ConnectDB()
	if err != nil {
		t.Errorf("Cannot init connection with DB: '%s'\n", err)
		return err
	}

	// clear the contents of the database
	if err := ClearDBTables(); err != nil {
		t.Errorf("Cannot clear DB: '%s'\n", err)
		return err
	}

	// mock bg tasks
	mock := &MockEnqueuer{}
	bgtasks.BgTasksEnqueuer = mock

	return nil
}

func TeardownTestEnvironment(t *testing.T) error {
	// clear the contents of the database
	if err := ClearDBTables(); err != nil {
		t.Errorf("Cannot clear DB: '%s'\n", err)
		return err
	}
	dbconn.CloseDB()

	return nil
}

/*
Helper function that wraps a test function with setup and teardown of the test environment
*/
func WithSetupAndTearDownTestEnvironment(t *testing.T, testFunc func(t *testing.T)) {
	if err := SetupTestEnvironment(t); err != nil {
		t.FailNow()
	}

	if err := SetupDBForTests(); err != nil {
		t.Errorf("Cannot setup DB for tests: '%s'\n", err)
		t.FailNow()
	}

	defer func() {
		// if err := TeardownTestEnvironment(t); err != nil {
		// 	t.FailNow()
		// }
	}()
	testFunc(t)
}
