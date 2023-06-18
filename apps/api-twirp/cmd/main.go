package main

import (
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl"
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	"github.com/twitchtv/twirp"
	"net/http"
)

func main() {
	twirpHandler := api.NewApiServer(
		&impl.Api{},
		twirp.WithServerPathPrefix("/v1"),
	)

	mux := http.NewServeMux()
	mux.Handle(
		twirpHandler.PathPrefix(),
		wrappers.Wrap(
			twirpHandler,
			wrappers.WithCors,
			wrappers.WithChannelId,
		),
	)
	http.ListenAndServe(":3002", mux)
}
