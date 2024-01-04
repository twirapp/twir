package proxy

import (
	"os"

	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "proxy",
	Usage: "Run https proxy",
	Action: func(context *cli.Context) error {
		return shell.ExecCommand(
			shell.ExecCommandOpts{
				Command: "caddy reverse-proxy --from dev.twir.app --to 127.0.0.1:3005 --insecure --internal-certs",
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
			},
		)
	},
}
