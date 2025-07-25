package dev

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/pterm/pterm"
	cfg "github.com/twirapp/twir/libs/config"
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
		},
		Before: func(context *cli.Context) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			config, err := cfg.NewWithEnvPath(filepath.Join(wd, ".env"))
			if err != nil {
				return err
			}

			if config.NgrokAuthToken == "" {
				return errors.New(
					"NGROK_AUTH_TOKEN is required in .env. Please set it to enable ngrok",
				)
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			proxyStartedChan, err := proxy.StartProxy(false)
			if err != nil {
				pterm.Fatal.Println(err)
				return err
			}
			// wait proxy to up
			<-proxyStartedChan

			skipDeps := c.Bool("skip-deps")
			isDebugEnabled := c.Bool("debug")

			fmt.Println("isDebugEnabled:", isDebugEnabled)
			fmt.Println(c.Args())

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
