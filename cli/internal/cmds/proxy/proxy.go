package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/pterm/pterm"
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

		if runtime.GOOS == "linux" {
			caddyFindCmd := exec.Command(
				"go",
				"tool",
				"-n",
				"github.com/caddyserver/caddy/v2/cmd/caddy",
			)
			caddyFindCmd.Dir = wd
			caddyFindCmdOutPut, err := caddyFindCmd.Output()
			if err != nil {
				return err
			}

			pterm.Warning.Println("!!! ATTENTION !!!")
			pterm.Info.Println("We need your sudo password to bind web server to port 443")

			if err := shell.ExecCommand(
				shell.ExecCommandOpts{
					Command: fmt.Sprintf(
						`sudo setcap 'cap_net_bind_service=+ep' %s`,
						string(caddyFindCmdOutPut),
					),
					Stdout: os.Stdout,
					Stderr: os.Stderr,
					Pwd:    wd,
				},
			); err != nil {
				return err
			}
		}

		return shell.ExecCommand(
			shell.ExecCommandOpts{
				Command: "go tool github.com/caddyserver/caddy/v2/cmd/caddy run --watch --config Caddyfile.dev --envfile .env",
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
				Pwd:     wd,
			},
		)
	},
}
