package connection

import (
	"fmt"

	"github.com/davidebianchi03/codebox/config"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectDB() error {
	// Open database connection
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

func IsItemInArray(item string, array []string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}
