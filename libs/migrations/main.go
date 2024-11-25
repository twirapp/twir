package main

import (
	"database/sql"
	"embed"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/satont/twir/libs/config"
	_ "github.com/satont/twir/libs/migrations/migrations"
	"github.com/satont/twir/libs/migrations/seeds"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const driver = "postgres"

func main() {
	config, err := cfg.New()
	if err != nil {
		panic(err)
	}

	opts, err := pq.ParseURL(config.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(driver, opts)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(driver); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations", goose.WithAllowMissing()); err != nil {
		panic(err)
	}

	if err := seeds.CreateDefaultBot(db, config); err != nil {
		panic(err)
	}

	if err := seeds.CreateIntegrations(db, config); err != nil {
		panic(err)
	}
}
