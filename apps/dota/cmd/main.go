package main

import (
	"github.com/twirapp/twir/apps/dota/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
	).Run()
}
