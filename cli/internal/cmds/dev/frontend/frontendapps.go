package frontend

import (
	"os"
	"path/filepath"
)

var appsForStart = []twirApp{
	{name: "dashboard"},
	{name: "web"},
	{name: "overlays"},
}

type FrontendApps struct {
	apps []*twirApp
}

func New() (*FrontendApps, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fa := &FrontendApps{}
	for _, app := range appsForStart {
		var path string
		if app.name == "web" {
			path = filepath.Join(wd, app.name)
		} else {
			path = filepath.Join(wd, "frontend", app.name)
		}

		application, err := newApplication(app.name, path)
		if err != nil {
			return nil, err
		}

		fa.apps = append(fa.apps, application)
	}

	return fa, nil
}

func (fa *FrontendApps) Start() error {
	for _, app := range fa.apps {
		if err := app.start(); err != nil {
			return err
		}
	}

	return nil
}

func (fa *FrontendApps) Stop() error {
	for _, app := range fa.apps {
		if err := app.stop(); err != nil {
			return err
		}
	}

	return nil
}
