package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectDB() error {
	dbURL := "./codebox.db"
	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	DB = db
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
