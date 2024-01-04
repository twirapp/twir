package dependencies

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
)

func installNodeDeps() error {
	spinner, _ := pterm.DefaultSpinner.Start("Install node deps...")

	wd, err := os.Getwd()
	if err != nil {
		spinner.Info(err)
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
		spinner.Fail(err)
		return err
	}

	spinner.Success("Nodejs deps installed")

	return nil
}
