package server

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC                 fx.Lifecycle
	Sessions           *auth.Auth
	CachedTwitchClient *twitch.CachedTwitchClient
	Logger             logger.Logger
	Middlewares        *middlewares.Middlewares
}

type Server struct {
	*gin.Engine
}

func New(opts Opts) *Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
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

	r.Use(otelgin.Middleware("api-gql"))
	r.Use(opts.Sessions.Middleware())
	r.Use(opts.Middlewares.Logging)
	r.Use(opts.Middlewares.DashboardID)
	r.Use(gin.Recovery())
	r.Use(gincontext.Middleware())

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
