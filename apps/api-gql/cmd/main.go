package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	"github.com/twirapp/twir/apps/api-gql/resolvers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(
		postgres.Open("postgres://tsuwari:tsuwari@localhost:54321/tsuwari?sslmode=disable"),
	)
	if err != nil {
		panic(err)
	}
	d, _ := db.DB()
	d.SetMaxIdleConns(1)
	d.SetMaxOpenConns(10)
	d.SetConnMaxLifetime(time.Hour)

	redisOpts, _ := redis.ParseURL("redis://localhost:6385/0")
	redisClient := redis.NewClient(redisOpts)

	sessionManager := sessions.New(redisClient)

	r := gin.Default()
	r.Use(
		cors.New(
			cors.Config{
				AllowAllOrigins:  true,
				AllowHeaders:     []string{"*"},
				AllowWebSockets:  true,
				AllowCredentials: true,
			},
		),
	)
	r.Use(ginContextToContextMiddleware())

	r.POST("/query", graphqlHandler(sessionManager))
	r.GET("/", playgroundHandler())

	// r.Run(":3009")

	handler := sessionManager.LoadAndSave(r)

	panic(http.ListenAndServe(":3011", handler))
}

func graphqlHandler(s *scs.SessionManager) gin.HandlerFunc {
	config := graph.Config{Resolvers: &resolvers.Resolver{}}
	config.Directives.IsAuthenticated = func(
		ctx context.Context,
		obj interface{},
		next graphql.Resolver,
	) (res interface{}, err error) {
		fmt.Println("1")
		user, ok := s.Get(ctx, "dbUser").(model.Users)
		if !ok {
			return nil, fmt.Errorf("not authenticated")
		}

		ctx = context.WithValue(ctx, "dbUser", user)

		return next(ctx)
	}

	h := handler.NewDefaultServer(graph.NewExecutableSchema(config))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
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
