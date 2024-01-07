package dev

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/cmds/build"
	"github.com/twirapp/twir/cli/internal/cmds/dependencies"
	"github.com/twirapp/twir/cli/internal/cmds/dev/frontend"
	"github.com/twirapp/twir/cli/internal/cmds/dev/golang"
	"github.com/twirapp/twir/cli/internal/cmds/dev/nodejs"
	"github.com/twirapp/twir/cli/internal/cmds/migrations"
	"github.com/urfave/cli/v2"
)

type app struct {
	Name     string
	Stack    string
	Port     int
	SkipWait bool
}

func CreateDevCommand() *cli.Command {
	var cmd = &cli.Command{
		Name:  "dev",
		Usage: "start project in dev mode",
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

			golangApps, err := golang.New()
			if err != nil {
				return err
			}

			frontendApps, err := frontend.New()
			if err != nil {
				return err
			}

			nodejsApps, err := nodejs.New()
			if err != nil {
				return err
			}

			if err := golangApps.Start(c.Context); err != nil {
				pterm.Error.Println(err)
				return err
			}

			if err := frontendApps.Start(); err != nil {
				pterm.Error.Println(err)
				return err
			}

			if err := nodejsApps.Start(); err != nil {
				pterm.Error.Println(err)
				return err
			}

			exitSignal := make(chan os.Signal, 1)
			signal.Notify(exitSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			<-exitSignal
			golangApps.Stop()
			frontendApps.Stop()
			nodejsApps.Stop()

			return nil
		},
	}

	return cmd
}
