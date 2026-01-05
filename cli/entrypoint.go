package cli

import (
	"fmt"
	"os"

	"gitlab.com/codebox4073715/codebox/cli/args"
	"gitlab.com/codebox4073715/codebox/cli/commands"
)

/*
Switch-case that selects the command to run
*/
func RunCommand(a CLIArgs) uint {
	switch a.Command {
	case "runserver":
		return commands.HandleRunServer()
	case "set-password":
		return commands.HandleSetPassword()
	case "reset-ratelimit":
		return commands.HandleResetRatelimits()
	case "approve-user":
		return commands.HandleApproveUser(a.Args.(args.ApproveUserCmdArgs))
	default:
		fmt.Printf(
			"Invalid command '%s'\n", os.Args[1],
		)
	}

	return 1
}
