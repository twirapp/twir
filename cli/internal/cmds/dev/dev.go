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
		},
		Action: func(c *cli.Context) error {
			// Ensure docker compose is running before starting anything
			skipDocker := c.Bool("skip-docker")
			if !skipDocker {
				if err := helpers.EnsureDockerComposeRunning(c.Context); err != nil {
					pterm.Fatal.Println("Failed to start docker compose:", err)
					return err
				}
			}

			proxyStartedChan, proxyInstance, err := proxy.StartProxy(false)
			if err != nil {
				pterm.Fatal.Println(err)
				return err
			}
			<-proxyStartedChan

			skipDeps := c.Bool("skip-deps")
			isDebugEnabled := c.Bool("debug")

			if !skipDeps {
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

			if err := frontendApps.Start(); err != nil {
				pterm.Error.Println(err)
				return err
			}

			if err := nodejsApps.Start(); err != nil {
				pterm.Error.Println(err)
				return err
			}

			exitSignal := make(chan os.Signal, 1)
			signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

			<-exitSignal
			if err := golangApps.Stop(); err != nil {
				pterm.Error.Println("Failed to stop golang apps:", err)
			}
			if err := frontendApps.Stop(); err != nil {
				pterm.Error.Println("Failed to stop frontend apps:", err)
			}
			if err := nodejsApps.Stop(); err != nil {
				pterm.Error.Println("Failed to stop nodejs apps:", err)
			}
			if err := proxyInstance.Stop(); err != nil {
				pterm.Error.Println("Failed to stop proxy:", err)
			}

			return nil
		},
	}

	return cmd
}
