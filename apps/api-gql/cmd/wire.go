//go:build wireinject

package main

import (
	"github.com/goforj/wire"
	"github.com/twirapp/twir/apps/api-gql/internal/di"
)

func initializeApplication() (*di.Application, error) {
	wire.Build(di.ProviderSet)
	return nil, nil
}
