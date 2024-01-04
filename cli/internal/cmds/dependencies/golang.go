package dependencies

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
)

var re = regexp.MustCompile(`use\s*\(\s*([\s\S]*?)\s*\)`)

func installGolangDeps() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	goWorkContent, err := os.ReadFile(filepath.Join(wd, "go.work"))
	if err != nil {
		return err
	}

	matches := re.FindStringSubmatch(string(goWorkContent))

	var paths []string

	// The extracted content is in the second submatch
	if len(matches) > 1 {
		content := matches[1]

		// Split the content into individual rows
		rows := strings.Split(content, "\n")

		// Print each row
		for _, row := range rows {
			row = strings.TrimSpace(row)
			if row != "" {
				paths = append(paths, row)
			}
		}
	}

	spinner, _ := pterm.DefaultSpinner.Start("Install golang deps...")

	for _, p := range paths {
		// name := strings.Split(p, "/")

		err = goInstallDepsForPath(p)
		if err != nil {
			pterm.Error.Println(err)
			return err
		}
	}

	spinner.Success("Golang deps installed")

	return nil
}

func goInstallDepsForPath(path string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return shell.ExecCommand(
		shell.ExecCommandOpts{
			Command: "go mod download",
			Pwd:     filepath.Join(wd, path),
		},
	)
}
