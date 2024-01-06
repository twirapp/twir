package main

import (
	"github.com/satont/twir/apps/events/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
	).Run()
}
