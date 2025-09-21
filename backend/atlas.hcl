env "dev" {
  url = "postgres://postgres:postgres@db:5432/financetracker_db?sslmode=disable"
  src = "file:///app/cmd/migrate/schema.hcl"
  dir = "file:///app/cmd/migrate/migrations"
  dev = "docker://postgres/15/dev"
  migration {
    exclude = ["atlas_schema_revisions"]
  }
}

env "prod" {
  url = env("DATABASE_URL")
  dir = "file://cmd/migrate/migrations"
}