package dependencies

import (
	"os"

	"github.com/twirapp/twir/cli/internal/shell"
)

func installGolangDeps() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return shell.ExecCommand(
		shell.ExecCommandOpts{
			Command: "go mod download",
			Pwd:     wd,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		},
	)
}
