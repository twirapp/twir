package goapps

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/twirapp/twir/cli/internal/watcher"
)

type application struct {
	name    string
	cmd     *exec.Cmd
	path    string
	watcher *watcher.Watcher
}

var apps = []application{
	// {Name: "tokens", FxModule: tokens.App},
	// {Name: "events", FxModule: events.App},
	// {Name: "emotes-cacher", FxModule: emotescacher.App},
	// {Name: "scheduler", FxModule: scheduler.App},
	{name: "api"},
	// {Name: "bots", FxModule: bots.App},
	// {Name: "discord", FxModule: discord.App},
	// {Name: "timers", FxModule: timers.App},
	// {Name: "websockets", FxModule: websockets.App},
	// {Name: "ytsr", FxModule: ytsr.App},
}

type GoApps struct {
	apps []*application
}

func New() (*GoApps, error) {
	ga := &GoApps{}
	for _, app := range apps {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		appPath := filepath.Join(wd, "apps", app.name)

		cmd, err := shell.CreateCommand(
			shell.ExecCommandOpts{
				Command: "go run ./cmd/main.go",
				Pwd:     appPath,
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
			},
		)
		if err != nil {
			return nil, err
		}

		app.cmd = cmd
		app.path = appPath
		app.watcher = watcher.New()

		ga.apps = append(ga.apps, &app)
	}

	return ga, nil
}

func (c *GoApps) Start(ctx context.Context) {
	for _, app := range c.apps {
		app := app
		pterm.Info.Println("Starting " + app.name)

		if err := app.cmd.Start(); err != nil {
			pterm.Fatal.Println(err)
		}

		go func() {
			chann, err := app.watcher.Start(app.path)
			if err != nil {
				pterm.Fatal.Println(err)
			}

			for range chann {
				if err := app.cmd.Process.Kill(); err != nil {
					pterm.Fatal.Println(err)
				}
				if err := app.cmd.Wait(); err != nil && err.Error() != "signal: killed" {
					pterm.Fatal.Println(err)
				}
				time.Sleep(5 * time.Second)
				if err := app.cmd.Start(); err != nil {
					pterm.Fatal.Println(err)
				}
			}
		}()
	}
}

func (c *GoApps) Stop(ctx context.Context) {
	for _, app := range c.apps {
		app.watcher.Stop()
		app.cmd.Process.Kill()
	}
}
