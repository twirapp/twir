package main

import (
	"github.com/twirapp/twir/apps/giveaways/app"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.FxOnlyErrorsLoggerOption(),
		app.App,
	).Run()
}
