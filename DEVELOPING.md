# Developing

Here you can find some getting started guides if you want to contribute to codebox.

## Create new DB migration

In order to create new DB migration follow this steps:

- Export the following environment variables:
```bash
export GOOSE_DRIVER='sqlite3'
export GOOSE_DBSTRING='sqlite://codebox.db'
export GOOSE_MIGRATION_DIR='./migrations'
```

- Now you can create new migration running:
```bash
goose -s create <custom_name_for_migration>
```