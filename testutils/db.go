package testutils

import "gorm.io/gorm"

/*
Method that removes all data from all tables in the database
*/
func ClearDB(db *gorm.DB) error {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return err
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	for _, table := range tables {
		err := db.Exec("TRUNCATE TABLE `" + table + "`;").Error
		if err != nil {
			return err
		}
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	return nil
}
