package db

import (
	"database/sql"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB         *gorm.DB
	SqlDB      *sql.DB
	DBMigrator gorm.Migrator
)

func InitDBConnection() error {
	dbDriver := os.Getenv("CODEBOX_DB_DRIVER")
	var err error
	if dbDriver == "mysql" {
		DB, err = gorm.Open(mysql.Open(os.Getenv("CODEBOX_DB_URL")), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		return fmt.Errorf("unknown db driver %s", dbDriver)
	}
	if err != nil {
		return fmt.Errorf("failed to connect database %s", err)
	}
	SqlDB, _ = DB.DB()
	DBMigrator = DB.Migrator()
	return nil
}
