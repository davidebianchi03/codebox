package main

import (
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	_ "database/sql"

	_ "codebox.com/migrations"

	"codebox.com/api"
	"codebox.com/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	godotenv.Load("codebox.env")

	// test della connessione con il database
	err := db.InitDBConnection()
	if err != nil {
		log.Fatalf("Cannot init connection with DB '%s'\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		log.Fatalln("A command is expected")
		os.Exit(1)
	}

	// apply migrations
	if err := goose.SetDialect(os.Getenv("CODEBOX_DB_DRIVER")); err != nil {
		panic(err)
	}
	if err := goose.Up(db.SqlDB, "migrations"); err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "runserver":
		// avvio del server http
		// gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		api.V1ApiRoutes(r)
		r.Run(":8080")
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

	// var users []db.User
	// result := dbConn.Find(&users)
	// if result.Error != nil {
	// 	// handle error
	// }

	// if len(users) == 0 {
	// 	// create first user
	// }

	// fmt.Println(result.RowsAffected)
	// fmt.Println(result.Error)
	// fmt.Println(users)

	// stmts, err := gormschema.New("sqlite").Load(&db.User{})
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
	// 	os.Exit(1)
	// }
	// io.WriteString(os.Stdout, stmts)

	// db_name := "codebox.sqlite3"

	// _, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	// if err != nil {
	// 	fmt.Errorf("failed to connect database %s", err)
	// }

	// r := gin.Default()
	// r.GET("/ping", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

}
