data "external_schema" "codebox" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./db/models",
    "--dialect", "mysql", // | postgres | sqlite | sqlserver
  ]
}

env "codebox" {
  src = data.external_schema.codebox.url
  dev = "mysql://${getenv("CODEBOX_DB_USER")}:${getenv("CODEBOX_DB_PASSWORD")}@${getenv("CODEBOX_DB_HOST")}:${getenv("CODEBOX_DB_PORT")}/${getenv("CODEBOX_DB_NAME")}-dev?charset=utf8mb4&parseTime=true"
  url = "mysql://${getenv("CODEBOX_DB_USER")}:${getenv("CODEBOX_DB_PASSWORD")}@${getenv("CODEBOX_DB_HOST")}:${getenv("CODEBOX_DB_PORT")}/${getenv("CODEBOX_DB_NAME")}?charset=utf8mb4&parseTime=true"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
  exclude  = [
    "singleton_models"
  ]
}

# curl -sSf https://atlasgo.sh | sh
# go get ariga.io/atlas-go-sdk/atlasexec
# atlas migrate diff  --env codebox
# atlas migrate apply --env codebox --url "sqlite://test.db"
# atlas migrate apply --env codebox
# atlas migrate down --env codebox
