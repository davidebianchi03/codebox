package config

import (
	"os"
	"testing"
)

func TestValidateExternalUrl(t *testing.T) {
	tests := []struct {
		name        string
		externalUrl string
		expectError bool
	}{
		{
			name:        "valid https url",
			externalUrl: "https://example.com",
			expectError: false,
		},
		{
			name:        "valid http url",
			externalUrl: "http://localhost:8080",
			expectError: false,
		},
		{
			name:        "empty url is valid",
			externalUrl: "",
			expectError: false,
		},
		{
			name:        "invalid scheme ftp",
			externalUrl: "ftp://example.com",
			expectError: true,
		},
		{
			name:        "invalid url format",
			externalUrl: "not a valid url",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				ExternalUrl: tt.externalUrl,
			}
			err := e.ValidateExternalUrl()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateExternalUrl() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateWildcardDomain(t *testing.T) {
	tests := []struct {
		name           string
		useSubDomains  bool
		wildcardDomain string
		expectError    bool
	}{
		{
			name:           "valid domain with subdomains enabled",
			useSubDomains:  true,
			wildcardDomain: "example.com",
			expectError:    false,
		},
		{
			name:           "empty domain with subdomains disabled",
			useSubDomains:  false,
			wildcardDomain: "",
			expectError:    false,
		},
		{
			name:           "missing domain when subdomains enabled",
			useSubDomains:  true,
			wildcardDomain: "",
			expectError:    true,
		},
		{
			name:           "invalid domain format",
			useSubDomains:  true,
			wildcardDomain: "example",
			expectError:    true,
		},
		{
			name:           "invalid domain - no dot",
			useSubDomains:  false,
			wildcardDomain: "localhost",
			expectError:    true,
		},
		{
			name:           "valid complex domain",
			useSubDomains:  true,
			wildcardDomain: "sub-domain.example.co.uk",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				UseSubDomains:  tt.useSubDomains,
				WildcardDomain: tt.wildcardDomain,
			}
			err := e.ValidateWildcardDomain()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateWildcardDomain() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateServerPort(t *testing.T) {
	tests := []struct {
		name        string
		port        int
		expectError bool
	}{
		{
			name:        "valid port 8080",
			port:        8080,
			expectError: false,
		},
		{
			name:        "valid port 1",
			port:        1,
			expectError: false,
		},
		{
			name:        "valid port 65535",
			port:        65535,
			expectError: false,
		},
		{
			name:        "invalid port 0",
			port:        0,
			expectError: true,
		},
		{
			name:        "invalid port 65536",
			port:        65536,
			expectError: true,
		},
		{
			name:        "invalid port negative",
			port:        -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				ServerPort: tt.port,
			}
			err := e.ValidateServerPort()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateServerPort() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateAuthCookieName(t *testing.T) {
	tests := []struct {
		name        string
		cookieName  string
		expectError bool
	}{
		{
			name:        "valid cookie name",
			cookieName:  "codebox_auth_token",
			expectError: false,
		},
		{
			name:        "empty cookie name",
			cookieName:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				AuthCookieName: tt.cookieName,
			}
			err := e.ValidateAuthCookieName()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateAuthCookieName() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateSubdomainAuthCookieName(t *testing.T) {
	tests := []struct {
		name        string
		cookieName  string
		expectError bool
	}{
		{
			name:        "valid subdomain cookie name",
			cookieName:  "subdomain_codebox_auth_token",
			expectError: false,
		},
		{
			name:        "empty subdomain cookie name",
			cookieName:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				SubdomainAuthCookieName: tt.cookieName,
			}
			err := e.ValidateSubdomainAuthCookieName()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateSubdomainAuthCookieName() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateUploadsPath(t *testing.T) {
	tests := []struct {
		name        string
		uploadsPath string
		expectError bool
	}{
		{
			name:        "valid path",
			uploadsPath: "./data",
			expectError: false,
		},
		{
			name:        "empty path",
			uploadsPath: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				UploadsPath: tt.uploadsPath,
			}
			err := e.ValidateUploadsPath()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateUploadsPath() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateCliBinariesPath(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		cliBinPath  string
		expectError bool
	}{
		{
			name:        "valid existing directory",
			cliBinPath:  tempDir,
			expectError: false,
		},
		{
			name:        "empty path",
			cliBinPath:  "",
			expectError: true,
		},
		{
			name:        "non-existent path",
			cliBinPath:  "/nonexistent/path",
			expectError: true,
		},
	}

	// Create a file to test if path is a directory
	tempFile := tempDir + "/test_file"
	f, _ := os.Create(tempFile)
	f.Close()

	tests = append(tests, struct {
		name        string
		cliBinPath  string
		expectError bool
	}{
		name:        "path is a file not directory",
		cliBinPath:  tempFile,
		expectError: true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				CliBinariesPath: tt.cliBinPath,
			}
			err := e.ValidateCliBinariesPath()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateCliBinariesPath() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateTemplatesFolder(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name            string
		templatesFolder string
		expectError     bool
	}{
		{
			name:            "valid existing directory",
			templatesFolder: tempDir,
			expectError:     false,
		},
		{
			name:            "empty path",
			templatesFolder: "",
			expectError:     true,
		},
		{
			name:            "non-existent path",
			templatesFolder: "/nonexistent/templates",
			expectError:     true,
		},
	}

	// Create a file to test if path is a directory
	tempFile := tempDir + "/test_file"
	f, _ := os.Create(tempFile)
	f.Close()

	tests = append(tests, struct {
		name            string
		templatesFolder string
		expectError     bool
	}{
		name:            "path is a file not directory",
		templatesFolder: tempFile,
		expectError:     true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				TemplatesFolder: tt.templatesFolder,
			}
			err := e.ValidateTemplatesFolder()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateTemplatesFolder() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateRunnerTokenHeader(t *testing.T) {
	tests := []struct {
		name        string
		tokenHeader string
		expectError bool
	}{
		{
			name:        "valid token header",
			tokenHeader: "X-Codebox-Runner-Token",
			expectError: false,
		},
		{
			name:        "empty token header",
			tokenHeader: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				RunnerTokenHeader: tt.tokenHeader,
			}
			err := e.ValidateRunnerTokenHeader()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateRunnerTokenHeader() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateRunnerTokenQueryParam(t *testing.T) {
	tests := []struct {
		name            string
		tokenQueryParam string
		expectError     bool
	}{
		{
			name:            "valid query param",
			tokenQueryParam: "runner_token",
			expectError:     false,
		},
		{
			name:            "empty query param",
			tokenQueryParam: "",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				RunnerTokenQueryParam: tt.tokenQueryParam,
			}
			err := e.ValidateRunnerTokenQueryParam()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateRunnerTokenQueryParam() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateDBDriver(t *testing.T) {
	tests := []struct {
		name        string
		driver      string
		expectError bool
	}{
		{
			name:        "valid mysql driver",
			driver:      "mysql",
			expectError: false,
		},
		{
			name:        "valid sqlite3 driver",
			driver:      "sqlite3",
			expectError: false,
		},
		{
			name:        "empty driver",
			driver:      "",
			expectError: true,
		},
		{
			name:        "unsupported postgresql driver",
			driver:      "postgresql",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				DBDriver: tt.driver,
			}
			err := e.ValidateDBDriver()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateDBDriver() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateDBHost(t *testing.T) {
	tests := []struct {
		name        string
		dbHost      string
		expectError bool
	}{
		{
			name:        "valid localhost",
			dbHost:      "localhost",
			expectError: false,
		},
		{
			name:        "valid ip address",
			dbHost:      "192.168.1.1",
			expectError: false,
		},
		{
			name:        "empty host",
			dbHost:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				DBHost: tt.dbHost,
			}
			err := e.ValidateDBHost()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateDBHost() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateDBPort(t *testing.T) {
	tests := []struct {
		name        string
		port        int
		expectError bool
	}{
		{
			name:        "valid port 3306",
			port:        3306,
			expectError: false,
		},
		{
			name:        "valid port 1",
			port:        1,
			expectError: false,
		},
		{
			name:        "valid port 65535",
			port:        65535,
			expectError: false,
		},
		{
			name:        "invalid port 0",
			port:        0,
			expectError: true,
		},
		{
			name:        "invalid port 65536",
			port:        65536,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				DBPort: tt.port,
			}
			err := e.ValidateDBPort()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateDBPort() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateDBName(t *testing.T) {
	tests := []struct {
		name        string
		dbName      string
		expectError bool
	}{
		{
			name:        "valid database name",
			dbName:      "codebox",
			expectError: false,
		},
		{
			name:        "empty database name",
			dbName:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				DBName: tt.dbName,
			}
			err := e.ValidateDBName()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateDBName() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateDBTestName(t *testing.T) {
	tests := []struct {
		name        string
		testDBName  string
		expectError bool
	}{
		{
			name:        "valid test database name",
			testDBName:  "codebox-test",
			expectError: false,
		},
		{
			name:        "empty test database name",
			testDBName:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				DBTestName: tt.testDBName,
			}
			err := e.ValidateDBTestName()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateDBTestName() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateDBUser(t *testing.T) {
	tests := []struct {
		name        string
		dbUser      string
		expectError bool
	}{
		{
			name:        "valid database user",
			dbUser:      "codebox",
			expectError: false,
		},
		{
			name:        "empty database user",
			dbUser:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				DBUser: tt.dbUser,
			}
			err := e.ValidateDBUser()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateDBUser() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateDBPassword(t *testing.T) {
	tests := []struct {
		name        string
		dbPassword  string
		expectError bool
	}{
		{
			name:        "valid password",
			dbPassword:  "secure_password",
			expectError: false,
		},
		{
			name:        "empty password",
			dbPassword:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				DBPassword: tt.dbPassword,
			}
			err := e.ValidateDBPassword()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateDBPassword() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateTasksConcurrency(t *testing.T) {
	tests := []struct {
		name             string
		tasksConcurrency int
		expectError      bool
	}{
		{
			name:             "valid concurrency 5",
			tasksConcurrency: 5,
			expectError:      false,
		},
		{
			name:             "valid concurrency 1",
			tasksConcurrency: 1,
			expectError:      false,
		},
		{
			name:             "invalid concurrency 0",
			tasksConcurrency: 0,
			expectError:      true,
		},
		{
			name:             "invalid concurrency negative",
			tasksConcurrency: -1,
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				TasksConcurrency: tt.tasksConcurrency,
			}
			err := e.ValidateTasksConcurrency()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateTasksConcurrency() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateRedisHost(t *testing.T) {
	tests := []struct {
		name        string
		redisHost   string
		expectError bool
	}{
		{
			name:        "valid redis host",
			redisHost:   "redis",
			expectError: false,
		},
		{
			name:        "empty redis host",
			redisHost:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				RedisHost: tt.redisHost,
			}
			err := e.ValidateRedisHost()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateRedisHost() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateRedisPort(t *testing.T) {
	tests := []struct {
		name        string
		port        int
		expectError bool
	}{
		{
			name:        "valid redis port 6379",
			port:        6379,
			expectError: false,
		},
		{
			name:        "valid port 1",
			port:        1,
			expectError: false,
		},
		{
			name:        "valid port 65535",
			port:        65535,
			expectError: false,
		},
		{
			name:        "invalid port 0",
			port:        0,
			expectError: true,
		},
		{
			name:        "invalid port 65536",
			port:        65536,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				RedisPort: tt.port,
			}
			err := e.ValidateRedisPort()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateRedisPort() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidateEmailSMTPPort(t *testing.T) {
	tests := []struct {
		name          string
		emailSMTPPort int
		expectError   bool
	}{
		{
			name:          "valid smtp port 587",
			emailSMTPPort: 587,
			expectError:   false,
		},
		{
			name:          "valid smtp port 25",
			emailSMTPPort: 25,
			expectError:   false,
		},
		{
			name:          "invalid port",
			emailSMTPPort: -1,
			expectError:   true,
		},
		{
			name:          "valid port 0 (not configured)",
			emailSMTPPort: 0,
			expectError:   false,
		},
		{
			name:          "invalid port 65536",
			emailSMTPPort: 65536,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EnvVars{
				EmailSMTPPort: tt.emailSMTPPort,
			}
			err := e.ValidateEmailSMTPPort()
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateEmailSMTPPort() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		expectValid bool
	}{
		{
			name:        "valid simple domain",
			domain:      "example.com",
			expectValid: true,
		},
		{
			name:        "valid subdomain",
			domain:      "sub.example.com",
			expectValid: true,
		},
		{
			name:        "valid complex domain",
			domain:      "sub-domain.example.co.uk",
			expectValid: true,
		},
		{
			name:        "invalid no dot",
			domain:      "localhost",
			expectValid: false,
		},
		{
			name:        "invalid too long label",
			domain:      "a" + string(make([]byte, 63)) + ".com",
			expectValid: false,
		},
		{
			name:        "invalid starts with hyphen",
			domain:      "-example.com",
			expectValid: false,
		},
		{
			name:        "invalid ends with hyphen",
			domain:      "example-.com",
			expectValid: false,
		},
		{
			name:        "valid with hyphens",
			domain:      "my-domain.example.com",
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidDomain(tt.domain)
			if valid != tt.expectValid {
				t.Errorf("isValidDomain(%s) = %v, expectValid %v", tt.domain, valid, tt.expectValid)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		envVars     *EnvVars
		expectError bool
	}{
		{
			name: "all valid environment variables",
			envVars: &EnvVars{
				ExternalUrl:             "https://example.com",
				UseSubDomains:           true,
				WildcardDomain:          "example.com",
				ServerPort:              8080,
				DebugEnabled:            true,
				AuthCookieName:          "codebox_auth_token",
				SubdomainAuthCookieName: "subdomain_codebox_auth_token",
				UploadsPath:             "./data",
				CliBinariesPath:         tempDir,
				TemplatesFolder:         tempDir,
				RunnerTokenHeader:       "X-Codebox-Runner-Token",
				RunnerTokenQueryParam:   "runner_token",
				DBDriver:                "mysql",
				DBHost:                  "localhost",
				DBPort:                  3306,
				DBName:                  "codebox",
				DBTestName:              "codebox-test",
				DBUser:                  "codebox",
				DBPassword:              "password",
				TasksConcurrency:        5,
				RedisHost:               "redis",
				RedisPort:               6379,
				EmailSMTPPort:           587,
			},
			expectError: false,
		},
		{
			name: "invalid server port",
			envVars: &EnvVars{
				ExternalUrl:             "https://example.com",
				UseSubDomains:           false,
				WildcardDomain:          "",
				ServerPort:              70000,
				DebugEnabled:            true,
				AuthCookieName:          "codebox_auth_token",
				SubdomainAuthCookieName: "subdomain_codebox_auth_token",
				UploadsPath:             "./data",
				CliBinariesPath:         tempDir,
				TemplatesFolder:         tempDir,
				RunnerTokenHeader:       "X-Codebox-Runner-Token",
				RunnerTokenQueryParam:   "runner_token",
				DBDriver:                "mysql",
				DBHost:                  "localhost",
				DBPort:                  3306,
				DBName:                  "codebox",
				DBTestName:              "codebox-test",
				DBUser:                  "codebox",
				DBPassword:              "password",
				TasksConcurrency:        5,
				RedisHost:               "redis",
				RedisPort:               6379,
			},
			expectError: true,
		},
		{
			name: "missing required cookie name",
			envVars: &EnvVars{
				ExternalUrl:             "https://example.com",
				UseSubDomains:           false,
				WildcardDomain:          "",
				ServerPort:              8080,
				DebugEnabled:            true,
				AuthCookieName:          "",
				SubdomainAuthCookieName: "subdomain_codebox_auth_token",
				UploadsPath:             "./data",
				CliBinariesPath:         tempDir,
				TemplatesFolder:         tempDir,
				RunnerTokenHeader:       "X-Codebox-Runner-Token",
				RunnerTokenQueryParam:   "runner_token",
				DBDriver:                "mysql",
				DBHost:                  "localhost",
				DBPort:                  3306,
				DBName:                  "codebox",
				DBTestName:              "codebox-test",
				DBUser:                  "codebox",
				DBPassword:              "password",
				TasksConcurrency:        5,
				RedisHost:               "redis",
				RedisPort:               6379,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.envVars.IsValid()
			if (err != nil) != tt.expectError {
				t.Errorf("IsValid() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestIsEmailConfigured(t *testing.T) {
	// Save current Environment
	oldEnv := Environment

	defer func() {
		Environment = oldEnv
	}()

	tests := []struct {
		name             string
		envVars          *EnvVars
		expectConfigured bool
	}{
		{
			name: "email fully configured",
			envVars: &EnvVars{
				EmailSMTPHost:     "smtp.example.com",
				EmailSMTPPort:     587,
				EmailSMTPUser:     "user@example.com",
				EmailSMTPPassword: "password123",
			},
			expectConfigured: true,
		},
		{
			name: "email not configured - empty host",
			envVars: &EnvVars{
				EmailSMTPHost:     "",
				EmailSMTPPort:     587,
				EmailSMTPUser:     "user@example.com",
				EmailSMTPPassword: "password123",
			},
			expectConfigured: false,
		},
		{
			name: "email not configured - zero port",
			envVars: &EnvVars{
				EmailSMTPHost:     "smtp.example.com",
				EmailSMTPPort:     0,
				EmailSMTPUser:     "user@example.com",
				EmailSMTPPassword: "password123",
			},
			expectConfigured: false,
		},
		{
			name: "email not configured - empty user",
			envVars: &EnvVars{
				EmailSMTPHost:     "smtp.example.com",
				EmailSMTPPort:     587,
				EmailSMTPUser:     "",
				EmailSMTPPassword: "password123",
			},
			expectConfigured: false,
		},
		{
			name: "email not configured - empty password",
			envVars: &EnvVars{
				EmailSMTPHost:     "smtp.example.com",
				EmailSMTPPort:     587,
				EmailSMTPUser:     "user@example.com",
				EmailSMTPPassword: "",
			},
			expectConfigured: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Environment = tt.envVars
			configured := IsEmailConfigured()
			if configured != tt.expectConfigured {
				t.Errorf("IsEmailConfigured() = %v, expectConfigured %v", configured, tt.expectConfigured)
			}
		})
	}
}
