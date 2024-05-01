package build

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/gqlgen/api"
	gqlconfig "github.com/99designs/gqlgen/codegen/config"
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
		if err := LibsCmd.Run(c); err != nil {
			return err
		}

		if err := GqlCmd.Run(c); err != nil {
			return err
		}

		return build(`turbo run build --filter=!./apps/dota`, true)
	},
	Subcommands: []*cli.Command{
		LibsCmd,
		GqlCmd,
		AppBuildCmd,
	},
}

var LibsCmd = &cli.Command{
	Name: "libs",
	Action: func(context *cli.Context) error {
		if err := build(`turbo run build --filter=./libs/*`, false); err != nil {
			return err
		}

		pterm.Success.Println("Builded libs")

		return nil
	},
}

var AppBuildCmd = &cli.Command{
	Name:      "app",
	Args:      true,
	ArgsUsage: "api",
	Action: func(context *cli.Context) error {
		argument := context.Args().First()

		var golangApp *goapp.TwirGoApp
		for _, a := range goapp.Apps {
			if a.Name != argument {
				continue
			}

			foundApp, err := goapp.NewApplication(a.Name, false)
			if err != nil {
				return err
			}
			golangApp = foundApp
		}

		if golangApp != nil {
			if err := golangApp.Build(); err != nil {
				return err
			}

			pterm.Success.Printfln("Builded %s", golangApp.Name)
			return nil
		}

		return build(fmt.Sprintf(`turbo run build --filter=@twir/%s`, argument), false)
	},
}

var GqlCmd = &cli.Command{
	Name: "gql",
	Action: func(context *cli.Context) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		configDir := filepath.Join(cwd, "apps", "api-gql", "gqlgen.yml")
		cfg, err := gqlconfig.LoadConfig(configDir)
		if err != nil {
			return err
		}

		if err := api.Generate(cfg); err != nil {
			return err
		}

		pterm.Success.Println("Generated gqlgen")

		return nil
	},
}

func build(cmd string, withGoApps bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	pterm.Info.Println("Building twir")

	startTime := time.Now()

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

	if withGoApps {
		for _, app := range goapp.Apps {
			pterm.Info.Printfln("Building %s", app.Name)

			a, err := goapp.NewApplication(app.Name, false)
			if err != nil {
				pterm.Fatal.Println(err)
			}

			if err := a.Build(); err != nil {
				pterm.Fatal.Println(err)
			}
		}
	}

	if time.Since(startTime).Milliseconds() < 1000 {
		pterm.Success.Println(rainbow(">>> FULL TWIR TURBO ") + "🤙 🤙 🤙")
	} else {
		pterm.Success.Println("Builded")
	}

	return nil
}
