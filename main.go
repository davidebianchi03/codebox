package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"codebox.com/bgtasks"
	"codebox.com/env"

	"codebox.com/api"
	"codebox.com/db"
	"github.com/gin-gonic/gin"
)

func main() {
	err := env.InitCodeBoxEnv()
	if err != nil {
		log.Fatalf("Failed to load server configuration from environment: '%s'\n", err)
		return
	}

	// test della connessione con il database
	err = db.ConnectDB()
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
		// avvio dei background tasks
		err = bgtasks.InitBgTasks(env.CodeBoxEnv.RedisHost, env.CodeBoxEnv.RedisPort, uint(env.CodeBoxEnv.WorkspaceRelatedTasksConcurrency), "")
		if err != nil {
			log.Fatalln("cannot start background tasks")
			return
		}

		// avvio del server http
		// if env.CodeBoxEnv.DebugEnabled {
		gin.SetMode(gin.DebugMode)
		// } else {
		// 	gin.SetMode(gin.ReleaseMode)
		// }
		r := gin.Default()
		api.V1ApiRoutes(r)
		r.Run(fmt.Sprintf(":%s", strconv.Itoa(env.CodeBoxEnv.ServerPort)))
	default:
		log.Fatalf("Invalid command '%s'", os.Args[1])
		os.Exit(1)
	}
}
