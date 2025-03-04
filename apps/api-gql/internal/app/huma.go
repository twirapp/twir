package app

import (
	"io"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
)

func NewHuma(router *server.Server, loader *dataloader.LoaderFactory) huma.API {
	var jsonFormat = huma.Format{
		Marshal: func(w io.Writer, v any) error {
			return json.NewEncoder(w).Encode(v)
		},
		Unmarshal: json.Unmarshal,
	}

	config := huma.DefaultConfig("Twir Api", "1.0.0")
	huma.DefaultArrayNullable = false
	config.Formats = map[string]huma.Format{
		"application/json": jsonFormat,
		"json":             jsonFormat,
	}

	// GET /greetings/{id}
	// PUT /greetings/{id}

	api := humagin.New(router.Engine, config)
	api.UseMiddleware(
		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithValue(ctx, dataloader.LoadersKey, loader.Load())

			next(ctx)
		},
	)

	return api
}
