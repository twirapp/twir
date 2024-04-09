package main

import (
	"log/slog"
	"net/http"

	"github.com/twirapp/twir/apps/api-gql/internal/gqlhandler"
	"github.com/twirapp/twir/apps/api-gql/internal/router"
)

func main() {
	mux := router.New()
	if err := gqlhandler.New(mux); err != nil {
		panic(err)
	}

	s := &http.Server{
		Addr:    ":3009",
		Handler: mux,
	}

	slog.Info("Server is running", slog.String("addr", s.Addr))

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
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
