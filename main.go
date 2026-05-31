package main

import (
	"fmt"
	"os"

	"gitlab.com/codebox4073715/codebox/cli"
)

// @title           Codebox API
// @version         {{version_placeholder}}
// @description     Codebox server

// @license.name  MIT
// @license.url   https://mit-license.org

// @host      localhost:8080
func main() {
	// parse cli args
	args, err := cli.ParseCLIArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// execute command and exit with status
	os.Exit(int(cli.RunCommand(args)))
}
