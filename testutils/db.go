package testutils

import (
	"errors"

	"gitlab.com/codebox4073715/codebox/config"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
)

/*
Method that removes all data from all tables in the database
*/
func ClearDBTables() error {
	switch config.Environment.DBDriver {
	case "sqlite3":
		return errors.New("ClearDBTables is not supported for sqlite3")
	case "mysql":
		// Disable foreign key checks temporarily
		if err := dbconn.DB.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error; err != nil {
			return err
		}

		var tables []string
		if err := dbconn.DB.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
			return err
		}

		for _, table := range tables {
			// Truncate each table to remove all rows and reset AUTO_INCREMENT
			if err := dbconn.DB.Exec("TRUNCATE TABLE `" + table + "`;").Error; err != nil {
				return err
			}
		}

		// Re-enable foreign key checks
		if err := dbconn.DB.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error; err != nil {
			return err
		}
	}
	return nil
}
