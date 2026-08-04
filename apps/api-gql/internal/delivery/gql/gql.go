package gql

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ravilushqa/otelgqlgen"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlerrors"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_groups"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	twitchservice "github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	"github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel/trace"
)

type Gql struct {
	*handler.Server
}

func New(resolver *resolvers.Resolver, directivesService *directives.Directives, config config.Config, tracer trace.Tracer, cachedTwitchClient *twitch.CachedTwitchClient, httpServer *server.Server, commandsGroupsService *commands_groups.Service, commandsResponseService *commands_responses.Service, twitchService *twitchservice.Service, dataLoaderFactory *data_loader.LoaderFactory, middlewaresService *middlewares.Middlewares, loggerInstance *slog.Logger) *Gql {
	graphConfig := graph.Config{
		Resolvers: resolver,
	}
	graphConfig.Directives.IsAuthenticated = directivesService.IsAuthenticated
	graphConfig.Directives.HasAccessToSelectedDashboard = directivesService.HasAccessToSelectedDashboard
	graphConfig.Directives.IsAdmin = directivesService.IsAdmin
	graphConfig.Directives.HasChannelRolesDashboardPermission = directivesService.HasChannelRolesDashboardPermission
	graphConfig.Directives.Validate = directivesService.Validate
	graphConfig.Directives.RateLimit = directivesService.RateLimit
	graphConfig.Directives.NoRateLimit = directivesService.NoRateLimit

	schema := graph.NewExecutableSchema(graphConfig)

	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(&rateLimitExtension{
		rateLimiter: middlewaresService.RateLimitInstance(),
	})
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

	srv.Use(
		otelgqlgen.Middleware(
			otelgqlgen.WithCreateSpanFromFields(
				func(ctx *graphql.FieldContext) bool {
					return ctx.IsMethod || ctx.IsResolver
				},
			),
		),
	)

	srv.SetErrorPresenter(
		func(ctx context.Context, err error) *gqlerror.Error {
			gerr := graphql.DefaultErrorPresenter(ctx, err)

			logErr := err
			if gerr.Extensions != nil {
				if cause, ok := gerr.Extensions[gqlerrors.InternalCauseKey]; ok {
					if causeErr, ok := cause.(error); ok {
						logErr = causeErr
					}
					delete(gerr.Extensions, gqlerrors.InternalCauseKey)
				}
			}

			loggerInstance.ErrorContext(
				ctx,
				"GraphQL error",
				slog.String("path", gerr.Path.String()),
				slog.Any("extensions", gerr.Extensions),
				logger.Error(logErr),
			)

			return gerr
		},
	)

	srv.Use(extension.Introspection{})

	playgroundHandler := playground.Handler("GraphQL", "/api/query")
	httpServer.Any(
		"/", func(c *gin.Context) {
			playgroundHandler.ServeHTTP(c.Writer, c.Request)
		},
	)

	httpServer.Any(
		"/query",
		dataLoaderFactory.LoadMiddleware,
		func(c *gin.Context) {
			srv.ServeHTTP(c.Writer, c.Request)
		},
	)

	return &Gql{srv}
}
