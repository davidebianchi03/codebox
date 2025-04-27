package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"

	"gitlab.com/codebox4073715/codebox/api"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
)

func main() {
	err := config.InitCodeBoxEnv()
	if err != nil {
		log.Fatalf("Failed to load server configuration from environment: '%s'\n", err)
		return
	}

	// test della connessione con il database
	err = dbconn.ConnectDB()
	if err != nil {
		log.Fatalf("Cannot init connection with DB: '%s'\n", err)
		return
	}

	if len(os.Args) < 2 {
		log.Fatalln("A command is expected")
		return
	}

	switch os.Args[1] {
	// TODO: command to reset the password for a user
	case "runserver":
		err = bgtasks.InitBgTasks(config.Environment.RedisHost, config.Environment.RedisPort, uint(config.Environment.TasksConcurrency), "")
		if err != nil {
			log.Fatalln("cannot start background tasks")
			return
		}

		r := api.SetupRouter()
		log.Printf("listening at 0.0.0.0:%d\n", config.Environment.ServerPort)
		r.Run(fmt.Sprintf(":%s", strconv.Itoa(config.Environment.ServerPort)))
	default:
		log.Fatalf("Invalid command '%s'", os.Args[1])
		os.Exit(1)
	}
}
