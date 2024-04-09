package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func New() *chi.Mux {
	router := chi.NewRouter()

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
