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
	"github.com/twirapp/twir/cli/internal/cmds/dev/helpers"
	"github.com/twirapp/twir/cli/internal/cmds/dev/nodejs"
	"github.com/twirapp/twir/cli/internal/cmds/migrations"
	"github.com/twirapp/twir/cli/internal/cmds/proxy"
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
		Subcommands: []*cli.Command{
			helpers.CleanPortsCmd,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "skip-deps",
				Usage: "skip dependencies installation",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "run backend in debug mode",
			},
			&cli.BoolFlag{
				Name:  "skip-docker",
				Usage: "skip docker compose check and startup",
			},
			&cli.BoolFlag{
				Name:  "skip-proxy",
				Usage: "skip proxy start",
			},
		},
		Action: func(c *cli.Context) error {
			// Ensure docker compose is running before starting anything
			if !c.Bool("skip-docker") {
				if err := helpers.EnsureDockerComposeRunning(c.Context); err != nil {
					pterm.Fatal.Println("Failed to start docker compose:", err)
					return err
				}
			}

			if !c.Bool("skip-proxy") {
				proxyStartedChan, proxyInstance, err := proxy.StartProxy(false)
				if err != nil {
					pterm.Fatal.Println(err)
					return err
				}
				<-proxyStartedChan

				defer proxyInstance.Stop()
			}

			isDebugEnabled := c.Bool("debug")

			if !c.Bool("skip-deps") {
				if err := dependencies.Cmd.Run(c); err != nil {
					return err
				}
			}

			if err := build.LibsCmd.Run(c); err != nil {
				pterm.Fatal.Println(err)
				return err
			}

			if err := build.GqlCmd.Run(c); err != nil {
				pterm.Fatal.Println(err)
				return err
			}

			if err := migrations.MigrateCmd.Run(c); err != nil {
				pterm.Fatal.Println(err)
				return err
			}

			golangApps, err := golang.New(isDebugEnabled)
			if err != nil {
				pterm.Fatal.Println(err)
				return err
			}

			frontendApps, err := frontend.New()
			if err != nil {
				pterm.Fatal.Println(err)
				return err
			}

			nodejsApps, err := nodejs.New()
			if err != nil {
				pterm.Fatal.Println(err)
				return err
			}

			if err := golangApps.Start(c.Context); err != nil {
				pterm.Error.Println(err)
				return err
			}
			defer golangApps.Stop()

			if err := frontendApps.Start(); err != nil {
				pterm.Error.Println(err)
				return err
			}
			defer frontendApps.Stop()

			if err := nodejsApps.Start(); err != nil {
				pterm.Error.Println(err)
				return err
			}
			defer nodejsApps.Stop()

			exitSignal := make(chan os.Signal, 1)
			signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

			<-exitSignal

			return nil
		},
	}

	return cmd
}
