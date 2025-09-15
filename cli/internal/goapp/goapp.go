package goapp

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/twirapp/twir/cli/internal/watcher"
)

type TwirGoApp struct {
	Name         string
	Cmd          *exec.Cmd
	Path         string
	Watcher      *watcher.Watcher
	debugEnabled bool
	Port         *int
	OnPortReady  func()
}

func NewApplication(name string, enableDebug bool, port *int, onPortReady func()) (
	*TwirGoApp,
	error,
) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	app := TwirGoApp{
		Name:         name,
		Cmd:          nil,
		Path:         filepath.Join(wd, "apps", name),
		Watcher:      watcher.New(),
		debugEnabled: enableDebug,
		Port:         port,
		OnPortReady:  onPortReady,
	}

	cmd, err := app.CreateAppCommand()
	if err != nil {
		return nil, err
	}
	app.Cmd = cmd

	return &app, nil
}

func (c *TwirGoApp) Stop() error {
	if c.Cmd != nil && c.Cmd.Process != nil {
		if err := c.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
			return err
		}
	}

	return nil
}

func (c *TwirGoApp) Start() error {
	if err := c.Stop(); err != nil {
		return err
	}

	if err := c.Build(); err != nil {
		return err
	}

	newCmd, err := c.CreateAppCommand()
	if err != nil {
		return err
	}

	c.Cmd = newCmd

	go func() {
		c.waitPortReady()
	}()

	return c.Cmd.Start()
}

func (c *TwirGoApp) getAppPath() string {
	return filepath.Join(c.Path, ".out", "twir-"+c.Name)
}

func (c *TwirGoApp) Build() error {
	pterm.Info.Println(
		fmt.Sprintf(
			"Building %s with debug = %v",
			c.Name,
			c.debugEnabled,
		),
	)

	args := []string{"build", "-o", c.getAppPath()}
	if c.debugEnabled {
		args = append(args, `-gcflags=all=-N -l`)
	} else {
		args = append(args, "-ldflags=-s -w")
	}
	args = append(args, "./cmd/main.go")

	buildCmd := exec.Command("go", args...)
	buildCmd.Dir = c.Path
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	if err := buildCmd.Run(); err != nil {
		return err
	}

	return nil
}

func (c *TwirGoApp) CreateAppCommand() (*exec.Cmd, error) {
	cmd, err := shell.CreateCommand(
		shell.ExecCommandOpts{
			Command: c.getAppPath(),
			Pwd:     c.Path,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		},
	)

	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (c *TwirGoApp) waitPortReady() {
	if c.Port == nil || c.OnPortReady == nil {
		return
	}

	for {
		_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", *c.Port))
		if err != nil {
			continue
		}

		pterm.Info.Println("Port " + c.Name + " is ready, running hook")
		c.OnPortReady()
		break
	}
}
