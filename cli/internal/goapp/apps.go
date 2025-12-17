package goapp

import (
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/twirapp/twir/cli/internal/shell"
)

var Apps = []TwirGoApp{
	{Name: "api", DebugPort: 2345},
	{Name: "tokens", DebugPort: 2346},
	{Name: "events", DebugPort: 2347},
	{Name: "emotes-cacher", DebugPort: 2348},
	{Name: "parser", DebugPort: 2349},
	{Name: "eventsub", DebugPort: 2350},
	{Name: "bots", DebugPort: 2351},
	{Name: "timers", DebugPort: 2352},
	{Name: "websockets", DebugPort: 2353},
	{Name: "scheduler", DebugPort: 2355},
	{Name: "discord", DebugPort: 2356},
	{
		Name:      "api-gql",
		Port:      lo.ToPtr(3009),
		DebugPort: 2359,
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
