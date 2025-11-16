package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"
	"syscall"

	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/router"
	"golang.org/x/term"
)

func prepareTerminal() *term.Terminal {
	if !term.IsTerminal(syscall.Stdin) {
		slog.Warn("std::cin is not a terminal")
	}
	if !term.IsTerminal(syscall.Stdout) {
		slog.Warn("std::cout is not a terminal")
	}

	oldState, err := term.MakeRaw(syscall.Stdin)
	defer func() {
		err = errors.Join(err, term.Restore(syscall.Stdin, oldState))
	}()

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	terminal := term.NewTerminal(screen, "")
	return terminal
}

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
	case "runserver":
		err = bgtasks.InitBgTasks(
			config.Environment.RedisHost,
			config.Environment.RedisPort,
			uint(config.Environment.TasksConcurrency),
			"",
		)
		if err != nil {
			log.Fatalln("cannot start background tasks")
			return
		}

		r := router.SetupRouter()
		log.Printf("listening at 0.0.0.0:%d\n", config.Environment.ServerPort)
		r.Run(fmt.Sprintf(":%s", strconv.Itoa(config.Environment.ServerPort)))
		os.Exit(0)
	case "set-password":
		terminal := prepareTerminal()
		fmt.Print("Enter email: ")
		email, err := terminal.ReadLine()
		if err != nil {
			os.Exit(1)
			return
		}

		password, err := terminal.ReadPassword("New password:")
		if err != nil {
			os.Exit(1)
			return
		}

		passwordConfirm, err := terminal.ReadPassword("Confirm the password")
		if err != nil {
			os.Exit(1)
			return
		}

		if password != passwordConfirm {
			fmt.Println("Passwords do not match")
			os.Exit(1)
			return
		}

		user, err := models.RetrieveUserByEmail(email)
		if err != nil {
			fmt.Println("Failed to retrieve user")
			os.Exit(1)
			return
		}

		if user == nil {
			fmt.Println("User not found")
			os.Exit(1)
			return
		}

		if err := models.ValidatePassword(password); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		password, err = models.HashPassword(password)
		if err != nil {
			fmt.Println("Unknown error")
			os.Exit(1)
			return
		}
		user.Password = password
		dbconn.DB.Save(&user)

		os.Exit(0)
	default:
		log.Fatalf("Invalid command '%s'", os.Args[1])
		os.Exit(1)
	}
}
