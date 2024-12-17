package gql

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ravilushqa/otelgqlgen"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	apq_cache "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/apq-cache"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_groups"
	"github.com/twirapp/twir/libs/cache/twitch"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type Gql struct {
	*handler.Server
}

type Opts struct {
	fx.In

	Resolver              *resolvers.Resolver
	Directives            *directives.Directives
	Config                config.Config
	ApqCache              *apq_cache.APQCache
	Tracer                trace.Tracer
	CachedTwitchClient    *twitch.CachedTwitchClient
	Server                *server.Server
	CommandsGroupsService *commands_groups.Service
}

func New(opts Opts) *Gql {
	graphConfig := graph.Config{
		Resolvers: opts.Resolver,
	}
	graphConfig.Directives.IsAuthenticated = opts.Directives.IsAuthenticated
	graphConfig.Directives.HasAccessToSelectedDashboard = opts.Directives.HasAccessToSelectedDashboard
	graphConfig.Directives.IsAdmin = opts.Directives.IsAdmin
	graphConfig.Directives.HasChannelRolesDashboardPermission = opts.Directives.HasChannelRolesDashboardPermission
	graphConfig.Directives.Validate = opts.Directives.Validate

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
			InitFunc: auth.WsGqlInitFunc,
		},
	)

	srv.Use(otelgqlgen.Middleware())

	if opts.Config.AppEnv != "production" {
		srv.Use(extension.Introspection{})
	}

	playgroundHandler := playground.Handler("GraphQL", "/api/query")
	opts.Server.Any(
		"/", func(c *gin.Context) {
			playgroundHandler.ServeHTTP(c.Writer, c.Request)
		},
	)

	opts.Server.Any(
		"/query",
		func(c *gin.Context) {
			loader := data_loader.New(
				data_loader.Opts{
					CachedTwitchClient:    opts.CachedTwitchClient,
					CommandsGroupsService: opts.CommandsGroupsService,
				},
			)

			c.Request = c.Request.WithContext(
				context.WithValue(c.Request.Context(), data_loader.LoadersKey, loader),
			)

			c.Next()
		},
		func(c *gin.Context) {
			srv.ServeHTTP(c.Writer, c.Request)
		},
	)

	return &Gql{srv}
}
