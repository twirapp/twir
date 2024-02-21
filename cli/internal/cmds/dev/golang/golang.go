package golang

import (
	"context"

	"github.com/pterm/pterm"
)

var appsForStart = []twirApp{
	// {Name: "tokens", FxModule: tokens.App},
	// {Name: "events", FxModule: events.App},
	// {Name: "emotes-cacher", FxModule: emotescacher.App},
	// {Name: "scheduler", FxModule: scheduler.App},
	{name: "api"},
	{name: "tokens"},
	{name: "events"},
	{name: "emotes-cacher"},
	{name: "parser"},
	{name: "eventsub"},
	{name: "bots"},
	{name: "timers"},
	{name: "websockets"},
	{name: "ytsr"},
	{name: "scheduler"},
	{name: "discord"},
	// {Name: "bots", FxModule: bots.App},
	// {Name: "discord", FxModule: discord.App},
	// {Name: "timers", FxModule: timers.App},
	// {Name: "websockets", FxModule: websockets.App},
	// {Name: "ytsr", FxModule: ytsr.App},
}

type GoApps struct {
	apps []*twirApp
}

func New() (*GoApps, error) {
	ga := &GoApps{}
	for _, app := range appsForStart {
		application, err := newApplication(app.name)
		if err != nil {
			return nil, err
		}

		ga.apps = append(ga.apps, application)
	}

	return ga, nil
}

func (c *GoApps) Start(ctx context.Context) error {
	for _, app := range c.apps {
		app := app

		for i := 0; i < 3; i++ {
			pterm.Info.Println("Starting " + app.name)
			if err := app.start(); err != nil {
				return err
			}

			break
		}

		go func() {
			chann, err := app.watcher.Start(app.path)
			if err != nil {
				pterm.Fatal.Println(err)
			}

			for range chann {
				pterm.Info.Println("ReStarting " + app.name)
				if err := app.start(); err != nil {
					pterm.Error.Println(err)
				}
			}
		}()
	}

	return nil
}

func (c *GoApps) Stop() {
	for _, app := range c.apps {
		app.watcher.Stop()
		app.stop()
	}
}
