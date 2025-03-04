package golang

import (
	"context"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/goapp"
)

type GoApps struct {
	apps         []*goapp.TwirGoApp
	debugEnabled bool
}

func New(enableDebug bool) (*GoApps, error) {
	ga := &GoApps{
		debugEnabled: enableDebug,
	}
	for _, app := range goapp.Apps {
		application, err := goapp.NewApplication(app.Name, enableDebug, app.Port, app.OnPortReady)
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
			pterm.Info.Println("Starting " + app.Name)
			if err := app.Start(); err != nil {
				return err
			}

			break
		}

		go func() {
			chann, err := app.Watcher.Start(app.Path)
			if err != nil {
				pterm.Fatal.Println(err)
			}

			for range chann {
				pterm.Info.Println("ReStarting " + app.Name)
				if err := app.Start(); err != nil {
					pterm.Error.Println(err)
				}
			}
		}()
	}

	return nil
}

func (c *GoApps) Stop() {
	for _, app := range c.apps {
		app.Watcher.Stop()
		app.Stop()
	}
}
