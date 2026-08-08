//go:build wireinject

package main

import (
	"github.com/goforj/wire"
	"github.com/twirapp/twir/apps/timers/app"
)

func initializeApplication() (*app.Application, error) {
	wire.Build(app.ProviderSet)
	return nil, nil
}
