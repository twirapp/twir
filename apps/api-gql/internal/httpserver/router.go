package httpserver

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/api-gql/internal/gql"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"github.com/twirapp/twir/libs/cache/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC                 fx.Lifecycle
	GqlHandler         *gql.Gql
	Sessions           *sessions.Sessions
	CachedTwitchClient *twitch.CachedTwitchClient
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

	r.Use(opts.Sessions.Middleware())

	// r.Use(
	// 	func(c *gin.Context) {
	// 		fmt.Println(opts.Sessions.GetAuthenticatedUser(c.Request.Context()))
	// 	},
	// )

	playgroundHandler := playground.Handler("GraphQL", "/query")

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
