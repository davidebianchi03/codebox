package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/cli/args"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Try to verify an email address using cli
*/
func TestVerifyEmailFromCLI(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		user, err := models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Errorf("Failed to retrieve test user: '%s'\n", err)
			t.FailNow()
		}

		user.EmailVerified = false
		if err := models.UpdateUser(user); err != nil {
			t.Errorf("Failed to retrieve test user: '%s'\n", err)
			t.FailNow()
		}

		user, err = models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Errorf("Failed to retrieve updated test user: '%s'\n", err)
			t.FailNow()
		}

		assert.False(t, user.EmailVerified)

		cmdArgs := args.VerifyEmailCmdArgs{
			Email: user.Email,
		}

		exitCode := HandleVerifyEmail(cmdArgs)

		assert.Equal(t, uint(0), exitCode)

		user, err = models.RetrieveUserByEmail("user1@user.com")
		if err != nil || user == nil {
			t.Errorf("Failed to retrieve updated test user: '%s'\n", err)
			t.FailNow()
		}

		assert.True(t, user.EmailVerified)
	})
}

/*
Try to verify an email address using cli,
in this test is verified the case of missing email address
*/
func TestVerifyEmailFromCLIMissingEmail(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		cmdArgs := args.VerifyEmailCmdArgs{}

		exitCode := HandleVerifyEmail(cmdArgs)

		assert.Equal(t, uint(1), exitCode)
	})
}

/*
Try to verify an email address using cli,
in this test is verified the case of an invalid
*/
func TestVerifyEmailFromCLIInvalidEmailAddress(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		cmdArgs := args.VerifyEmailCmdArgs{
			Email: "invalid-email",
		}

		exitCode := HandleVerifyEmail(cmdArgs)

		assert.Equal(t, uint(1), exitCode)
	})
}

/*
Try to verify an email address using cli,
in this test is verified the case of a user that does not exist
*/
func TestVerifyEmailFromCLIInvalidUserDoesNotExist(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {

		cmdArgs := args.VerifyEmailCmdArgs{
			Email: "user@notexist.com",
		}

		exitCode := HandleVerifyEmail(cmdArgs)

		assert.Equal(t, uint(1), exitCode)
	})
}
