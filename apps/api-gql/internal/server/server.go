package server

import (
	"context"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	"github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
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
	Config             config.Config
}

type Server struct {
	*gin.Engine
}

func New(opts Opts) (*Server, error) {
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

	err := r.SetTrustedProxies(
		append(
			opts.Config.TrustedProxies,
			"127.0.0.1",
			"::1",
			"172.17.0.0/16",
			"172.18.0.0/16",
			// docker 10
			"10.0.2.0/24",
			"10.0.1.0/24",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}

	r.ForwardedByClientIP = true
	r.RemoteIPHeaders = append(r.RemoteIPHeaders, "Cf-Connecting-IP", "X-Forwarded-For", "X-Real-IP")

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

	return server, nil
}

func (c *Server) StartServer() {
	c.Run(":3009")
}

func (c *Server) StopServer() {

}
