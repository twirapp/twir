package main

import (
	"github.com/twirapp/twir/apps/timers/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
	).Run()
}
