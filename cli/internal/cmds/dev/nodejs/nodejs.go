package nodejs

import (
	"errors"
	"sync"
)

var appsForStart = []twirApp{
	{name: "integrations"},
	{name: "executron"},
	{name: "ytsub"},
}

type NodejsApps struct {
	apps []*twirApp
}

func New() (*NodejsApps, error) {
	fa := &NodejsApps{}
	for _, app := range appsForStart {
		application, err := newApplication(app.name)
		if err != nil {
			return nil, err
		}

		fa.apps = append(fa.apps, application)
	}

	return fa, nil
}

func (fa *NodejsApps) Start() error {
	for _, app := range fa.apps {
		if err := app.start(); err != nil {
			return err
		}
	}

	return nil
}

func (fa *NodejsApps) Stop() error {
	var wg sync.WaitGroup
	stopErrors := make([]error, len(fa.apps))
	for i, app := range fa.apps {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stopErrors[i] = app.stop()
		}()
	}
	wg.Wait()

	return errors.Join(stopErrors...)
}
