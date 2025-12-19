package main

import (
	"log"
	"os"

	"github.com/twirapp/twir/cli/internal/cmds/build"
	"github.com/twirapp/twir/cli/internal/cmds/dependencies"
	"github.com/twirapp/twir/cli/internal/cmds/dev"
	"github.com/twirapp/twir/cli/internal/cmds/execbin"
	"github.com/twirapp/twir/cli/internal/cmds/generate"
	"github.com/twirapp/twir/cli/internal/cmds/kill"
	"github.com/twirapp/twir/cli/internal/cmds/migrations"
	"github.com/twirapp/twir/cli/internal/cmds/proxy"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "go run cmd/main.go",
		Description: "TwirApp cli for helping in manage project",
		Commands: []*cli.Command{
			dependencies.Cmd,
			migrations.Cmd,
			proxy.Cmd,
			generate.Cmd,
			build.Cmd,
			dev.CreateDevCommand(),
			execbin.Cmd,
			kill.Cmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
