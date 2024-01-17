package execbin

import (
	"fmt"
	"os"
	"strings"

	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "bin",
	Usage: "exec binary from .bin folder",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "pwd",
		},
	},
	Action: func(c *cli.Context) error {
		var pwd string
		if c.IsSet("pwd") {
			pwd = c.String("pwd")
		}

		fmt.Println(c.Args().Slice(), pwd)

		err := shell.ExecCommand(
			shell.ExecCommandOpts{
				Command: strings.Join(c.Args().Slice(), " "),
				Pwd:     pwd,
				Stdout:  os.Stdout,
				Stderr:  os.Stdin,
			},
		)

		return err
	},
}
