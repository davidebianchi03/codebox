package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/cli/args"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Try to approve a user using cli
*/
func TestApproveUserFromCLI(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Errorf("Failed to retrieve test user: '%s'\n", err)
			t.FailNow()
		}

		user.Approved = false
		if err := models.UpdateUser(user); err != nil {
			t.Errorf("Failed to retrieve test user: '%s'\n", err)
			t.FailNow()
		}

		cmdArgs := args.ApproveUserCmdArgs{
			UserEmail: user.Email,
		}

		exitCode := HandleApproveUser(cmdArgs)

		assert.Equal(t, uint(0), exitCode)

		user, err = models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Errorf("Failed to retrieve updated test user: '%s'\n", err)
			t.FailNow()
		}

		assert.True(t, user.Approved)
	})
}

/*
Try to approve a user using cli,
in this test is verified the case of missing email address
*/
func TestApproveUserFromCLIMissingEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		cmdArgs := args.ApproveUserCmdArgs{}

		exitCode := HandleApproveUser(cmdArgs)

		assert.Equal(t, uint(1), exitCode)
	})
}

/*
Try to approve a user using cli,
in this test is verified the case of an invalid
*/
func TestApproveUserFromCLIInvalidEmailAddress(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		cmdArgs := args.ApproveUserCmdArgs{
			UserEmail: "invalid-email",
		}

		exitCode := HandleApproveUser(cmdArgs)

		assert.Equal(t, uint(1), exitCode)
	})
}

/*
Try to approve a user using cli,
in this test is verified the case of a user that does not exist
*/
func TestApproveUserFromCLIInvalidUserDoesNotExist(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {

		cmdArgs := args.ApproveUserCmdArgs{
			UserEmail: "user@notexist.com",
		}

		exitCode := HandleApproveUser(cmdArgs)

		assert.Equal(t, uint(1), exitCode)
	})
}
