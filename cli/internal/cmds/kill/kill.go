package kill

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:    "kill",
	Usage:   "kill any runed twir services (except docker)",
	Aliases: []string{"k"},
	Action: func(c *cli.Context) error {
		pterm.Info.Println("Killing")

		shell.ExecCommand(
			shell.ExecCommandOpts{
				Command: "kill -9 $(pgrep twir-)",
				Stderr:  os.Stderr,
				Stdout:  os.Stdout,
			},
		)

		return nil
	},
}
