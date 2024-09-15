package migrations

import (
	"context"
	"database/sql"

	"codebox.com/db"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateWorkspaceAndRelatedTables, downCreateWorkspaceAndRelatedTables)
}

func upCreateWorkspaceAndRelatedTables(ctx context.Context, tx *sql.Tx) error {
	// creazione della tabella workspace
	err := db.DBMigrator.CreateTable(&db.Workspace{})
	if err != nil {
		return err
	}

	// creazione della tabella workspace container
	err = db.DBMigrator.CreateTable(&db.WorkspaceContainer{})
	if err != nil {
		return err
	}

	// creazione della tabella workspace forwarded port
	err = db.DBMigrator.CreateTable(&db.ForwardedPort{})
	if err != nil {
		return err
	}

	return nil
}

func downCreateWorkspaceAndRelatedTables(ctx context.Context, tx *sql.Tx) error {
	// drop della tabella workspace
	err := db.DBMigrator.DropTable(&db.Workspace{})
	if err != nil {
		return err
	}

	// drop della tabella workspace container
	err = db.DBMigrator.DropTable(&db.WorkspaceContainer{})
	if err != nil {
		return err
	}

	// drop della tabella forwarded port
	err = db.DBMigrator.DropTable(&db.ForwardedPort{})
	if err != nil {
		return err
	}
	return nil
}
