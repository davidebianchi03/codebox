package commands

import (
	"fmt"

	"gitlab.com/codebox4073715/codebox/cli/args"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
This function handles the command to approve a user
*/
func HandleApproveUser(args args.ApproveUserCmdArgs) uint {
	user, err := models.RetrieveUserByEmail(args.UserEmail)
	if err != nil {
		fmt.Println("Failed to retrieve user by email, unknown error")
		return 1
	}

	if user == nil {
		fmt.Println("No user found with the given email")
		return 1
	}

	if user.Approved {
		fmt.Println("Nothing to do here, user was already approved")
	} else {
		user.Approved = true
		if err := models.UpdateUser(user); err != nil {
			fmt.Println("Failed to update user, unknown error")
			return 1
		}

		fmt.Println("User has been approved")
	}

	return 0
}
