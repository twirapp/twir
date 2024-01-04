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
	Name:  "create",
	Usage: "Create new migration",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "golang",
			Usage:   "create golang migration",
			Aliases: []string{"g"},
		},
	},
	ArgsUsage:              "migrationName",
	UseShortOptionHandling: true,
	Action: func(c *cli.Context) error {
		arg := c.Args().Get(0)
		if arg == "" {
			return fmt.Errorf("name of migration cannot be empty")
		}

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get working directory: %w", err)
		}

		dir := filepath.Join(wd, "libs", "migrations", "migrations")

		fmt.Println(c.FlagNames(), c.Bool("golang"))
		migrationType := "sql"
		if c.Bool("golang") {
			migrationType = "go"
		}

		log.SetOutput(&emptyLogWriter{})
		err = goose.Create(nil, dir, arg, migrationType)
		if err != nil {
			return fmt.Errorf("cannot create migration: %w", err)
		}

		pterm.Info.Println(fmt.Sprintf(`Migration "%s" created`, arg))

		return nil
	},
}
