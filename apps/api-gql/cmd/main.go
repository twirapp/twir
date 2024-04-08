package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/graph"
	"github.com/twirapp/twir/apps/api-gql/resolvers"
)

func main() {
	c := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		},
	)
	config := graph.Config{
		Resolvers: &resolvers.Resolver{
			NewCommandChann: make(chan *gqlmodel.Command),
		},
	}
	config.Directives.IsAuthenticated = func(
		ctx context.Context,
		obj interface{},
		next graphql.Resolver,
	) (res interface{}, err error) {
		// user, ok := s.Get(ctx, "dbUser").(model.Users)
		// if !ok {
		// 	return nil, fmt.Errorf("not authenticated")
		// }
		//
		// ctx = context.WithValue(ctx, "dbUser", user)

		return next(ctx)
	}

	srv := handler.New(graph.NewExecutableSchema(config))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(
		transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	)
	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", c.Handler(srv))

	log.Fatal(http.ListenAndServe(":3013", nil))
}

// func main() {
// 	db, err := gorm.Open(
// 		postgres.Open("postgres://tsuwari:tsuwari@localhost:54321/tsuwari?sslmode=disable"),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	d, _ := db.DB()
// 	d.SetMaxIdleConns(1)
// 	d.SetMaxOpenConns(10)
// 	d.SetConnMaxLifetime(time.Hour)
//
// 	redisOpts, _ := redis.ParseURL("redis://localhost:6385/0")
// 	redisClient := redis.NewClient(redisOpts)
//
// 	sessionManager := sessions.New(redisClient)
//
// 	r := gin.Default()
// 	r.Use(
// 		cors.New(
// 			cors.Config{
// 				AllowAllOrigins:  true,
// 				AllowHeaders:     []string{"*"},
// 				AllowWebSockets:  true,
// 				AllowCredentials: true,
// 			},
// 		),
// 	)
// 	r.Use(ginContextToContextMiddleware())
//
// 	r.Any("/query", gin.WrapH(graphqlHandler(sessionManager)))
// 	r.GET("/", playgroundHandler())
//
// 	// r.Run(":3009")
//
// 	handler := sessionManager.LoadAndSave(r)
//
// 	panic(http.ListenAndServe(":3013", handler))
// }

func graphqlHandler(s *scs.SessionManager) *handler.Server {
	config := graph.Config{
		Resolvers: &resolvers.Resolver{
			NewCommandChann: make(chan *gqlmodel.Command),
		},
	}
	config.Directives.IsAuthenticated = func(
		ctx context.Context,
		obj interface{},
		next graphql.Resolver,
	) (res interface{}, err error) {
		user, ok := s.Get(ctx, "dbUser").(model.Users)
		if !ok {
			return nil, fmt.Errorf("not authenticated")
		}

		ctx = context.WithValue(ctx, "dbUser", user)

		return next(ctx)
	}

	h := handler.NewDefaultServer(graph.NewExecutableSchema(config))
	h.AddTransport(
		transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	)

	return h
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func ginContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
