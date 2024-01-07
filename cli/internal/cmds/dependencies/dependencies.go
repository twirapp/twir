package dependencies

import (
	"github.com/pterm/pterm"
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

		if !skipGo {
			goSpinner, _ := pterm.DefaultSpinner.Start("Install golang deps...")
			if err := installGolangDeps(); err != nil {
				goSpinner.Fail(err)
				return err
			}
			goSpinner.Success("Golang deps installed")
		}

		if !skipNode {
			nodeSpinner, _ := pterm.DefaultSpinner.Start("Install node deps...")
			if err := installNodeDeps(); err != nil {
				nodeSpinner.Fail(err)
				return err
			}
			nodeSpinner.Success("Nodejs deps installed")
		}

		var errwg errgroup.Group

		binariesSpinner, _ := pterm.DefaultSpinner.Start("Install binaries...")

		errwg.Go(installGoBinaries)
		errwg.Go(installProtoc)

		if err := errwg.Wait(); err != nil {
			binariesSpinner.Fail(err)
			return err
		}

		binariesSpinner.Success("Binaries installed")

		return nil
	},
}
