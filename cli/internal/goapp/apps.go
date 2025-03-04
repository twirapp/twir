package goapp

import (
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/twirapp/twir/cli/internal/shell"
)

var Apps = []TwirGoApp{
	{Name: "api"},
	{Name: "tokens"},
	{Name: "events"},
	{Name: "emotes-cacher"},
	{Name: "parser"},
	{Name: "eventsub"},
	{Name: "bots"},
	{Name: "timers"},
	{Name: "websockets"},
	{Name: "ytsr"},
	{Name: "scheduler"},
	{Name: "discord"},
	{
		Name: "api-gql",
		Port: lo.ToPtr(3009),
		OnPortReady: func() {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			pwd := filepath.Join(wd, "libs", "api")

			cmd, err := shell.CreateCommand(
				shell.ExecCommandOpts{
					Command: "bun run build:openapi",
					Pwd:     pwd,
				},
			)
			if err != nil {
				panic(err)
			}

			cmd.Run()
		},
	},
}
