package connection

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/codebox4073715/codebox/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// Open connection with db
// Codebox supports sqlite3 and mysql as dbms
// Connection will be stored in DB var and will be
// accessible from any point of the code
func ConnectDB() error {
	dbName := config.Environment.DBName

	// if we are running tests, use the test db name
	if flag.Lookup("test.v") != nil || strings.HasSuffix(os.Args[0], ".test") {
		dbName = config.Environment.DBTestName
	}

	if config.Environment.DBDriver == "sqlite3" {
		db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		DB = db
	} else if config.Environment.DBDriver == "mysql" {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
			config.Environment.DBUser,
			config.Environment.DBPassword,
			config.Environment.DBHost,
			config.Environment.DBPort,
			dbName,
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		DB = db
	} else {
		return errors.New("unsupported db engine")
	}
	return nil
}

// Close connection with db
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
