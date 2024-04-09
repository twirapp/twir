package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/router/logger"
)

func New() *chi.Mux {
	log := logger.New()

	router := chi.NewRouter()
	router.Use(httplog.RequestLogger(log))

	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:     []string{"*"},
				AllowedMethods:     []string{"*"},
				AllowedHeaders:     []string{"*"},
				AllowCredentials:   true,
				MaxAge:             0,
				OptionsPassthrough: true,
			},
		),
	)

	return router
}
