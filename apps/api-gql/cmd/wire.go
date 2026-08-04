//go:build wireinject

package main

import "github.com/goforj/wire"

func initializeApplication() (*Application, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
