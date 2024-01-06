package dev

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pterm/pterm"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/twirapp/twir/cli/internal/cmds/build"
	"github.com/twirapp/twir/cli/internal/cmds/dependencies"
	"github.com/twirapp/twir/cli/internal/cmds/dev/goapps"
	"github.com/twirapp/twir/cli/internal/cmds/migrations"
	"github.com/urfave/cli/v2"
)

type app struct {
	Name     string
	Stack    string
	Port     int
	SkipWait bool
}

var apps = []app{
	{Stack: "go", Name: "parser", Port: constants.PARSER_SERVER_PORT},
	{Stack: "go", Name: "eventsub", Port: constants.EVENTSUB_SERVER_PORT},
	{Stack: "node", Name: "eval", Port: constants.EVAL_SERVER_PORT},
	{Stack: "node", Name: "integrations", Port: constants.INTEGRATIONS_SERVER_PORT},
	{Stack: "frontend", Name: "dashboard", Port: 3006},
	{Stack: "frontend", Name: "landing", Port: 3005},
	{Stack: "frontend", Name: "overlays", Port: 3008},
	{Stack: "frontend", Name: "public-page", Port: 3007},
}

func CreateDevCommand() *cli.Command {
	defaultApps := make([]string, 0, len(apps))
	for _, a := range apps {
		defaultApps = append(defaultApps, a.Name)
	}

	var cmd = &cli.Command{
		Name:  "dev",
		Usage: "start project in dev mode",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "app",
				Usage:   "multiple app flags",
				Value:   cli.NewStringSlice(defaultApps...),
				Aliases: []string{"a"},
			},
		},
		Action: func(c *cli.Context) error {
			if err := dependencies.Cmd.Run(c); err != nil {
				return err
			}

			if err := build.LibsCmd.Run(c); err != nil {
				return err
			}

			if err := migrations.MigrateCmd.Run(c); err != nil {
				return err
			}

			golangApps, err := goapps.New()
			if err != nil {
				return err
			}
			defer golangApps.Stop(c.Context)

			if err := golangApps.Start(c.Context); err != nil {
				pterm.Error.Println(err)
				return err
			}

			exitSignal := make(chan os.Signal, 1)
			signal.Notify(exitSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			<-exitSignal

			return nil
		},
	}

	return cmd
}
