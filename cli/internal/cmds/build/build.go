package build

import (
	"fmt"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/goapp"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:    "build",
	Usage:   "build application",
	Aliases: []string{"b"},
	Action: func(c *cli.Context) error {
		return build(`turbo run build --filter=!./apps/dota`)
	},
	Subcommands: []*cli.Command{
		LibsCmd,
	},
}

var LibsCmd = &cli.Command{
	Name: "libs",
	Action: func(context *cli.Context) error {
		return build(`turbo run build --filter=./libs/*`)
	},
}

func build(cmd string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	pterm.Info.Println("Building twir")

	startTime := time.Now()
	fmt.Println(wd)

	err = shell.ExecCommand(
		shell.ExecCommandOpts{
			Command: cmd,
			Pwd:     wd,
			Stderr:  os.Stderr,
			Stdout:  os.Stdout,
		},
	)
	if err != nil {
		pterm.Fatal.Println(err)
		return err
	}

	for _, app := range goapp.Apps {
		pterm.Info.Printfln("Building %s", app.Name)

		a, err := goapp.NewApplication(app.Name)
		if err != nil {
			pterm.Fatal.Println(err)
		}

		if err := a.Build(); err != nil {
			pterm.Fatal.Println(err)
		}
	}

	if time.Since(startTime).Milliseconds() < 1000 {
		pterm.Success.Println(rainbow(">>> FULL TWIR TURBO ") + "🤙 🤙 🤙")
	} else {
		pterm.Success.Println("Builded")
	}

	return nil
}
