package main

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/gql"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/resolvers"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	"github.com/twirapp/twir/apps/api-gql/internal/redis"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			cfg.NewFx,
			redis.New,
			sessions.New,
			resolvers.New,
			directives.New,
			gql.New,
			httpserver.New,
		),
		fx.Invoke(
			httpserver.New,
		),
	).Run()
}
