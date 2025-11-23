package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

		// fix package name in go migrations
		if migrationType == "go" {
			files, err := os.ReadDir(dir)
			if err != nil {
				return fmt.Errorf("cannot read migrations directory: %w", err)
			}

			var latestMigrationFile string
			for _, file := range files {
				if file.IsDir() || filepath.Ext(file.Name()) != ".go" {
					continue
				}

				latestMigrationFile = file.Name()
			}

			if latestMigrationFile != "" {
				migrationFilePath := filepath.Join(dir, latestMigrationFile)
				content, err := os.ReadFile(migrationFilePath)
				if err != nil {
					return fmt.Errorf("cannot read migration file: %w", err)
				}

				fixedContent := strings.Replace(string(content), "package migrations", fmt.Sprintf("package %s", migrationDb), 1)
				err = os.WriteFile(migrationFilePath, []byte(fixedContent), 0644)
				if err != nil {
					return fmt.Errorf("cannot write fixed migration file: %w", err)
				}

				pterm.Info.Println("Fixed package name in migration file")
			}
		}

		pterm.Info.Println(fmt.Sprintf(`Migration "%s" created`, migrationName))

		return nil
	},
}
