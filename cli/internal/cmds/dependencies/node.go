package dependencies

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/goccy/go-json"
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
			Command: "bun install",
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

type packageJson struct {
	Engines struct {
		Node string `json:"node"`
		Pnpm string `json:"pnpm"`
		Bun  string `json:"bun"`
	} `json:"engines"`
}

func checkBunVersion() error {
	data, err := exec.Command("bun", "--version").Output()
	if err != nil {
		return err
	}

	installedVersion := string(data)
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

	rootPackageJsonPath := filepath.Join(wd, "package.json")
	packageJsonBytes, err := os.ReadFile(rootPackageJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var packageJsonData packageJson
	err = json.Unmarshal(packageJsonBytes, &packageJsonData)
	if err != nil {
		return fmt.Errorf("failed to parse package.json: %w", err)
	}

	contraint, err := semver.NewConstraint(packageJsonData.Engines.Bun)
	if err != nil {
		return fmt.Errorf("failed to parse bin constraint: %w", err)
	}

	if !contraint.Check(parsedInstalledVersion) {
		return fmt.Errorf(
			"installed bun version %s does not satisfy constraint %s",
			installedVersion,
			contraint,
		)
	}

	return nil
}
