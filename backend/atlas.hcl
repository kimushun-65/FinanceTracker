env "dev" {
  url = "postgres://postgres:postgres@db:5432/financetracker_db?sslmode=disable"
  dir = "file://cmd/migrate/migrations"
}

env "prod" {
  url = env("DATABASE_URL")
  dir = "file://cmd/migrate/migrations"
}