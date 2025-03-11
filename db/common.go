package db

import (
	"fmt"
	"strings"

	"github.com/davidebianchi03/codebox/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectDB() error {
	dbUrl := config.Environment.DBUrl
	dbEngine := ""

	if strings.Index(dbUrl, "sqlite://") == 0 {
		dbEngine = "sqlite"
		dbUrl = strings.ReplaceAll(dbUrl, "sqlite://", "")
	}

	if dbEngine == "sqlite" {
		// Open database connection
		db, err := gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		DB = db
	} else {
		panic("unsupported db engine")
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
