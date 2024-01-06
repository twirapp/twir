package proxy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "proxy",
	Usage: "Run https proxy",
	Action: func(context *cli.Context) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		caddyPath := filepath.Join(wd, ".bin", "caddy")

		return shell.ExecCommand(
			shell.ExecCommandOpts{
				Command: fmt.Sprintf(
					"%s reverse-proxy --from dev.twir.app --to 127.0.0.1:3005 --insecure --internal-certs",
					caddyPath,
				),
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			},
		)
	},
}
