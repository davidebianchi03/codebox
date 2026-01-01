package main

import (
	"log"
	"os"

	"gitlab.com/codebox4073715/codebox/cli"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
)

// @title           Codebox API
// @version         {{version_placeholder}}
// @description     Codebox server

// @license.name  MIT
// @license.url   https://mit-license.org

// @host      localhost:8080
func main() {
	err := config.InitCodeBoxEnv()
	if err != nil {
		log.Fatalf("Failed to load server configuration from environment: '%s'\n", err)
		return
	}

	// test db connection
	err = dbconn.ConnectDB()
	if err != nil {
		log.Fatalf("Cannot init connection with DB: '%s'\n", err)
		return
	}

	// TODO: test redis connection

	// parse cli args
	args, err := cli.ParseCLIArgs()
	if err != nil {
		log.Fatal(err)
		return
	}

	// execute command and exit with status
	os.Exit(int(cli.RunCommand(args)))
}
