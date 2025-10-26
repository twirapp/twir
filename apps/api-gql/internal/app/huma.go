package app

import (
	"io"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type HumaOpts struct {
	fx.In

	Router   *server.Server
	Cfg      config.Config
	Loader   *dataloader.LoaderFactory
	Sessions *auth.Auth
}

func NewHuma(opts HumaOpts) (
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

	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		// Example alternative describing the use of JWTs without documenting how
		// they are issued or which flows might be supported. This is simpler but
		// tells clients less information.
		"api-key": {
			Type:        "apiKey",
			In:          "header",
			Description: "Api key from twir dashboard",
		},
	}

	serverUrl, err := url.Parse(opts.Cfg.SiteBaseUrl)
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

	opts.Router.GET(
		"/docs", func(c *gin.Context) {
			c.Header("Content-Type", "text/html")
			_, _ = c.Writer.Write(docs)
		},
	)

	api := humagin.New(opts.Router.Engine, humaConfig)
	api.UseMiddleware(
		func(ctx huma.Context, next func(huma.Context)) {
			ctx = huma.WithValue(ctx, dataloader.LoadersKey, opts.Loader.Load())

			next(ctx)
		},
	)
	api.UseMiddleware(NewAuthMiddleware(api, opts.Sessions))

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
