package testutils

import (
	"testing"

	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
)

// func SetupTestEnvironment(t *testing.T) error {
// 	// err := config.InitCodeBoxEnv()
// 	// if err != nil {
// 	// 	t.Errorf("Failed to load server configuration from environment: '%s'\n", err)
// 	// 	return err
// 	// }

// 	// // test della connessione con il database
// 	// err = dbconn.ConnectDB()
// 	// if err != nil {
// 	// 	t.Errorf("Cannot init connection with DB: '%s'\n", err)
// 	// 	return err
// 	// }

// 	// // clear the contents of the database
// 	// if err := ClearDBTables(); err != nil {
// 	// 	t.Errorf("Cannot clear DB: '%s'\n", err)
// 	// 	return err
// 	// }

// 	// // mock bg tasks
// 	// mock := &MockEnqueuer{}
// 	// bgtasks.BgTasksEnqueuer = mock

// 	return nil
// }

// func TeardownTestEnvironment(t *testing.T) error {
// 	// clear the contents of the database
// 	if err := ClearDBTables(); err != nil {
// 		t.Errorf("Cannot clear DB: '%s'\n", err)
// 		return err
// 	}
// 	dbconn.CloseDB()

// 	return nil
// }

// Helper function that wraps a test function with
// setup and teardown of the test environment
// The setup is done before calling the test function, it initializes the db and
// clears its contents
// The teardown is done after the test function returns, it clears the db contents
// and closes the db connection
// If setup or teardown fail, the test will fail immediately
func WithSetupAndTearDownTestEnvironment(t *testing.T, testFunc func(t *testing.T)) {
	// if err := SetupTestEnvironment(t); err != nil {
	// 	t.FailNow()
	// 	return
	// }
	err := config.InitCodeBoxEnv()
	if err != nil {
		t.Errorf("Failed to load server configuration from environment: '%s'\n", err)
		t.FailNow()
		return
	}

	err = dbconn.ConnectDB()
	if err != nil {
		t.Errorf("Cannot init connection with DB: '%s'\n", err)
		t.FailNow()
		return
	}

	if err := ClearDB(dbconn.DB); err != nil {
		t.Errorf("Cannot clear DB: '%s'\n", err)
		t.FailNow()
		return
	}

	// transaction := dbconn.DB.Begin()
	// if transaction.Error != nil {
	// 	t.Errorf("Cannot start transaction: '%s'\n", transaction.Error)
	// 	t.FailNow()
	// 	return
	// }

	// t.Cleanup(func() {
	// 	transaction.Rollback()
	// })

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
