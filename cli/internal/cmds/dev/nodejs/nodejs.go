package nodejs

var appsForStart = []twirApp{
	{name: "integrations"},
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
	for _, app := range fa.apps {
		if err := app.stop(); err != nil {
			return err
		}
	}

	return nil
}
