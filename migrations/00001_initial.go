package migrations

import (
	"context"
	"database/sql"

	"codebox.com/db"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUser, downCreateUser)
}

func upCreateUser(ctx context.Context, tx *sql.Tx) error {
	// creazione della tabella user
	err := db.DBMigrator.CreateTable(&db.User{})
	if err != nil {
		return err
	}

	// creazione della tabella token
	err = db.DBMigrator.CreateTable(&db.Token{})
	if err != nil {
		return err
	}

	return nil
}

func downCreateUser(ctx context.Context, tx *sql.Tx) error {
	// drop della tabella user
	err := db.DBMigrator.DropTable(&db.User{})
	if err != nil {
		return err
	}

	// drop della tabella token
	err = db.DBMigrator.DropTable(&db.Token{})
	if err != nil {
		return err
	}

	return nil
}
