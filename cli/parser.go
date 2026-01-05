package cli

import (
	"errors"
	"flag"
	"fmt"
	"net/mail"
	"os"

	"gitlab.com/codebox4073715/codebox/cli/args"
)

type CLIArgs struct {
	Command string
	Args    any
}

/*
Parse cli args, return an object with command and parsed args.
Return an error if some argument is missing or if command does
not exist
*/
func ParseCLIArgs() (CLIArgs, error) {
	if len(os.Args) < 2 {
		return CLIArgs{}, errors.New("A command is expected")
	}

	switch os.Args[1] {
	case "runserver":
		return CLIArgs{
			Command: "runserver",
			Args:    nil,
		}, nil
	case "set-password":
		return CLIArgs{
			Command: "set-password",
			Args:    nil,
		}, nil
	case "reset-ratelimit":
		return CLIArgs{
			Command: "reset-ratelimit",
			Args:    nil,
		}, nil
	case "approve-user":
		var approveUserArgs args.ApproveUserCmdArgs
		approveUserCmd := flag.NewFlagSet("approve-user", flag.ExitOnError)
		approveUserCmd.StringVar(&approveUserArgs.UserEmail, "user-email", "", "email address of the user to approve")
		approveUserCmd.Parse(os.Args[2:])

		if approveUserArgs.UserEmail == "" {
			return CLIArgs{}, errors.New("arg 'user-email' is required")
		}

		if _, err := mail.ParseAddress(approveUserArgs.UserEmail); err != nil {
			return CLIArgs{}, errors.New("provided value for 'user-email' is not a valid email address")
		}

		return CLIArgs{
			Command: "approve-user",
			Args:    approveUserArgs,
		}, nil
	case "verify-email":
		var verifyEmailArgs args.VerifyEmailCmdArgs
		verifyEmailCmd := flag.NewFlagSet("verify-email", flag.ExitOnError)
		verifyEmailCmd.StringVar(&verifyEmailArgs.Email, "email-address", "", "email address to verify")
		verifyEmailCmd.Parse(os.Args[2:])

		if verifyEmailArgs.Email == "" {
			return CLIArgs{}, errors.New("arg 'email-address' is required")
		}

		if _, err := mail.ParseAddress(verifyEmailArgs.Email); err != nil {
			return CLIArgs{}, errors.New("provided value for 'email-address' is not a valid email address")
		}

		return CLIArgs{
			Command: "verify-email",
			Args:    verifyEmailArgs,
		}, nil
	default:
		return CLIArgs{}, fmt.Errorf("Invalid command '%s'", os.Args[1])
	}
}
