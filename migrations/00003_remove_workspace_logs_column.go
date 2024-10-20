package migrations

import (
	"context"
	"database/sql"

	"codebox.com/db"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRemoveWorkspaceLogsColumn, downRemoveWorkspaceLogsColumn)
}

func upRemoveWorkspaceLogsColumn(ctx context.Context, tx *sql.Tx) error {
	// drop column 'logs'
	err := db.DBMigrator.DropColumn(&db.Workspace{}, "logs")
	if err != nil {
		return err
	}
	return nil
}

func downRemoveWorkspaceLogsColumn(ctx context.Context, tx *sql.Tx) error {
	// restore column 'logs'
	err := db.DBMigrator.AddColumn(&db.Workspace{}, "logs")
	if err != nil {
		return err
	}
	return nil
}
