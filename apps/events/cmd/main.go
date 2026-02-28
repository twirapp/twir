package main

import (
	"github.com/twirapp/twir/apps/events/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
	).Run()
}
