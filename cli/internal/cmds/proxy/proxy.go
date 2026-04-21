package proxy

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "proxy",
	Usage: "Run https proxy",
	Action: func(context *cli.Context) error {
		_, _, err := StartProxy(true)
		return err
	},
}

type Proxy struct {
	cmd  *exec.Cmd
	done chan struct{}
}

func (p *Proxy) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	pid := p.cmd.Process.Pid

	pgid, err := syscall.Getpgid(pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGTERM)
	} else {
		p.cmd.Process.Signal(syscall.SIGTERM)
	}

	select {
	case <-p.done:
	case <-time.After(5 * time.Second):
		if pgid != 0 {
			syscall.Kill(-pgid, syscall.SIGKILL)
		} else {
			p.cmd.Process.Kill()
		}
		<-p.done
	}

	p.cmd = nil
	return nil
}

func StartProxy(block bool) (<-chan struct{},
	*Proxy,
	error,
) {
	startChannel := make(chan struct{})

	wd, err := os.Getwd()
	if err != nil {
		return startChannel, nil, err
	}

	caddyFindCmd := exec.Command(
		"go",
		"tool",
		"-n",
		"github.com/caddyserver/caddy/v2/cmd/caddy",
	)
	caddyFindCmd.Dir = wd
	caddyPathBytes, err := caddyFindCmd.Output()
	if err != nil {
		return startChannel, nil, fmt.Errorf("failed to find Caddy path: %v", err)
	}
	caddyPath := strings.TrimSpace(string(caddyPathBytes))

	if runtime.GOOS == "linux" {
		getcapCmd := exec.Command("getcap", caddyPath)
		getcapCmd.Dir = wd
		getcapOutput, err := getcapCmd.Output()
		if err != nil {
			pterm.Warning.Println("Could not check capabilities; assuming they need to be set")
		}

		if !strings.Contains(string(getcapOutput), "cap_net_bind_service") {
			pterm.Warning.Println("!!! ATTENTION !!!")
			pterm.Warning.Println("We need your sudo password to bind web server to port 443 (this is a one-time setup)")

			setcapCmd := fmt.Sprintf("sudo setcap 'cap_net_bind_service=+ep' %s", caddyPath)
			if err := shell.ExecCommand(
				shell.ExecCommandOpts{
					Command: setcapCmd,
					Stdout:  os.Stdout,
					Stderr:  os.Stderr,
					Pwd:     wd,
				},
			); err != nil {
				return startChannel, nil, fmt.Errorf("failed to set capability: %v", err)
			}
			pterm.Success.Println("Capability set successfully; no further sudo prompts needed unless Caddy binary changes")
		}
	}

	go func() {
		for !checkIsProxyStarted(80) {
			pterm.Info.Println("Waiting for proxy to start")
			time.Sleep(500 * time.Millisecond)
		}

		pterm.Success.Println("Proxy started")
		startChannel <- struct{}{}
		close(startChannel)
	}()

	cmd, err := shell.CreateCommand(
		shell.ExecCommandOpts{
			Command: "go tool github.com/caddyserver/caddy/v2/cmd/caddy run --watch --config Caddyfile.dev --envfile .env",
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
			Pwd:     wd,
		},
	)
	if err != nil {
		return startChannel, nil, err
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	p := &Proxy{cmd: cmd, done: make(chan struct{})}

	if block {
		if err := cmd.Start(); err != nil {
			return startChannel, p, err
		}
		if err := cmd.Wait(); err != nil {
			close(p.done)
			return startChannel, p, err
		}
		close(p.done)
	} else {
		if err := cmd.Start(); err != nil {
			return startChannel, p, err
		}
		go func() {
			err := cmd.Wait()
			close(p.done)
			if err != nil && !isSignalTermination(err) {
				pterm.Error.Println("Proxy exited with error:", err)
			}
		}()
	}

	return startChannel, p, nil
}

func isSignalTermination(err error) bool {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	return exitErr.ProcessState != nil && exitErr.ProcessState.ExitCode() < 0
}

func checkIsProxyStarted(port int) bool {
	_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	return err == nil
}
