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

		var wg errgroup.Group
		if !skipGo {
			wg.Go(
				func() error {
					if err := installGolangDeps(); err != nil {
						pterm.Fatal.Println(err)
						return err
					}
					pterm.Success.Println("Golang deps installed")
					return nil
				},
			)
		} else {
			go pterm.Warning.Println("Golang deps skipped")
		}

		if !skipNode {
			wg.Go(
				func() error {
					if err := installNodeDeps(); err != nil {
						pterm.Fatal.Println(err)
						return err
					}
					pterm.Success.Println("Nodejs deps installed")
					return nil
				},
			)
		} else {
			go pterm.Warning.Println("Nodejs deps skipped")
		}

		wg.Go(
			func() error {
				var binariesWg errgroup.Group

				if err := binaries.CreateDir(); err != nil {
					return err
				}

				binariesWg.Go(binaries.InstallProtoc)

				if err := binariesWg.Wait(); err != nil {
					pterm.Fatal.Println(err)
					return err
				}

				pterm.Success.Println("Binaries installed")
				return nil
			},
		)

		if err := wg.Wait(); err != nil {
			return err
		}

		pterm.Success.Println("All dependencies installed")

		return nil
	},
}
