package testutils

import "gitlab.com/codebox4073715/codebox/db/models"

// SetupDBForTests initializes the database with necessary data for tests
func SetupDBForTests() error {
	// create an admin user
	if _, err := models.CreateUser(
		"admin@admin.com",
		"Admin",
		"User",
		"password",
		true,
		true,
	); err != nil {
		return err
	}

	// create a regular user
	if _, err := models.CreateUser(
		"user1@user.com",
		"User1",
		"User",
		"password",
		false,
		false,
	); err != nil {
		return err
	}

	// create another regular user
	if _, err := models.CreateUser(
		"user2@user.com",
		"User2",
		"User",
		"password",
		false,
		false,
	); err != nil {
		return err
	}

	// create a runner for docker based workspaces
	if _, err := models.CreateRunner(
		"docker-runner",
		"docker",
		false,
		"",
	); err != nil {
		return err
	}

	return nil
}
