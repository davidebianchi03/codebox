package commands

import (
	"fmt"
	"log"

	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
Handle command to reset the password for a user
This is an interactive command
*/
func HandleSetPassword() uint {
	// load config from env vars
	err := config.InitCodeBoxEnv()
	if err != nil {
		log.Fatalf("Failed to load server configuration from environment: '%s'\n", err)
		return 1
	}

	// init db connection
	if err = dbconn.ConnectDB(); err != nil {
		log.Fatalf("Cannot init connection with DB: '%s'\n", err)
		return 1
	}

	terminal := PrepareTerminal()
	fmt.Print("Enter email: ")
	email, err := terminal.ReadLine()
	if err != nil {
		return 1
	}

	password, err := terminal.ReadPassword("New password:")
	if err != nil {
		return 1
	}

	passwordConfirm, err := terminal.ReadPassword("Confirm the password")
	if err != nil {
		return 1
	}

	if password != passwordConfirm {
		fmt.Println("Passwords do not match")
		return 1
	}

	user, err := models.RetrieveUserByEmail(email)
	if err != nil {
		fmt.Println("Failed to retrieve user")
		return 1
	}

	if user == nil {
		fmt.Println("User not found")
		return 1
	}

	if err := models.ValidatePassword(password); err != nil {
		fmt.Println(err.Error())
		return 1
	}

	password, err = models.HashPassword(password)
	if err != nil {
		fmt.Println("Unknown error")
		return 1
	}
	user.Password = password

	if err := models.UpdateUser(user); err != nil {
		fmt.Printf("Fail to update user, %s", err)
		return 1
	}

	return 0
}
