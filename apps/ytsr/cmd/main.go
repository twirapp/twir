package main

import (
	"github.com/satont/twir/apps/ytsr/app"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.FxDiOnlyErrors,
		app.App,
	).Run()
}
