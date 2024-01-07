package frontendapps

var appsForStart = []twirApp{
	{name: "dashboard"},
	{name: "landing"},
	{name: "overlays"},
	{name: "public-page"},
}

type FrontendApps struct {
	apps []*twirApp
}

func New() (*FrontendApps, error) {
	fa := &FrontendApps{}
	for _, app := range appsForStart {
		application, err := newApplication(app.name)
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
