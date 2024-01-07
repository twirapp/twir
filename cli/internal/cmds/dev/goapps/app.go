package goapps

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/twirapp/twir/cli/internal/watcher"
)

type twirApp struct {
	name    string
	cmd     *exec.Cmd
	path    string
	watcher *watcher.Watcher
}

func newApplication(name string) (*twirApp, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	app := twirApp{
		name:    name,
		cmd:     nil,
		path:    filepath.Join(wd, "apps", name),
		watcher: watcher.New(),
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

func (c *twirApp) start() error {
	if err := c.stop(); err != nil {
		return err
	}

	if err := c.build(); err != nil {
		return err
	}

	newCmd, err := c.createAppCommand()
	if err != nil {
		return err
	}

	c.cmd = newCmd
	return c.cmd.Start()
}

func (c *twirApp) getTempPath() string {
	tmp := os.TempDir()
	return filepath.Join(tmp, "twir-"+c.name)
}

func (c *twirApp) build() error {
	tmpFilePath := c.getTempPath()

	buildCmd := exec.Command("go", "build", "-o", tmpFilePath, "./cmd/main.go")
	buildCmd.Dir = c.path

	if err := buildCmd.Run(); err != nil {
		return err
	}

	return nil
}

func (c *twirApp) createAppCommand() (*exec.Cmd, error) {
	cmd, err := shell.CreateCommand(
		shell.ExecCommandOpts{
			Command: c.getTempPath(),
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
