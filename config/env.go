package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"regexp"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var ServerVersion = "dbg-9.9.9"

type EnvVars struct {
	// domains
	ExternalUrl    string `env:"CODEBOX_EXTERNAL_URL,required"`
	UseSubDomains  bool   `env:"CODEBOX_USE_SUBDOMAINS" envDefault:"true"`
	WildcardDomain string `env:"CODEBOX_WILDCARD_DOMAIN"`
	// server
	ServerPort   int  `env:"CODEBOX_SERVER_PORT" envDefault:"8080"`
	DebugEnabled bool `env:"CODEBOX_DEBUG" envDefault:"true"`
	// cookies
	AuthCookieName          string `env:"CODEBOX_AUTH_COOKIE_NAME" envDefault:"codebox_auth_token"`
	SubdomainAuthCookieName string `env:"CODEBOX_SUBDOMAIN_AUTH_COOKIE_NAME" envDefault:"subdomain_codebox_auth_token"`
	// paths
	UploadsPath     string `env:"CODEBOX_DATA_PATH" envDefault:"./data"`
	CliBinariesPath string `env:"CODEBOX_CLI_BINARIES_PATH" envDefault:"./cli"`
	TemplatesFolder string `env:"CODEBOX_TEMPLATES_FOLDER" envDefault:"./templates"`
	// runner
	RunnerTokenHeader     string `env:"CODEBOX_RUNNER_TOKEN_HEADER" envDefault:"X-Codebox-Runner-Token"`
	RunnerTokenQueryParam string `env:"CODEBOX_RUNNER_TOKEN_QUERY_PARAM" envDefault:"runner_token"`
	// database
	DBDriver   string `env:"CODEBOX_DB_DRIVER" envDefault:"mysql"`
	DBHost     string `env:"CODEBOX_DB_HOST" envDefault:"db"`
	DBPort     int    `env:"CODEBOX_DB_PORT" envDefault:"3306"`
	DBName     string `env:"CODEBOX_DB_NAME" envDefault:"codebox"`
	DBTestName string `env:"CODEBOX_TEST_DB_NAME" envDefault:"codebox-test"`
	DBUser     string `env:"CODEBOX_DB_USER" envDefault:"codebox"`
	DBPassword string `env:"CODEBOX_DB_PASSWORD" envDefault:"password"`
	// bg tasks
	TasksConcurrency int `env:"CODEBOX_BG_TASKS_CONCURRENCY" envDefault:"5"`
	// redis
	RedisHost string `env:"CODEBOX_REDIS_HOST" envDefault:"redis"`
	RedisPort int    `env:"CODEBOX_REDIS_PORT" envDefault:"6379"`
	// email related settings
	EmailSMTPHost     string `env:"CODEBOX_EMAIL_SMTP_HOST"`
	EmailSMTPPort     int    `env:"CODEBOX_EMAIL_SMTP_PORT"`
	EmailSMTPUser     string `env:"CODEBOX_EMAIL_SMTP_USER"`
	EmailSMTPPassword string `env:"CODEBOX_EMAIL_SMTP_PASSWORD"`
}

var Environment *EnvVars
var SkipPathValidation = false

func (e *EnvVars) ValidateExternalUrl() error {
	if e.ExternalUrl != "" {
		parsedURL, err := url.Parse(e.ExternalUrl)
		if err != nil {
			return fmt.Errorf("CODEBOX_EXTERNAL_URL is not a valid URL: %w", err)
		}
		if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
			return fmt.Errorf("CODEBOX_EXTERNAL_URL must use HTTPS or HTTP scheme, got: %s", parsedURL.Scheme)
		}
	}
	return nil
}

func (e *EnvVars) ValidateWildcardDomain() error {
	if e.UseSubDomains {
		if e.WildcardDomain == "" {
			return fmt.Errorf("CODEBOX_WILDCARD_DOMAIN must be set when CODEBOX_USE_SUBDOMAINS is true")
		}
	}

	if e.WildcardDomain != "" && !isValidDomain(e.WildcardDomain) {
		return fmt.Errorf("CODEBOX_WILDCARD_DOMAIN is not a valid domain: %s", e.WildcardDomain)
	}

	return nil
}

func (e *EnvVars) ValidateServerPort() error {
	if e.ServerPort < 1 || e.ServerPort > 65535 {
		return errors.New("CODEBOX_SERVER_PORT is not valid")
	}
	return nil
}

func (e *EnvVars) ValidateAuthCookieName() error {
	if e.AuthCookieName == "" {
		return errors.New("CODEBOX_AUTH_COOKIE_NAME cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateSubdomainAuthCookieName() error {
	if e.SubdomainAuthCookieName == "" {
		return errors.New("CODEBOX_SUBDOMAIN_AUTH_COOKIE_NAME cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateUploadsPath() error {
	if e.UploadsPath == "" {
		return errors.New("CODEBOX_DATA_PATH cannot be empty")
	}

	return nil
}

func (e *EnvVars) ValidateCliBinariesPath() error {
	if e.CliBinariesPath == "" {
		return errors.New("CODEBOX_CLI_BINARIES_PATH cannot be empty")
	}

	if SkipPathValidation {
		return nil
	}

	info, err := os.Stat(e.CliBinariesPath)
	if err != nil {
		return fmt.Errorf("invalid value for CODEBOX_CLI_BINARIES_PATH %s", err.Error())
	}

	if !info.IsDir() {
		return errors.New("CODEBOX_CLI_BINARIES_PATH is not a directory")
	}

	return nil
}

func (e *EnvVars) ValidateTemplatesFolder() error {
	if e.TemplatesFolder == "" {
		return errors.New("CODEBOX_TEMPLATES_FOLDER cannot be empty")
	}

	if SkipPathValidation {
		return nil
	}

	info, err := os.Stat(e.TemplatesFolder)
	if err != nil {
		return fmt.Errorf("invalid value for CODEBOX_TEMPLATES_FOLDER %s", err.Error())
	}

	if !info.IsDir() {
		return errors.New("CODEBOX_TEMPLATES_FOLDER is not a directory")
	}

	return nil
}

func (e *EnvVars) ValidateRunnerTokenHeader() error {
	if e.RunnerTokenHeader == "" {
		return errors.New("CODEBOX_RUNNER_TOKEN_HEADER cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateRunnerTokenQueryParam() error {
	if e.RunnerTokenQueryParam == "" {
		return errors.New("CODEBOX_RUNNER_TOKEN_QUERY_PARAM cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateDBDriver() error {
	if e.DBDriver == "" {
		return errors.New("CODEBOX_DB_DRIVER cannot be empty")
	}

	if e.DBDriver != "sqlite3" && e.DBDriver != "mysql" {
		return errors.New("CODEBOX_DB_DRIVER unsupported db driver")
	}

	return nil
}

func (e *EnvVars) ValidateDBHost() error {
	if e.DBHost == "" {
		return errors.New("CODEBOX_DB_HOST cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateDBPort() error {
	if e.DBPort < 1 || e.DBPort > 65535 {
		return errors.New("CODEBOX_DB_PORT is not valid")
	}
	return nil
}

func (e *EnvVars) ValidateDBName() error {
	if e.DBName == "" {
		return errors.New("CODEBOX_DB_NAME cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateDBTestName() error {
	if e.DBTestName == "" {
		return errors.New("CODEBOX_TEST_DB_NAME cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateDBUser() error {
	if e.DBUser == "" {
		return errors.New("CODEBOX_DB_USER cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateDBPassword() error {
	if e.DBPassword == "" {
		return errors.New("CODEBOX_DB_PASSWORD cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateTasksConcurrency() error {
	if e.TasksConcurrency < 1 {
		return errors.New("CODEBOX_BG_TASKS_CONCURRENCY cannot be less than 1")
	}
	return nil
}

func (e *EnvVars) ValidateRedisHost() error {
	if e.RedisHost == "" {
		return errors.New("CODEBOX_REDIS_HOST cannot be empty")
	}
	return nil
}

func (e *EnvVars) ValidateRedisPort() error {
	if e.RedisPort < 1 || e.RedisPort > 65535 {
		return errors.New("CODEBOX_REDIS_PORT is not valid")
	}
	return nil
}
func (e *EnvVars) ValidateEmailSMTPPort() error {
	if e.EmailSMTPPort != -1 {
		if e.EmailSMTPPort < 1 || e.EmailSMTPPort > 65535 {
			return errors.New("CODEBOX_EMAIL_SMTP_PORT is not valid")
		}
	}
	return nil
}

/*
IsValid automatically calls all ValidateXXX methods for fields in the struct
using reflection to discover and invoke field-specific validators
*/
func (e *EnvVars) IsValid() error {
	eVal := reflect.ValueOf(e).Elem()
	eType := eVal.Type()

	// Iterate over all fields in the struct
	for i := 0; i < eType.NumField(); i++ {
		fieldName := eType.Field(i).Name

		// Look for the Validate<FieldName> method
		methodName := "Validate" + fieldName
		method := reflect.ValueOf(e).MethodByName(methodName)

		// If the method exists, call it
		if method.IsValid() {
			results := method.Call([]reflect.Value{})

			// Check if there's an error returned
			if len(results) > 0 && !results[0].IsNil() {
				if err, ok := results[0].Interface().(error); ok {
					return err
				}
			}
		}
	}

	return nil
}

/*
Check if a string is a valid domain name
Domain rules:
- Must contain at least one dot (e.g., example.com)
- Each label must start and end with alphanumeric character
- Labels can contain hyphens
- Labels must be 1-63 characters
- Total length must not exceed 253 characters
*/
func isValidDomain(domain string) bool {
	if len(domain) > 253 {
		return false
	}

	// Domain must contain at least one dot
	if !regexp.MustCompile(`\.`).MatchString(domain) {
		return false
	}

	// Each label (part between dots) must be valid
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
	return domainRegex.MatchString(domain)
}

/*
Load configuration from environment variables or from a .env file
*/
func InitCodeBoxEnv() error {
	codeboxEnvFilename := os.Getenv("CODEBOX_ENV_FILE")
	if codeboxEnvFilename == "" {
		codeboxEnvFilename = "codebox.env"
	}
	err := godotenv.Load(codeboxEnvFilename)
	if err != nil {
		return fmt.Errorf("failed to load environement variables from %s", codeboxEnvFilename)
	}

	e := EnvVars{}
	err = env.Parse(&e)
	if err != nil {
		return err
	}

	// Validate configuration
	if err = e.IsValid(); err != nil {
		return err
	}

	if err = os.MkdirAll(e.UploadsPath, 0777); err != nil {
		return err
	}

	Environment = &e
	return nil
}

/*
Get if the email service is configured
*/
func IsEmailConfigured() bool {
	return Environment.EmailSMTPHost != "" && Environment.EmailSMTPPort != 0 &&
		Environment.EmailSMTPUser != "" && Environment.EmailSMTPPassword != ""
}
