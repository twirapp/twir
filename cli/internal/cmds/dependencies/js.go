package dependencies

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/twirapp/twir/cli/internal/shell"
)

func installNodeDeps() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = checkBunVersion()
	if err != nil {
		return err
	}

	err = shell.ExecCommand(
		shell.ExecCommandOpts{
			Command: "bun install --frozen-lockfile",
			Pwd:     wd,
			Stderr:  os.Stderr,
			Stdout:  os.Stdout,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func checkBunVersion() error {
	currentVersionData, err := exec.Command("bun", "--version").Output()
	if err != nil {
		return err
	}

	installedVersion := string(currentVersionData)
	installedVersion = strings.Replace(installedVersion, "\n", "", -1)
	installedVersion = strings.Replace(installedVersion, "\r", "", -1)

	parsedInstalledVersion, err := semver.NewVersion(installedVersion)
	if err != nil {
		return fmt.Errorf("failed to parse installed pnpm version: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	bunVersionFileContent, err := os.ReadFile(filepath.Join(wd, ".bun-version"))
	if err != nil {
		return fmt.Errorf("failed to read .bun-version file: %w", err)
	}

	constraint, err := semver.NewConstraint(fmt.Sprintf(">=%s", string(bunVersionFileContent)))
	if err != nil {
		return fmt.Errorf("failed to parse bin constraint: %w", err)
	}

	if !constraint.Check(parsedInstalledVersion) {
		return fmt.Errorf(
			"installed bun version %s does not satisfy constraint %s",
			installedVersion,
			constraint,
		)
	}

	return nil
}
