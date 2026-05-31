package commands

import (
	"fmt"

	"gitlab.com/codebox4073715/codebox/config"
)

/*
This function handles the command to check the environment variables
*/
func HandleCheckEnv() uint {
	if err := config.InitCodeBoxEnv(); err != nil {
		fmt.Printf("Failed to load server configuration from environment: '%s'\n", err)
		return 1
	}

	fmt.Println("Config is valid")
	return 0
}
