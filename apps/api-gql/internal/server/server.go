package server

import (
	"context"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	"github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Opts struct {
	LC                 *lifecycle.Lifecycle
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
	r.Use(newGlobalCORS(r))

	r.Use(otelgin.Middleware("api-gql"))
	r.Use(opts.Sessions.Middleware())
	r.Use(opts.Middlewares.DashboardID)
	r.Use(gin.Recovery())
	r.Use(gincontext.Middleware())
	r.Use(opts.Middlewares.RateLimit("global", 1000, 60*time.Second))

	r.NoRoute(
		func(c *gin.Context) {
			c.JSON(
				http.StatusNotFound, gin.H{
					"error":   "Not Found",
					"message": "The requested route does not exist",
					"path":    c.Request.URL.Path,
					"host":    c.Request.Host,
					"method":  c.Request.Method,
					"ip":      c.ClientIP(),
				},
			)
		},
	)

	server := &Server{
		r,
	}

	opts.LC.Append(
		lifecycle.Hook{
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
	port := os.Getenv("API_GQL_PORT")
	if port == "" {
		port = "3009"
	}
	c.Run(":" + port)
}

func (c *Server) StopServer() {
}

func newGlobalCORS(router *gin.Engine) gin.HandlerFunc {
	fallback := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	})

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			if hasExplicitRoute(router, http.MethodOptions, c.Request.URL.Path) {
				c.Next()
				return
			}

			fallback(c)
			return
		}

		fallback(c)
	}
}

func hasExplicitRoute(router *gin.Engine, method, path string) bool {
	for _, route := range router.Routes() {
		if route.Method == method && route.Path == path {
			return true
		}
	}

	return false
}
