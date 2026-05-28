package config

import (
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var ServerVersion = "dbg-9.9.9"

type EnvVars struct {
	ExternalUrl             string `env:"CODEBOX_EXTERNAL_URL,required"`
	DebugEnabled            bool   `env:"CODEBOX_DEBUG" envDefault:"true"`
	DBDriver                string `env:"CODEBOX_DB_DRIVER" envDefault:"mysql"`
	DBHost                  string `env:"CODEBOX_DB_HOST" envDefault:"db"`
	DBPort                  int    `env:"CODEBOX_DB_PORT" envDefault:"3306"`
	DBName                  string `env:"CODEBOX_DB_NAME" envDefault:"codebox"`
	DBTestName              string `env:"CODEBOX_TEST_DB_NAME" envDefault:"codebox-test"`
	DBUser                  string `env:"CODEBOX_DB_USER" envDefault:"codebox"`
	DBPassword              string `env:"CODEBOX_DB_PASSWORD" envDefault:"password"`
	ServerPort              int    `env:"CODEBOX_SERVER_PORT" envDefault:"8080"`
	TasksConcurrency        int    `env:"CODEBOX_BG_TASKS_CONCURRENCY" envDefault:"5"`
	RedisHost               string `env:"CODEBOX_REDIS_HOST" envDefault:"redis"`
	RedisPort               int    `env:"CODEBOX_REDIS_PORT" envDefault:"6379"`
	UploadsPath             string `env:"CODEBOX_DATA_PATH" envDefault:"./data"`
	UseSubDomains           bool   `env:"CODEBOX_USE_SUBDOMAINS" envDefault:"true"`
	WildcardDomain          string `env:"CODEBOX_WILDCARD_DOMAIN"`
	AuthCookieName          string `env:"CODEBOX_AUTH_COOKIE_NAME" envDefault:"codebox_auth_token"`
	SubdomainAuthCookieName string `env:"CODEBOX_SUBDOMAIN_AUTH_COOKIE_NAME" envDefault:"subdomain_codebox_auth_token"`
	CliBinariesPath         string `env:"CODEBOX_CLI_BINARIES_PATH" envDefault:"./cli"`
	TemplatesFolder         string `env:"CODEBOX_TEMPLATES_FOLDER" envDefault:"./templates"`
	RunnerTokenHeader       string `env:"CODEBOX_RUNNER_TOKEN_HEADER" envDefault:"X-Codebox-Runner-Token"`
	RunnerTokenQueryParam   string `env:"CODEBOX_RUNNER_TOKEN_QUERY_PARAM" envDefault:"runner_token"`
	// email related settings
	EmailSMTPHost     string `env:"CODEBOX_EMAIL_SMTP_HOST"`
	EmailSMTPPort     int    `env:"CODEBOX_EMAIL_SMTP_PORT"`
	EmailSMTPUser     string `env:"CODEBOX_EMAIL_SMTP_USER"`
	EmailSMTPPassword string `env:"CODEBOX_EMAIL_SMTP_PASSWORD"`
}

var Environment *EnvVars

/*
Validate environment variables after parsing
*/
func validateEnvVars(e *EnvVars) error {
	// Validate ExternalUrl is a valid HTTPS URL
	if e.ExternalUrl != "" {
		parsedURL, err := url.Parse(e.ExternalUrl)
		if err != nil {
			return fmt.Errorf("CODEBOX_EXTERNAL_URL is not a valid URL: %w", err)
		}
		if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
			return fmt.Errorf("CODEBOX_EXTERNAL_URL must use HTTPS or HTTP scheme, got: %s", parsedURL.Scheme)
		}
	}

	// Validate WildcardDomain when UseSubDomains is enabled
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
	if err = validateEnvVars(&e); err != nil {
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
