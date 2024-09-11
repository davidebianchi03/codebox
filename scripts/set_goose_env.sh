export GOOSE_DRIVER='sqlite3'
export GOOSE_DBSTRING='sqlite://codebox.db'
export GOOSE_MIGRATION_DIR='./migrations'
goose -s create create_user