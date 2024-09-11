package db

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB         *gorm.DB
	SqlDB      *sql.DB
	DBMigrator gorm.Migrator
)

func InitDBConnection() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("codebox.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database %s", err)
	}
	SqlDB, _ = DB.DB()
	DBMigrator = DB.Migrator()
	return nil
}
