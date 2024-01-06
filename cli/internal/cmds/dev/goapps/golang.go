package goapps

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/rjeczalik/notify"
	"github.com/samber/lo"
	api "github.com/satont/twir/apps/api/app"
	cfg "github.com/satont/twir/libs/config"

	// bots "github.com/satont/twir/apps/bots/app"
	// discord "github.com/satont/twir/apps/discord/app"
	// emotescacher "github.com/satont/twir/apps/emotes-cacher/app"
	// events "github.com/satont/twir/apps/events/app"
	// scheduler "github.com/satont/twir/apps/scheduler/app"
	// timers "github.com/satont/twir/apps/timers/app"
	// tokens "github.com/satont/twir/apps/tokens/app"
	// websockets "github.com/satont/twir/apps/websockets/app"
	// ytsr "github.com/satont/twir/apps/ytsr/app"
	"go.uber.org/fx"
)

type twirApplication struct {
	Name     string
	FxModule fx.Option
}

var apps = []twirApplication{
	// {Name: "tokens", FxModule: tokens.App},
	// {Name: "events", FxModule: events.App},
	// {Name: "emotes-cacher", FxModule: emotescacher.App},
	// {Name: "scheduler", FxModule: scheduler.App},
	{Name: "api", FxModule: api.App},
	// {Name: "bots", FxModule: bots.App},
	// {Name: "discord", FxModule: discord.App},
	// {Name: "timers", FxModule: timers.App},
	// {Name: "websockets", FxModule: websockets.App},
	// {Name: "ytsr", FxModule: ytsr.App},
}

type applicationWithWatcher struct {
	*twirApplication
	fxApp *fx.App
}

type GoApps struct {
	apps []*applicationWithWatcher
}

func New() *GoApps {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ga := &GoApps{}

	envPath := filepath.Join(wd, ".env")
	for _, app := range apps {
		fxApp := fx.New(
			app.FxModule,
			fx.Decorate(
				func() cfg.Config {
					return cfg.NewFxWithPath(envPath)
				},
			),
			fx.NopLogger,
		)

		ga.apps = append(
			ga.apps,
			&applicationWithWatcher{
				twirApplication: &app,
				fxApp:           fxApp,
			},
		)
	}

	return ga
}

func (c *GoApps) Start(ctx context.Context) {
	for _, app := range c.apps {
		app := app
		pterm.Info.Println("Starting " + app.Name)
		go app.fxApp.Start(ctx) //nolint:errcheck

		go func() {
			pterm.Info.Println("Starting watcher for " + app.Name)
			watcher, err := c.watchAppFsUpdate(app.Name)
			if err != nil {
				panic(err)
			}

			for range watcher {
				if err := app.fxApp.Stop(ctx); err != nil {
					panic(err)
				}
				app.fxApp.Start(ctx) //nolint:errcheck
			}
		}()
	}
}

func (c *GoApps) Stop(ctx context.Context) {
	for _, app := range c.apps {
		app.fxApp.Stop(ctx) //nolint:errcheck
	}
}

func (c *GoApps) watchAppFsUpdate(appName string) (chan struct{}, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	watchPath := filepath.Join(wd, "apps", appName) + "..."

	notifyChan := make(chan notify.EventInfo, 1)
	chann := make(chan struct{}, 1)

	if err := notify.Watch(watchPath, notifyChan, notify.All); err != nil {
		return nil, err
	}

	reload, _ := lo.NewDebounce(
		1*time.Second,
		func() {
			chann <- struct{}{}
		},
	)

	go func() {
		for event := range notifyChan {
			if strings.HasSuffix(event.Path(), "~") {
				continue
			}

			reload()
		}
	}()

	return chann, nil
}
