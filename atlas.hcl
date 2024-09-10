data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    ".",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://codebox.db"
  migration {
    dir = "file://db/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}