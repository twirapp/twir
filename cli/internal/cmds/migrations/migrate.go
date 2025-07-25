package migrations

import (
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/cli/internal/cmds/migrations/clickhouse"
	"github.com/twirapp/twir/cli/internal/cmds/migrations/postgres"
	"github.com/urfave/cli/v2"
)

var MigrateCmd = &cli.Command{
	Name:  "run",
	Usage: "Run pending migrations",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "migrations-path",
			Value: "./libs/migrations",
		},
		&cli.BoolFlag{
			Name:  "skip-clickhouse",
			Value: false,
			Usage: "skip clickhouse migrations",
		},
		&cli.BoolFlag{
			Name:  "skip-postgres",
			Value: false,
			Usage: "skip postgres migrations",
		},
	},
	Action: func(c *cli.Context) error {
		pterm.Info.Println("Running migrations")

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		migrationsPath := filepath.Join(wd, c.String("migrations-path"))

		config, err := cfg.NewWithEnvPath(filepath.Join(wd, ".env"))
		if err != nil {
			return err
		}

		if !c.Bool("skip-postgres") {
			pterm.Info.Println("Running postgres migrations...")
			if err := postgres.Migrate(
				c.Context,
				config,
				filepath.Join(migrationsPath, "postgres"),
			); err != nil {
				return err
			}
			pterm.Success.Println("Postgres migrations succeed")

		} else {
			pterm.Warning.Println("Postgres migrations skipped")
		}

		if !c.Bool("skip-clickhouse") {
			pterm.Info.Println("Running clickhouse migrations...")
			if err := clickhouse.Migrate(
				c.Context,
				config,
				filepath.Join(migrationsPath, "clickhouse"),
			); err != nil {
				return err
			}
			pterm.Success.Println("Clickhouse migrations succeed")
		} else {
			pterm.Warning.Println("Clickhouse migrations skipped")
		}

		pterm.Success.Println("Migration succeed")
		return nil
	},
}
