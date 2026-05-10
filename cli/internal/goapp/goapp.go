package goapp

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

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
	DebugPort    int
	mu           sync.Mutex
	stdout       *shell.PrefixWriter
	stderr       *shell.PrefixWriter
}

func NewApplication(name string, enableDebug bool, port *int, debugPort int, onPortReady func()) (
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
		DebugPort:    debugPort,
		stdout:       shell.StdoutFor(name),
		stderr:       shell.StderrFor(name),
	}

	cmd, err := app.CreateAppCommand()
	if err != nil {
		return nil, err
	}
	app.Cmd = cmd

	return &app, nil
}

func (c *TwirGoApp) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stopLocked()
	return nil
}

func (c *TwirGoApp) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Cmd != nil {
		c.stopLocked()
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

func (c *TwirGoApp) stopLocked() {
	if c.Cmd == nil || c.Cmd.Process == nil {
		return
	}

	pid := c.Cmd.Process.Pid

	pgid, err := syscall.Getpgid(pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGTERM)
	} else {
		c.Cmd.Process.Signal(syscall.SIGTERM)
	}

	done := make(chan error, 1)
	go func() {
		done <- c.Cmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		if pgid != 0 {
			syscall.Kill(-pgid, syscall.SIGKILL)
		} else {
			c.Cmd.Process.Kill()
		}
		<-done
	}

	orphanKill := exec.Command("pkill", "-9", "-f", c.getAppPath())
	orphanKill.Run()

	if c.stdout != nil {
		c.stdout.Flush()
	}
	if c.stderr != nil {
		c.stderr.Flush()
	}

	c.Cmd = nil
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
	// dlv exec .out/twir-emotes-cacher --headless=true --api-version=2 --check-go-version=false --only-same-user=false --listen=:2345 --log

	cmd, err := shell.CreateCommand(
		shell.ExecCommandOpts{
			Command: fmt.Sprintf(
				"go tool dlv exec %s --headless=true --api-version=2 --check-go-version=false --only-same-user=false --listen=:%d --log --continue --accept-multiclient",
				c.getAppPath(),
				c.DebugPort,
			),
			Pwd:    c.Path,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
	)

	if err != nil {
		return nil, err
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd, nil
}

func (c *TwirGoApp) waitPortReady() {
	if c.Port == nil || c.OnPortReady == nil {
		return
	}

	for {
		_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", *c.Port))
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		pterm.Info.Println("Port " + c.Name + " is ready, running hook")
		c.OnPortReady()
		break
	}
}
