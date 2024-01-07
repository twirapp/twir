package main

import (
	"github.com/satont/twir/apps/timers/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
	).Run()
}
