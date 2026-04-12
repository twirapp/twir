package helpers

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/pterm/pterm"
)

// EnsureDockerComposeRunning checks if docker compose is running and starts it if not
func EnsureDockerComposeRunning(ctx context.Context) error {
	composeFile := "docker-compose.dev.yml"

	pterm.Info.Println("Checking docker compose status...")

	// Try to start docker compose with up -d
	// This command is idempotent - it will start only stopped containers
	pterm.Info.Println("Starting docker compose services...")

	upCmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "up", "-d", "--build")
	output, err := upCmd.CombinedOutput()

	if err != nil {
		pterm.Error.Printfln("Failed to start docker compose: %v", err)
		pterm.Error.Printfln("Output: %s", string(output))
		return fmt.Errorf("failed to start docker compose: %w", err)
	}

	pterm.Success.Println("Docker compose is running")
	return nil
}
