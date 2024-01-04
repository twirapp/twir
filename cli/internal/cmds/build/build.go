package build

import (
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "build",
	Usage: "build application",
	Action: func(c *cli.Context) error {
		return build("turbo run build --filter=!./apps/dota")
	},
	Subcommands: []*cli.Command{
		LibsCmd,
	},
}

var LibsCmd = &cli.Command{
	Name: "libs",
	Action: func(context *cli.Context) error {
		return build("turbo run build --filter='./libs/*'")
	},
}

func build(cmd string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	pterm.Info.Println("Building twir")
	spinner, _ := pterm.DefaultSpinner.Start("Building...")

	startTime := time.Now()

	err = shell.ExecCommand(
		shell.ExecCommandOpts{
			Command: cmd,
			Pwd:     wd,
		},
	)
	if err != nil {
		spinner.Fail(err)
		return err
	}

	if time.Since(startTime).Minutes() <= 1 {
		spinner.Success("Full turbo 🤙 🤙 🤙")
	} else {
		spinner.Success("Builded")
	}

	return nil
}
