package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type envVars struct {
	DebugEnabled           bool   `env:"CODEBOX_DEBUG" envDefault:"true"`
	DBUrl                  string `env:"CODEBOX_DB_URL" envDefault:"sqlite://./codebox.db"`
	ServerPort             int    `env:"CODEBOX_SERVER_PORT" envDefault:"8080"`
	TasksConcurrency       int    `env:"CODEBOX_WORKSPACE_CONCURRENCY" envDefault:"1"`
	RedisHost              string `env:"CODEBOX_REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort              int    `env:"CODEBOX_REDIS_PORT" envDefault:"6379"`
	UploadsPath            string `env:"CODEBOX_DATA_PATH" envDefault:"./data"`
	AllowSignUp            bool   `env:"CODEBOX_ALLOW_SIGNUP" envDefault:"false"`
	UseGravatar            bool   `env:"CODEBOX_USE_GRAVATAR" envDefault:"true"`
	UseSubDomains          bool   `env:"CODEBOX_USE_SUBDOMAINS" envDefault:"true"`
	WorkspaceObjectsPrefix string `env:"CODEBOX_WORKSPACE_OBJECTS_PREFIX" envDefault:"codebox"` // TODO: remove
	DevcontainerCmd        string `env:"CODEBOX_DEVCONTAINERS_COMMAND" envDefault:"devcontainer"`
}

var Environment *envVars

func InitCodeBoxEnv() error {
	codeboxEnvFilename := "codebox.env"
	err := godotenv.Load(codeboxEnvFilename)
	if err != nil {
		return fmt.Errorf("failed to load environement variables from %s", codeboxEnvFilename)
	}

	e := envVars{}
	err = env.Parse(&e)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(e.UploadsPath, 0777); err != nil {
		return err
	}

	Environment = &e
	return nil
}
