package tunnel

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"regexp"
	"strings"
	"time"

	cfg "github.com/twirapp/twir/libs/config"
)

var urlRegex = regexp.MustCompile(`https?://[a-zA-Z0-9\-]+\.[a-zA-Z0-9\-]+\.pinggy\.link`)

type Manager interface {
	Start(ctx context.Context) (string, error)
	Stop(ctx context.Context) error
}

type PinggyManager struct {
	config cfg.Config
	logger *slog.Logger
	cmd    *exec.Cmd
}

var _ Manager = (*PinggyManager)(nil)

func NewPinggyManager(config cfg.Config, logger *slog.Logger) *PinggyManager {
	return &PinggyManager{
		config: config,
		logger: logger,
	}
}

func (m *PinggyManager) Start(ctx context.Context) (string, error) {
	port := m.config.EventsubHttpPort
	if port == 0 {
		port = 3030
	}

	sshHost := m.config.TunnelSshHost
	if sshHost == "" {
		sshHost = "a.pinggy.io"
	}

	sshPort := m.config.TunnelSshPort
	if sshPort == 0 {
		sshPort = 443
	}

	remoteForward := fmt.Sprintf("-R0:localhost:%d", port)

	m.logger.InfoContext(
		ctx,
		"Starting pinggy.io tunnel",
		slog.String("ssh_host", sshHost),
		slog.Int("ssh_port", sshPort),
		slog.Int("local_port", port),
	)

	m.cmd = exec.CommandContext(
		ctx,
		"ssh",
		"-o", "StrictHostKeyChecking=accept-new",
		"-o", "ServerAliveInterval=30",
		"-o", "ServerAliveCountMax=3",
		"-p", fmt.Sprintf("%d", sshPort),
		remoteForward,
		sshHost,
	)

	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("tunnel: create stdout pipe: %w", err)
	}

	stderr, err := m.cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("tunnel: create stderr pipe: %w", err)
	}

	if err := m.cmd.Start(); err != nil {
		return "", fmt.Errorf("tunnel: start ssh process: %w", err)
	}

	combined := make(chan string, 100)
	go scanLines(bufio.NewReader(stdout), combined)
	go scanLines(bufio.NewReader(stderr), combined)

	timeout := time.NewTimer(30 * time.Second)
	defer timeout.Stop()

	for {
		select {
		case line, ok := <-combined:
			if !ok {
				return "", fmt.Errorf("tunnel: ssh process closed output before URL was received")
			}

			m.logger.DebugContext(ctx, "tunnel output", slog.String("line", line))

			if url := extractURL(line); url != "" {
				m.logger.InfoContext(
					ctx,
					"Tunnel established",
					slog.String("public_url", url),
				)
				return url, nil
			}

		case <-timeout.C:
			_ = m.Stop(ctx)
			return "", fmt.Errorf("tunnel: timed out waiting for public URL")
		}
	}
}

func (m *PinggyManager) Stop(_ context.Context) error {
	if m.cmd == nil || m.cmd.Process == nil {
		return nil
	}

	if err := m.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("tunnel: kill ssh process: %w", err)
	}

	return nil
}

func scanLines(r *bufio.Reader, out chan<- string) {
	defer close(out)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		out <- strings.TrimSpace(line)
	}
}

func extractURL(line string) string {
	matches := urlRegex.FindStringSubmatch(line)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}
