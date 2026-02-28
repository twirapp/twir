package main

import (
	"github.com/twirapp/twir/apps/bots/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(

		app.App,
	).Run()
}
