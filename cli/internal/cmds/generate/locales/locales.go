package locales

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	twiri18n "github.com/twirapp/twir/libs/i18n"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:    "locales",
	Usage:   "generate locales for application",
	Aliases: []string{"l"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "app",
			Aliases:  []string{"a"},
			Usage:    "application name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		app := c.String("app")
		pterm.Info.Printfln("Generating locales for application %s", app)

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %w", err)
		}

		localesDir := filepath.Join(wd, "apps", app, "locales")

		store, err := twiri18n.NewStore(localesDir)
		if err != nil {
			return fmt.Errorf("error creating locales store: %w", err)
		}

		content, err := twiri18n.GenerateKeysFileContent(
			twiri18n.GenerateKeysOptions{
				Locales:    store,
				Package:    "locales",
				BaseLocale: "en",
				LocalesDir: localesDir,
			},
		)
		if err != nil {
			return fmt.Errorf("error generating locales file content: %w", err)
		}

		err = os.WriteFile(filepath.Join(localesDir, "locales.go"), []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("error writing locales.go file: %w", err)
		}

		return nil
	},
}
