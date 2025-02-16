package frontend

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/twirapp/twir/cli/internal/shell"
)

type twirApp struct {
	name string
	cmd  *exec.Cmd
	path string
}

func newApplication(name, path string) (*twirApp, error) {
	app := twirApp{
		name: name,
		cmd:  nil,
		path: path,
	}

	cmd, err := app.createAppCommand()
	if err != nil {
		return nil, err
	}
	app.cmd = cmd

	return &app, nil
}

func (c *twirApp) stop() error {
	if c.cmd != nil && c.cmd.Process != nil {
		if err := c.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			return err
		}
	}

	return nil
}

func (c *twirApp) createAppCommand() (*exec.Cmd, error) {
	cmd, err := shell.CreateCommand(
		shell.ExecCommandOpts{
			Command: "bun run dev",
			Pwd:     c.path,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		},
	)

	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (c *twirApp) start() error {
	return c.cmd.Start()
}
