package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	"github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC                 fx.Lifecycle
	Sessions           *auth.Auth
	CachedTwitchClient *twitch.CachedTwitchClient
	Logger             *slog.Logger
	Middlewares        *middlewares.Middlewares
	Config             config.Config
}

type Server struct {
	*gin.Engine
}

func New(opts Opts) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(opts.Middlewares.Logger())
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

	r.Use(otelgin.Middleware("api-gql"))
	r.Use(opts.Sessions.Middleware())
	r.Use(opts.Middlewares.DashboardID)
	r.Use(gin.Recovery())
	r.Use(gincontext.Middleware())
	r.Use(opts.Middlewares.RateLimit("global", 1000, 60*time.Second))

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

	return server, nil
}

func (c *Server) StartServer() {
	c.Run(":3009")
}

func (c *Server) StopServer() {

}
