package db

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB         *gorm.DB
	SqlDB      *sql.DB
	DBMigrator gorm.Migrator
)

func InitDBConnection(dbDriver string, dbUrl string) error {
	// connect to DB
	var err error
	if dbDriver == "mysql" {
		DB, err = gorm.Open(mysql.Open(dbUrl), &gorm.Config{
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

	// apply migrations
	if err := goose.SetDialect(dbDriver); err != nil {
		return err
	}

	migrationsFolder := "migrations"
	if err := goose.Up(SqlDB, migrationsFolder); err != nil {
		return err
	}

	return nil
}

func isItemInArray(item string, array []string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}
