package gql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/resolvers"
	"go.uber.org/fx"
)

type Gql struct {
	*handler.Server
}

type Opts struct {
	fx.In

	Resolver   *resolvers.Resolver
	Directives *directives.Directives
	Config     config.Config
}

func New(opts Opts) *Gql {
	graphConfig := graph.Config{
		Resolvers: opts.Resolver,
	}
	graphConfig.Directives.IsAuthenticated = opts.Directives.IsAuthenticated
	graphConfig.Directives.HasAccessToSelectedDashboard = opts.Directives.HasAccessToSelectedDashboard
	graphConfig.Directives.IsAdmin = opts.Directives.IsAdmin

	schema := graph.NewExecutableSchema(graphConfig)

	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(
		transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	)

	if opts.Config.AppEnv != "production" {
		srv.Use(extension.Introspection{})
	}

	srv.Use(extension.FixedComplexityLimit(5))

	return &Gql{srv}
}
