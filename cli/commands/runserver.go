package commands

import (
	"fmt"
	"log"
	"strconv"

	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/httpserver"
)

/*
Handle command to start codebox http server
Returns an error if something fails
*/
func HandleRunServer() uint {
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

	// init bg tasks
	if err := bgtasks.InitBgTasks(
		uint(config.Environment.TasksConcurrency),
		"",
	); err != nil {
		log.Println("Cannot start background tasks")
		return 1
	}

	// run server
	r := httpserver.SetupRouter()
	listeningAddress := fmt.Sprintf(":%s", strconv.Itoa(config.Environment.ServerPort))
	log.Printf(
		"listening at %s\n",
		listeningAddress,
	)

	r.Run(
		listeningAddress,
	)

	return 0
}
