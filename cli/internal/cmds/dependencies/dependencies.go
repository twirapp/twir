package dependencies

import (
	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/cmds/dependencies/binaries"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

var Cmd = &cli.Command{
	Name:    "dependencies",
	Usage:   "install golang and nodejs dependencies",
	Aliases: []string{"deps"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "skip-node",
			Value: false,
			Usage: "skip nodejs dependencies installation",
		},
		&cli.BoolFlag{
			Name:  "skip-go",
			Value: false,
			Usage: "skip golang dependencies installation",
		},
	},
	Action: func(c *cli.Context) error {
		skipNode := c.Bool("skip-node")
		skipGo := c.Bool("skip-go")

		multiPrinter := pterm.DefaultMultiPrinter
		goSpinner, _ := pterm.DefaultSpinner.
			WithWriter(multiPrinter.NewWriter()).
			Start("Install golang deps...")
		nodeSpinner, _ := pterm.DefaultSpinner.
			WithWriter(multiPrinter.NewWriter()).
			Start("Install node deps...")
		binariesSpinner, _ := pterm.DefaultSpinner.
			WithWriter(multiPrinter.NewWriter()).
			Start("Install binaries...")

		if _, err := multiPrinter.Start(); err != nil {
			return err
		}

		var wg errgroup.Group
		if !skipGo {
			wg.Go(
				func() error {
					if err := installGolangDeps(); err != nil {
						goSpinner.Fail(err)
						return err
					}
					goSpinner.Success("Golang deps installed")
					return nil
				},
			)
		} else {
			go goSpinner.Warning("Golang deps skipped")
		}

		if !skipNode {
			wg.Go(
				func() error {
					if err := installNodeDeps(); err != nil {
						nodeSpinner.Fail(err)
						return err
					}
					nodeSpinner.Success("Nodejs deps installed")
					return nil
				},
			)
		} else {
			go nodeSpinner.Warning("Nodejs deps skipped")
		}

		wg.Go(
			func() error {
				var binariesWg errgroup.Group

				if err := binaries.CreateDir(); err != nil {
					return err
				}

				binariesWg.Go(binaries.InstallGoBinaries)
				binariesWg.Go(binaries.InstallProtoc)

				if err := binariesWg.Wait(); err != nil {
					binariesSpinner.Fail(err)
					return err
				}

				binariesSpinner.Success("Binaries installed")
				return nil
			},
		)

		if err := wg.Wait(); err != nil {
			return err
		}

		if _, err := multiPrinter.Stop(); err != nil {
			return err
		}

		return nil
	},
}
