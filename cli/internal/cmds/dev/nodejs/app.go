package nodejs

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/twirapp/twir/cli/internal/shell"
)

type twirApp struct {
	name   string
	cmd    *exec.Cmd
	path   string
	stdout *shell.PrefixWriter
	stderr *shell.PrefixWriter
}

func newApplication(name string) (*twirApp, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	app := twirApp{
		name:   name,
		cmd:    nil,
		path:   filepath.Join(wd, "apps", name),
		stdout: shell.StdoutFor(name),
		stderr: shell.StderrFor(name),
	}

	cmd, err := app.createAppCommand()
	if err != nil {
		return nil, err
	}
	app.cmd = cmd

	return &app, nil
}

func (c *twirApp) stop() error {
	if c.cmd == nil || c.cmd.Process == nil {
		return nil
	}

	pid := c.cmd.Process.Pid

	pgid, err := syscall.Getpgid(pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGTERM)
	} else {
		c.cmd.Process.Signal(syscall.SIGTERM)
	}

	done := make(chan error, 1)
	go func() {
		done <- c.cmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		if pgid != 0 {
			syscall.Kill(-pgid, syscall.SIGKILL)
		} else {
			c.cmd.Process.Kill()
		}
		<-done
	}

	if c.stdout != nil {
		c.stdout.Flush()
	}
	if c.stderr != nil {
		c.stderr.Flush()
	}

	c.cmd = nil
	return nil
}

func (c *twirApp) createAppCommand() (*exec.Cmd, error) {
	cmd, err := shell.CreateCommand(
		shell.ExecCommandOpts{
			Command: "bun run dev",
			Pwd:     c.path,
			Stdout:  c.stdout,
			Stderr:  c.stderr,
		},
	)

	if err != nil {
		return nil, err
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd, nil
}

func (c *twirApp) start() error {
	return c.cmd.Start()
}
