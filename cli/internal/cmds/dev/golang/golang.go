package golang

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/goapp"
	"github.com/twirapp/twir/cli/internal/modutil"
	"github.com/twirapp/twir/cli/internal/watcher"
)

type GoApps struct {
	apps         []*goapp.TwirGoApp
	debugEnabled bool
	libWatcher   *watcher.Watcher
	resolver     *modutil.LibDependencyResolver
}

func New(enableDebug bool) (*GoApps, error) {
	ga := &GoApps{
		debugEnabled: enableDebug,
	}

	resolver, err := modutil.NewLibDependencyResolver()
	if err != nil {
		return nil, err
	}

	for _, app := range goapp.Apps {
		application, err := goapp.NewApplication(
			app.Name,
			enableDebug,
			app.Port,
			app.DebugPort,
			app.OnPortReady,
		)
		if err != nil {
			return nil, err
		}

		ga.apps = append(ga.apps, application)

		if err := resolver.ResolveApp(app.Name); err != nil {
			pterm.Warning.Printf("Failed to resolve deps for %s: %v\n", app.Name, err)
		}
	}

	ga.resolver = resolver
	ga.libWatcher = watcher.New()

	return ga, nil
}

func (c *GoApps) Start(ctx context.Context) error {
	for _, app := range c.apps {
		app := app

		pterm.Info.Println("Starting " + app.Name)
		if err := app.Start(); err != nil {
			return err
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

	if err := c.startLibWatcher(); err != nil {
		pterm.Warning.Printf("Failed to start lib watcher: %v\n", err)
	}

	return nil
}

func (c *GoApps) startLibWatcher() error {
	paths := c.resolver.GetAllWatchedPaths()
	if len(paths) == 0 {
		return nil
	}

	chann, err := c.libWatcher.StartWithPaths(paths)
	if err != nil {
		return err
	}

	go func() {
		for eventPath := range chann {
			if eventPath == "" {
				continue
			}

			if strings.HasSuffix(eventPath, "~") || strings.Contains(eventPath, ".out") {
				continue
			}

			apps := c.resolver.GetAppsForFile(eventPath)
			if len(apps) == 0 {
				continue
			}

			for _, appName := range apps {
				for _, app := range c.apps {
					if app.Name != appName {
						continue
					}

					pterm.Info.Printf("ReStarting %s (lib change: %s)\n", app.Name, filepath.Base(eventPath))
					if err := app.Start(); err != nil {
						pterm.Error.Println(err)
					}
					break
				}
			}
		}
	}()

	return nil
}

func (c *GoApps) Stop() error {
	if c.libWatcher != nil {
		c.libWatcher.Stop()
	}

	for _, app := range c.apps {
		app.Watcher.Stop()
		if err := app.Stop(); err != nil {
			return err
		}
	}

	return nil
}
