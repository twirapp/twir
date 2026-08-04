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

type Server struct {
	*gin.Engine
}

func New(lc *lifecycle.Lifecycle, sessions *auth.Auth, cachedTwitchClient *twitch.CachedTwitchClient, logger *slog.Logger, middlewaresService *middlewares.Middlewares, config config.Config) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(middlewaresService.Logger())
	r.Use(newGlobalCORS(r))

	r.Use(otelgin.Middleware("api-gql"))
	r.Use(sessions.Middleware())
	r.Use(middlewaresService.DashboardID)
	r.Use(gin.Recovery())
	r.Use(gincontext.Middleware())
	r.Use(middlewaresService.RateLimit("global", 1000, 60*time.Second))

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

	srv := &Server{
		r,
	}

	lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("Starting server")
				go func() {
					srv.StartServer()
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				srv.StopServer()
				return nil
			},
		},
	)

	return srv, nil
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
