package gqlhandler

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/twirapp/twir/apps/api-gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/graph"
	"github.com/twirapp/twir/apps/api-gql/resolvers"
)

func New(mux *chi.Mux) error {
	config := graph.Config{
		Resolvers: &resolvers.Resolver{
			NewCommandChann: make(chan *gqlmodel.Command),
		},
	}
	schema := graph.NewExecutableSchema(config)

	srv := handler.New(schema)
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

	mux.Handle("/", playground.Handler("Todo", "/query"))
	mux.Handle("/query", srv)

	return nil
}
