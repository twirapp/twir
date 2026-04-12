package main

import (
	"github.com/twirapp/twir/apps/twitch-mock/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.Module,
	).Run()
}
