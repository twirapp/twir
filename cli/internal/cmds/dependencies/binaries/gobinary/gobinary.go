package gobinary

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type GoBinary struct {
	Url string
}

func (c GoBinary) GetNameAndVersionFromUrl() (string, string) {
	splittedBinaryName := strings.Split(c.Url, "/")
	binaryName := splittedBinaryName[len(splittedBinaryName)-1]
	name := strings.Split(binaryName, "@")[0]
	version := strings.Split(binaryName, "@")[1]

	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	return name, version
}

func (c GoBinary) IsInstalled() (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	name, _ := c.GetNameAndVersionFromUrl()

	if _, err := os.Stat(filepath.Join(wd, ".bin", name)); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

var ErrCannotGetVersion = errors.New("binary not found")

func (c GoBinary) GetGolangBinaryVersion(name string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	path := filepath.Join(wd, ".bin", name)

	cmd := exec.Command("go", "version", "-m", path)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cannot get %s version: %w", name, err)
	}

	output := string(cmdOutput)

	r := regexp.MustCompile(`\s(path|mod)\s+([^\s]+)(\s+([^\s]+))?`)

	var version string

	for _, l := range strings.Split(output, "\n") {
		matches := r.FindStringSubmatch(l)
		if len(matches) == 0 {
			continue
		}

		switch matches[1] {
		case "mod":
			version = matches[4]
		default:
			continue
		}
	}

	if version == "" {
		return "", ErrCannotGetVersion
	}

	// replace everything after "+"
	return strings.Split(version, "+")[0], nil
}

func (c GoBinary) Install() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "install", c.Url)

	cmd.Dir = wd
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOBIN="+filepath.Join(wd, ".bin"), "GOFLAGS=")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cannot Install %s: %w", c.Url, err)
	}

	return nil
}
