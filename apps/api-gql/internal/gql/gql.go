package gql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
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
}

func New(opts Opts) *Gql {
	config := graph.Config{
		Resolvers: opts.Resolver,
	}
	config.Directives.IsAuthenticated = opts.Directives.IsAuthenticated
	config.Directives.HasAccessToSelectedDashboard = opts.Directives.HasAccessToSelectedDashboard
	config.Directives.IsAdmin = opts.Directives.IsAdmin

	schema := graph.NewExecutableSchema(config)

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
	srv.Use(extension.Introspection{})
	srv.Use(extension.FixedComplexityLimit(5))

	return &Gql{srv}
}
