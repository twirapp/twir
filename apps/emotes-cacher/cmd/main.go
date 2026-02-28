package main

import (
	"github.com/twirapp/twir/apps/emotes-cacher/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.App,
	).Run()
}
