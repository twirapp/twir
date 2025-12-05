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
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "Name of the migration",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "db",
			Usage:    "Database type (postgres, clickhouse)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "type",
			Usage:    "Migration type (sql, go)",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		name, db, migrationType := c.String("name"), c.String("db"), c.String("type")

		if name == "" {
			migrationNameInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)
			migrationName, err := migrationNameInput.Show("Migration name")
			if err != nil {
				return err
			}

			name = strings.TrimSpace(migrationName)
		}

		if db == "" {
			migrationDb, err := pterm.DefaultInteractiveSelect.WithOptions(
				[]string{
					"postgres",
					"clickhouse",
				},
			).Show()
			if err != nil {
				return err
			}

			db = migrationDb
		}

		if db != "postgres" && db != "clickhouse" {
			return fmt.Errorf("invalid db type: %s", db)
		}

		if migrationType == "" {
			mType, err := pterm.DefaultInteractiveSelect.WithOptions([]string{"sql", "go"}).Show()
			if err != nil {
				return err
			}

			migrationType = mType
		}

		if migrationType != "sql" && migrationType != "go" {
			return fmt.Errorf("invalid migration type: %s", migrationType)
		}

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get working directory: %w", err)
		}

		dir := filepath.Join(wd, "libs", "migrations", db)

		log.SetOutput(&emptyLogWriter{})
		err = goose.Create(nil, dir, name, migrationType)
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

				fixedContent := strings.Replace(
					string(content),
					"package migrations",
					fmt.Sprintf("package %s", db),
					1,
				)
				err = os.WriteFile(migrationFilePath, []byte(fixedContent), 0644)
				if err != nil {
					return fmt.Errorf("cannot write fixed migration file: %w", err)
				}

				pterm.Info.Println("Fixed package name in migration file")
			}
		}

		pterm.Info.Println(fmt.Sprintf(`Migration "%s" created`, name))

		return nil
	},
}
