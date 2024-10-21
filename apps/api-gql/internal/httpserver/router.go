package httpserver

import (
	"context"
	"log/slog"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/gql"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gincontext"
	"github.com/twirapp/twir/libs/cache/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC                 fx.Lifecycle
	GqlHandler         *gql.Gql
	Sessions           *auth.Auth
	CachedTwitchClient *twitch.CachedTwitchClient
	Logger             logger.Logger
}

type Server struct {
	*gin.Engine
}

func New(opts Opts) *Server {
	r := gin.New()

	r.Use(
		cors.New(
			cors.Config{
				AllowAllOrigins:  true,
				AllowMethods:     []string{"*"},
				AllowHeaders:     []string{"*"},
				ExposeHeaders:    []string{"*"},
				AllowCredentials: true,
			},
		),
	)

	r.ForwardedByClientIP = true
	r.RemoteIPHeaders = append(r.RemoteIPHeaders, "Cf-Connecting-IP")

	r.Use(opts.Sessions.Middleware())

	r.Use(
		sloggin.NewWithConfig(
			opts.Logger.GetSlog(),
			sloggin.Config{
				WithSpanID:         true,
				WithTraceID:        true,
				DefaultLevel:       slog.LevelError,
				WithResponseBody:   true,
				WithRequestBody:    true,
				WithRequestHeader:  true,
				WithUserAgent:      true,
				WithResponseHeader: true,
			},
		),
	)

	r.Use(
		func(c *gin.Context) {
			user, err := opts.Sessions.GetAuthenticatedUser(c.Request.Context())
			if err == nil {
				sloggin.AddCustomAttributes(c, slog.String("userId", user.ID))
			}
		},
	)

	r.Use(gin.Recovery())

	playgroundHandler := playground.Handler("GraphQL", "/api/query")

	r.Use(gincontext.Middleware())

	r.Any(
		"/", func(c *gin.Context) {
			playgroundHandler.ServeHTTP(c.Writer, c.Request)
		},
	)

	r.Any(
		"/query",
		func(c *gin.Context) {
			loader := data_loader.New(
				data_loader.Opts{
					CachedTwitchClient: opts.CachedTwitchClient,
				},
			)

			c.Request = c.Request.WithContext(
				context.WithValue(c.Request.Context(), data_loader.LoadersKey, loader),
			)

			c.Next()
		},
		func(c *gin.Context) {
			opts.GqlHandler.ServeHTTP(c.Writer, c.Request)
		},
	)

	server := &Server{
		r,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				opts.Logger.Info("Starting server")
				go func() {
					server.StartServer()
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				server.StopServer()
				return nil
			},
		},
	)

	return server
}

func (c *Server) StartServer() {
	c.Run(":3009")
}

func (c *Server) StopServer() {

}
