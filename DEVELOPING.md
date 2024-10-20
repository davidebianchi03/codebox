# Developing

Here you can find some getting started guides if you want to contribute to codebox.

## Create new DB migration

In order to create new DB migration follow this steps:

- Export the following environment variables:
```bash
export GOOSE_DRIVER='mysql'
export GOOSE_DBSTRING='root:password@tcp(pc-taverna.fritz.box:3306)/yourdb?charset=utf8mb4&parseTime=True&loc=Local'
export GOOSE_MIGRATION_DIR='./migrations'
```

- Now you can create new migration running:
```bash
goose -s create <custom_name_for_migration>
```