data "external_schema" "codebox" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./db/models",
    "--dialect", "sqlite", // | postgres | sqlite | sqlserver
  ]
}

env "codebox" {
  src = data.external_schema.codebox.url
  dev = "sqlite://dev.db?_pragma=encoding=UTF-8"
  url = "sqlite://codebox.db?_pragma=encoding=UTF-8"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

# curl -sSf https://atlasgo.sh | sh
# go get ariga.io/atlas-go-sdk/atlasexec
# atlas migrate diff  --env codebox
# atlas migrate apply --env codebox --url "sqlite://test.db"
# atlas migrate apply --env codebox