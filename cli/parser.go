package cli

import (
	"errors"
	"fmt"
	"os"
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
	default:
		return CLIArgs{}, fmt.Errorf("Invalid command '%s'", os.Args[1])
	}
}
