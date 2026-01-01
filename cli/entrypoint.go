package cli

import (
	"fmt"
	"os"

	"gitlab.com/codebox4073715/codebox/cli/commands"
)

/*
Switch-case that selects the command to run
*/
func RunCommand(args CLIArgs) uint {
	switch os.Args[1] {
	case "runserver":
		return commands.HandleRunServer()
	case "set-password":
		return commands.HandleSetPassword()
	case "reset-ratelimit":
		return commands.HandleResetRatelimits()
	default:
		fmt.Printf(
			"Invalid command '%s'\n", os.Args[1],
		)
	}

	return 1
}
