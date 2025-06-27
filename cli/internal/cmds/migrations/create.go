package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var createCmd = &cli.Command{
	Name:                   "create",
	Usage:                  "Create new migration",
	UseShortOptionHandling: true,
	Action: func(c *cli.Context) error {
		migrationNameInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)
		migrationName, err := migrationNameInput.Show("Migration name")
		if err != nil {
			return err
		}

		migrationDb, err := pterm.DefaultInteractiveSelect.WithOptions(
			[]string{
				"postgres",
				"clickhouse",
			},
		).Show()
		if err != nil {
			return err
		}

		migrationType, err := pterm.DefaultInteractiveSelect.WithOptions([]string{"sql", "go"}).Show()
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get working directory: %w", err)
		}

		dir := filepath.Join(wd, "libs", "migrations", migrationDb)

		log.SetOutput(&emptyLogWriter{})
		err = goose.Create(nil, dir, migrationName, migrationType)
		if err != nil {
			return fmt.Errorf("cannot create migration: %w", err)
		}

		pterm.Info.Println(fmt.Sprintf(`Migration "%s" created`, migrationName))

		return nil
	},
}
