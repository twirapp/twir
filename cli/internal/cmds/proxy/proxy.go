package proxy

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

		if runtime.GOOS == "windows" {
			caddyPath += ".exe"
		}

		return shell.ExecCommand(
			shell.ExecCommandOpts{
				Command: fmt.Sprintf(
					"%s reverse-proxy --from twir.localhost --to 127.0.0.1:3005",
					caddyPath,
				),
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			},
		)
	},
}
