package main

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/gql"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/resolvers"
	subscriptions_store "github.com/twirapp/twir/apps/api-gql/internal/gql/subscriptions-store"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	"github.com/twirapp/twir/apps/api-gql/internal/minio"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"github.com/twirapp/twir/libs/baseapp"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp("api-gql"),
		fx.Provide(
			sessions.New,
			func(config cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
			minio.New,
			twitchcache.New,
			subscriptions_store.New,
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
