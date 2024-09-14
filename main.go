package main

import (
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strconv"
	"strings"

	_ "database/sql"

	"codebox.com/bgtasks"
	"codebox.com/env"
	_ "codebox.com/migrations"

	"codebox.com/api"
	"codebox.com/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh/terminal"
)

func a() {
	fmt.Println("Hello world")
}

func main() {
	err := env.InitCodeBoxEnv()
	if err != nil {
		log.Fatalf("Failed to load server configuration from environment: '%s'\n", err)
		return
	}

	// test della connessione con il database
	err = db.InitDBConnection(env.CodeBoxEnv.DbDriver, env.CodeBoxEnv.DbURL)
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
		// avvio dei background tasks
		err = bgtasks.InitBgTasks(env.CodeBoxEnv.RedisHost, env.CodeBoxEnv.RedisPort, uint(env.CodeBoxEnv.WorkspaceRelatedTasksConcurrency), "")
		if err != nil {
			log.Fatalln("cannot start background tasks")
			return
		}

		// avvio del server http
		r := gin.Default()
		api.V1ApiRoutes(r)
		r.Run(fmt.Sprintf(":%s", strconv.Itoa(env.CodeBoxEnv.ServerPort)))
	case "create-superuser":
		// creazione superuser
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		// controllo se la l'indirizzo email è valido
		_, err := mail.ParseAddress(email)
		if err != nil {
			fmt.Println("Invalid email address")
			os.Exit(1)
		}

		// check if email already exists
		var foundUser db.User
		result := db.DB.Where("email=?", email).Find(&foundUser)
		if result.Error != nil {
			fmt.Println("An error occured creating new superuser")
			os.Exit(1)
		}

		if foundUser.Id != 0 {
			fmt.Println("User with this email already exists")
			os.Exit(1)
		}

		password := "password"
		confirm_password := "confirm_password"
		mismatching_passwords_count := 0

		for password != confirm_password && mismatching_passwords_count < 3 {
			fmt.Print("Enter password: ")
			passwordB, _ := terminal.ReadPassword(0)
			fmt.Println("")
			fmt.Print("Confirm password: ")
			confirmPasswordB, _ := terminal.ReadPassword(0)

			password = string(passwordB)
			confirm_password = string(confirmPasswordB)

			password = strings.TrimSpace(password)
			confirm_password = strings.TrimSpace(confirm_password)

			if password != confirm_password {
				fmt.Println("Mismatching passwords")
				mismatching_passwords_count += 1
			}
		}

		if password != confirm_password {
			fmt.Println("Too many attemps")
			os.Exit(1)
		}

		user := db.User{Email: email, FirstName: "", LastName: "", Password: password}
		result = db.DB.Create(&user)

		if result.Error != nil {
			fmt.Println("An error occured creating new superuser")
			os.Exit(1)
		} else {
			fmt.Println("Superuser has been successfully created")
			os.Exit(0)
		}
	default:
		log.Fatalf("Invalid command '%s'", os.Args[1])
		os.Exit(1)
	}
}
