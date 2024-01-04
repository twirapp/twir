package dev

import (
	"fmt"

	"github.com/satont/twir/libs/grpc/constants"
	"github.com/twirapp/twir/cli/internal/cmds/build"
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
	{Stack: "go", Name: "tokens", Port: constants.TOKENS_SERVER_PORT},
	{Stack: "go", Name: "events", Port: constants.EVENTS_SERVER_PORT},
	{Stack: "go", Name: "emotes-cacher", Port: constants.EMOTES_CACHER_SERVER_PORT},
	{Stack: "go", Name: "parser", Port: constants.PARSER_SERVER_PORT},
	{Stack: "go", Name: "eventsub", Port: constants.EVENTSUB_SERVER_PORT},
	{Stack: "node", Name: "eval", Port: constants.EVAL_SERVER_PORT},
	{Stack: "go", Name: "bots", Port: constants.BOTS_SERVER_PORT},
	{Stack: "go", Name: "timers", Port: constants.TIMERS_SERVER_PORT},
	{Stack: "go", Name: "websockets", Port: constants.WEBSOCKET_SERVER_PORT},
	{Stack: "go", Name: "ytsr", Port: constants.YTSR_SERVER_PORT},
	{Stack: "node", Name: "integrations", Port: constants.INTEGRATIONS_SERVER_PORT},
	{Stack: "go", Name: "api", Port: 3002},
	{Stack: "go", Name: "scheduler", Port: constants.SCHEDULER_SERVER_PORT},
	{Stack: "go", Name: "discord", Port: constants.DISCORD_SERVER_PORT, SkipWait: true},
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
		Name: "dev",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "app",
				Usage:   "multiple app flags",
				Value:   cli.NewStringSlice(defaultApps...),
				Aliases: []string{"a"},
			},
		},
		Action: func(c *cli.Context) error {
			if err := build.LibsCmd.Run(c); err != nil {
				return err
			}

			if err := migrations.MigrateCmd.Run(c); err != nil {
				return err
			}

			fmt.Println(c.StringSlice("app"))
			for _, app := range c.StringSlice("app") {
				fmt.Println(app)
			}

			return nil
		},
	}

	return cmd
}
