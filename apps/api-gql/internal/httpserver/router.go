package httpserver

import (
	"context"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/api-gql/internal/gql"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC         fx.Lifecycle
	GqlHandler *gql.Gql
	Sessions   *sessions.Sessions
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
		"/query", func(c *gin.Context) {
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
