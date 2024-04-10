package main

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/gqlhandler"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	"github.com/twirapp/twir/apps/api-gql/internal/redis"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"github.com/twirapp/twir/apps/api-gql/resolvers"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			cfg.NewFx,
			redis.New,
			sessions.New,
			resolvers.New,
			gqlhandler.New,
			httpserver.New,
		),
		fx.Invoke(
			httpserver.New,
		),
	).Run()
}
