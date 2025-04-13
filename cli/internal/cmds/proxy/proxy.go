package proxy

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "proxy",
	Usage: "Run https proxy",
	Action: func(context *cli.Context) error {
		_, err := StartProxy()
		return err
	},
}

func StartProxy() (<-chan struct{}, error) {
	startChannel := make(chan struct{})

	wd, err := os.Getwd()
	if err != nil {
		return startChannel, err
	}

	caddyFindCmd := exec.Command(
		"go",
		"tool",
		"-n", // -n prints the command without running it, giving us the path
		"github.com/caddyserver/caddy/v2/cmd/caddy",
	)
	caddyFindCmd.Dir = wd
	caddyPathBytes, err := caddyFindCmd.Output()
	if err != nil {
		return startChannel, fmt.Errorf("failed to find Caddy path: %v", err)
	}
	caddyPath := strings.TrimSpace(string(caddyPathBytes))

	if runtime.GOOS == "linux" {
		// Check if the capability is already set
		getcapCmd := exec.Command("getcap", caddyPath)
		getcapCmd.Dir = wd
		getcapOutput, err := getcapCmd.Output()
		if err != nil {
			// If getcap fails (e.g., command not found), proceed cautiously
			pterm.Warning.Println("Could not check capabilities; assuming they need to be set")
		}

		// Check if cap_net_bind_service is present
		if !strings.Contains(string(getcapOutput), "cap_net_bind_service") {
			pterm.Warning.Println("!!! ATTENTION !!!")
			pterm.Warning.Println("We need your sudo password to bind web server to port 443 (this is a one-time setup)")

			// Set the capability if missing
			setcapCmd := fmt.Sprintf("sudo setcap 'cap_net_bind_service=+ep' %s", caddyPath)
			if err := shell.ExecCommand(
				shell.ExecCommandOpts{
					Command: setcapCmd,
					Stdout:  os.Stdout,
					Stderr:  os.Stderr,
					Pwd:     wd,
				},
			); err != nil {
				return startChannel, fmt.Errorf("failed to set capability: %v", err)
			}
			pterm.Success.Println("Capability set successfully; no further sudo prompts needed unless Caddy binary changes")
		} else {
			pterm.Info.Println("Caddy already has permission to bind to port 443; no sudo required")
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

	go func() {
		err = shell.ExecCommand(
			shell.ExecCommandOpts{
				Command: "go tool github.com/caddyserver/caddy/v2/cmd/caddy run --watch --config Caddyfile.dev --envfile .env",
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
				Pwd:     wd,
			},
		)
		if err != nil {
			panic(err)
		}
	}()

	return startChannel, err
}

func checkIsProxyStarted(port int) bool {
	_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	return err == nil
}
