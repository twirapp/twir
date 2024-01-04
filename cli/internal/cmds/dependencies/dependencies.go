package dependencies

import (
	"github.com/urfave/cli/v2"
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
			if err := installGolangDeps(); err != nil {
				return err
			}
		}

		if !skipNode {
			if err := installNodeDeps(); err != nil {
				return err
			}
		}

		return nil
	},
}
