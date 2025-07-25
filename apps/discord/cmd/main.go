package main

import (
	"github.com/twirapp/twir/apps/discord/app"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
		logger.FxDiOnlyErrors,
	).Run()
}
