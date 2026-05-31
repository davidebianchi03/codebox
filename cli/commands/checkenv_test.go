package commands

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"gitlab.com/codebox4073715/codebox/config"
)

var baseEnv = config.EnvVars{
	ExternalUrl:             "https://codebox.example.com",
	DebugEnabled:            false,
	DBDriver:                "mysql",
	DBHost:                  "mysql",
	DBPort:                  3306,
	DBName:                  "codebox",
	DBTestName:              "",
	DBUser:                  "codebox",
	DBPassword:              "password",
	ServerPort:              8080,
	TasksConcurrency:        5,
	RedisHost:               "",
	RedisPort:               6379,
	UploadsPath:             "",
	UseSubDomains:           true,
	WildcardDomain:          "codebox.example.com",
	AuthCookieName:          "",
	SubdomainAuthCookieName: "",
	CliBinariesPath:         "",
	TemplatesFolder:         "",
	RunnerTokenHeader:       "",
	RunnerTokenQueryParam:   "",
	EmailSMTPHost:           "",
	EmailSMTPPort:           587,
	EmailSMTPUser:           "",
	EmailSMTPPassword:       "",
}

/*
Try to check environment configuration with valid config
*/
func TestHandleCheckEnvValid(t *testing.T) {
	if err := ExportEnvVars(baseEnv); err != nil {
		t.Fatalf("failed to export env vars: %v", err)
	}

	status := HandleCheckEnv()

	if status != 0 {
		t.Errorf("expected status 0, got %d", status)
	}
}

/*
Try to check environment configuration with invalid config
*/
func TestHandleCheckEnvInvalid(t *testing.T) {
	invalidEnv := OverrideEnv(
		t,
		baseEnv,
		map[string]any{
			"DBDriver": "myfakedb",
		},
	)

	if err := ExportEnvVars(invalidEnv); err != nil {
		t.Fatalf("failed to export env vars: %v", err)
	}

	status := HandleCheckEnv()

	if status != 1 {
		t.Errorf("expected status 1, got %d", status)
	}
}

/*
Export environment variables
*/
func ExportEnvVars(cfg config.EnvVars) error {
	v := reflect.ValueOf(cfg)
	t := reflect.TypeOf(cfg)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("env")
		if tag == "" {
			continue
		}

		// supporta "NAME,required"
		parts := strings.Split(tag, ",")
		envName := parts[0]

		var strValue string

		switch value.Kind() {
		case reflect.String:
			strValue = value.String()

		case reflect.Int, reflect.Int32, reflect.Int64:
			strValue = strconv.FormatInt(value.Int(), 10)

		case reflect.Bool:
			strValue = strconv.FormatBool(value.Bool())

		default:
			return fmt.Errorf("unsupported type: %s", field.Name)
		}

		// se vuoto, prova envDefault
		if strValue == "" {
			if def := field.Tag.Get("envDefault"); def != "" {
				strValue = def
			}
		}

		// required check
		if strings.Contains(tag, "required") && strValue == "" {
			return fmt.Errorf("missing required env var: %s", envName)
		}

		if err := os.Setenv(envName, strValue); err != nil {
			return err
		}
	}

	return nil
}

/*
Changed the value of an environment variable
*/
func OverrideEnv(t *testing.T, env config.EnvVars, overrides map[string]any) config.EnvVars {
	v := reflect.ValueOf(&env).Elem()

	for fieldName, value := range overrides {
		field := v.FieldByName(fieldName)

		if !field.IsValid() {
			t.Error(fmt.Errorf("field %s does not exist", fieldName))
			return env
		}

		if !field.CanSet() {
			t.Error(fmt.Errorf("cannot set field %s", fieldName))
			return env
		}

		val := reflect.ValueOf(value)

		// handle type mismatch safely
		if val.Type() != field.Type() {
			if val.Type().ConvertibleTo(field.Type()) {
				val = val.Convert(field.Type())
			} else {
				t.Error(fmt.Errorf("invalid type for field %s", fieldName))
				return env
			}
		}

		field.Set(val)
	}

	return env
}
