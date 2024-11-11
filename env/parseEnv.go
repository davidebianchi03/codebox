package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type codeBoxEnv struct {
	DebugEnabled                     bool
	ServerPort                       int
	WorkspaceRelatedTasksConcurrency int
	DbDriver                         string
	DbURL                            string
	RedisHost                        string
	RedisPort                        int
	UploadsPath                      string
	UseGravatar                      bool
	FrontendPath                     string
}

var CodeBoxEnv *codeBoxEnv

func envVarOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func InitCodeBoxEnv() error {
	codeboxEnvFilename := "codebox.env"
	err := godotenv.Load(codeboxEnvFilename)
	if err != nil {
		return fmt.Errorf("failed to load environement variables from %s", codeboxEnvFilename)
	}

	CodeBoxEnv = new(codeBoxEnv)

	CodeBoxEnv.ServerPort, err = strconv.Atoi(envVarOrDefault("CODEBOX_SERVER_PORT", "8080"))
	if err != nil {
		return fmt.Errorf("invalid value for CODEBOX_SERVER_PORT env var")
	}

	// background tasks
	CodeBoxEnv.WorkspaceRelatedTasksConcurrency, err = strconv.Atoi(envVarOrDefault("CODEBOX_WORKSPACE_TASKS_CONCURRENCY", "4"))
	if err != nil {
		return fmt.Errorf("invalid value for CODEBOX_WORKSPACE_TASKS_CONCURRENCY env var")
	}

	// database
	CodeBoxEnv.DbDriver = envVarOrDefault("CODEBOX_DB_DRIVER", "sqlite3")
	CodeBoxEnv.DbURL = envVarOrDefault("CODEBOX_DB_URL", "")
	if CodeBoxEnv.DbURL == "" {
		return fmt.Errorf("CODEBOX_DB_URL not defined")
	}

	// redis
	CodeBoxEnv.RedisHost = envVarOrDefault("CODEBOX_REDIS_HOST", "127.0.0.1")
	CodeBoxEnv.RedisPort, err = strconv.Atoi(envVarOrDefault("CODEBOX_REDIS_PORT", "6379"))
	if err != nil {
		return fmt.Errorf("invalid value for CODEBOX_REDIS_PORT env var")
	}

	// uploads
	CodeBoxEnv.UploadsPath = envVarOrDefault("CODEBOX_DATA_PATH", "./data")

	// use gravatar
	CodeBoxEnv.UseGravatar = strings.ToLower(envVarOrDefault("CODEBOX_USE_GRAVATAR", "true")) == "true"

	// debug
	CodeBoxEnv.DebugEnabled = strings.ToLower(envVarOrDefault("CODEBOX_DEBUG", "true")) == "true"

	return nil
}
