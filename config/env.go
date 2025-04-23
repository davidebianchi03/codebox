package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var ServerVersion = "dbg-9.9.9"

type EnvVars struct {
	ExternalUrl      string `env:"CODEBOX_EXTERNAL_URL,required"`
	DebugEnabled     bool   `env:"CODEBOX_DEBUG" envDefault:"true"`
	DBDriver         string `env:"CODEBOX_DB_DRIVER" envDefault:"postgres"`
	DBHost           string `env:"CODEBOX_DB_HOST" envDefault:"postgres"`
	DBPort           int    `env:"CODEBOX_DB_PORT" envDefault:"5432"`
	DBName           string `env:"CODEBOX_DB_NAME" envDefault:"codebox"`
	DBUser           string `env:"CODEBOX_DB_USER" envDefault:"codebox"`
	DBPassword       string `env:"CODEBOX_DB_PASSWORD" envDefault:"password"`
	ServerPort       int    `env:"CODEBOX_SERVER_PORT" envDefault:"8080"`
	TasksConcurrency int    `env:"CODEBOX_WORKSPACE_CONCURRENCY" envDefault:"5"`
	RedisHost        string `env:"CODEBOX_REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort        int    `env:"CODEBOX_REDIS_PORT" envDefault:"6379"`
	UploadsPath      string `env:"CODEBOX_DATA_PATH" envDefault:"./data"`
	UseGravatar      bool   `env:"CODEBOX_USE_GRAVATAR" envDefault:"true"`
	UseSubDomains    bool   `env:"CODEBOX_USE_SUBDOMAINS" envDefault:"true"`
	AuthCookieName   string `env:"CODEBOX_AUTH_COOKIE_NAME" envDefault:"jwtToken"`
}

var Environment *EnvVars

func InitCodeBoxEnv() error {
	codeboxEnvFilename := "codebox.env"
	err := godotenv.Load(codeboxEnvFilename)
	if err != nil {
		return fmt.Errorf("failed to load environement variables from %s", codeboxEnvFilename)
	}

	e := EnvVars{}
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
