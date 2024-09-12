export GOOSE_DRIVER='mysql'
export GOOSE_DBSTRING='root:password@tcp(pc-taverna.fritz.box:3306)/yourdb?charset=utf8mb4&parseTime=True&loc=Local'
export GOOSE_MIGRATION_DIR='./migrations'
goose -s create create_user