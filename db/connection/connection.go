package connection

import (
	"fmt"

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
	if config.Environment.DBDriver == "sqlite3" {
		db, err := gorm.Open(sqlite.Open(config.Environment.DBName), &gorm.Config{})
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
			config.Environment.DBName,
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
