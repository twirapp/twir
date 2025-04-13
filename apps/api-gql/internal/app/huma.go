package app

import (
	"io"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
)

func NewHuma(router *server.Server, cfg config.Config, loader *dataloader.LoaderFactory) (
	huma.API,
	error,
) {
	var jsonFormat = huma.Format{
		Marshal: func(w io.Writer, v any) error {
			return json.NewEncoder(w).Encode(v)
		},
		Unmarshal: json.Unmarshal,
	}

	humaConfig := huma.DefaultConfig("Twir Api", "1.0.0")
	huma.DefaultArrayNullable = false

	serverUrl, err := url.Parse(cfg.SiteBaseUrl)
	if err != nil {
		return nil, err
	}

	serverUrl.Path = "/api"

	humaConfig.OpenAPI.Servers = []*huma.Server{
		{
			URL: serverUrl.String(),
		},
	}
	humaConfig.DocsPath = ""
	humaConfig.Formats = map[string]huma.Format{
		"application/json": jsonFormat,
		"json":             jsonFormat,
	}
	humaConfig.OpenAPIPath = "/docs/openapi"

	router.GET(
		"/docs", func(c *gin.Context) {
			c.Header("Content-Type", "text/html")
			_, _ = c.Writer.Write(docs)
		},
	)

	api := humagin.New(router.Engine, humaConfig)
	api.UseMiddleware(
		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithValue(ctx, dataloader.LoadersKey, loader.Load())

			next(ctx)
		},
	)

	return api, nil
}

var docs = []byte(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/api/docs/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`)
