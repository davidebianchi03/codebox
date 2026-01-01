package commands

import (
	"fmt"
	"log"
	"strconv"

	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/httpserver"
)

/*
Handle command to start codebox http server
Returns an error if something fails
*/
func HandleRunServer() uint {
	// init bg tasks
	err := bgtasks.InitBgTasks(
		uint(config.Environment.TasksConcurrency),
		"",
	)
	if err != nil {
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
