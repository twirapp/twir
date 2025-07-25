package main

import (
	"github.com/twirapp/twir/apps/api/app"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.FxDiOnlyErrors,
		app.App,
	).Run()
}
