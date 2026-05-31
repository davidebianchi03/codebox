package commands

import (
	"fmt"
	"log"

	"gitlab.com/codebox4073715/codebox/cli/args"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
This function handles the command to verify an email
*/
func HandleVerifyEmail(args args.VerifyEmailCmdArgs) uint {
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

	user, err := models.RetrieveUserByEmail(args.Email)
	if err != nil {
		fmt.Println("Failed to retrieve user by email, unknown error")
		return 1
	}

	if user == nil {
		fmt.Println("No user found with the given email")
		return 1
	}

	if user.EmailVerified {
		fmt.Println("Nothing to do here, email was already verified")
	} else {
		user.EmailVerified = true
		if err := models.UpdateUser(user); err != nil {
			fmt.Println("Failed to update user, unknown error")
			return 1
		}

		fmt.Println("Email has been verified")
	}

	return 0
}
