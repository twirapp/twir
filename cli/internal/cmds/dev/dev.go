package dev

import (
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/pterm/pterm"
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/cli/internal/cmds/build"
	"github.com/twirapp/twir/cli/internal/cmds/dependencies"
	"github.com/twirapp/twir/cli/internal/cmds/dev/frontend"
	"github.com/twirapp/twir/cli/internal/cmds/dev/golang"
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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "skip-deps",
				Value: false,
				Usage: "skip dependencies installation",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "run backend in debug mode",
			},
			&cli.BoolFlag{
				Name:  "proxy",
				Value: false,
				Usage: "start with proxy",
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
				return errors.New("NGROK_AUTH_TOKEN is required in .env. Please set it to enable ngrok")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
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

			if c.Bool("proxy") {
				go func() {
					if err := proxy.Cmd.Run(c); err != nil {
						pterm.Error.Println(err)
						return
					}
				}()
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
