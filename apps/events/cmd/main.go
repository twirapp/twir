package main

import (
	"github.com/twirapp/twir/apps/events/app"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.FxDiOnlyErrors,
		app.App,
	).Run()
}
