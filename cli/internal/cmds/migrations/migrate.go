package migrations

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/pterm/pterm"
	cfg "github.com/satont/twir/libs/config"
	migrationsSeeds "github.com/satont/twir/libs/migrations/seeds"
	"github.com/urfave/cli/v2"
)

var MigrateCmd = &cli.Command{
	Name:  "run",
	Usage: "Run pending migrations",
	Action: func(c *cli.Context) error {
		pterm.Info.Println("Runing migrations")

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		config, err := cfg.NewWithEnvPath(filepath.Join(wd, ".env"))
		if err != nil {
			return err
		}

		opts, err := pq.ParseURL(config.DatabaseUrl)
		if err != nil {
			panic(err)
		}

		const driver = "postgres"
		db, err := sql.Open(driver, opts)
		if err != nil {
			return err
		}

		if err := goose.SetDialect(driver); err != nil {
			panic(err)
		}

		migrationsDir := filepath.Join(wd, "libs", "migrations", "migrations")

		log.SetOutput(&emptyLogWriter{})
		if err := goose.Up(db, migrationsDir); err != nil {
			pterm.Error.Println(err)
			return err
		}

		slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))

		if err := migrationsSeeds.CreateDefaultBot(db, config); err != nil {
			panic(err)
		}

		if err := migrationsSeeds.CreateIntegrations(db, config); err != nil {
			panic(err)
		}

		pterm.Success.Println("Migrations succeed")

		return nil
	},
}
