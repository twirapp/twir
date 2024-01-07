package dependencies

import (
	"os"

	"github.com/twirapp/twir/cli/internal/shell"
)

func installNodeDeps() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = shell.ExecCommand(
		shell.ExecCommandOpts{
			Command: "pnpm install --frozen-lockfile",
			Pwd:     wd,
			Stderr:  os.Stderr,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
