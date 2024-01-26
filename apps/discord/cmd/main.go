package main

import (
	"github.com/satont/twir/apps/discord/app"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
		logger.FxDiOnlyErrors,
	).Run()
}
