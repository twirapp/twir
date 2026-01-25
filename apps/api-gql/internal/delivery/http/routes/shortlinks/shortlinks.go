package shortlinks

import (
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Api      huma.API
	Config   config.Config
	Service  *shortenedurls.Service
	Sessions *auth.Auth
	Logger   *slog.Logger
}

var FxModule = fx.Provide(
	httpbase.AsFxRoute(newCreate),
	httpbase.AsFxRoute(newInfo),
	httpbase.AsFxRoute(newRedirect),
	httpbase.AsFxRoute(newProfile),
	httpbase.AsFxRoute(newStatistics),
	httpbase.AsFxRoute(newTopCountries),
	httpbase.AsFxRoute(newUpdate),
	httpbase.AsFxRoute(newDelete),
	httpbase.AsFxRoute(newGetCustomDomain),
	httpbase.AsFxRoute(newCreateCustomDomain),
	httpbase.AsFxRoute(newVerifyCustomDomain),
	httpbase.AsFxRoute(newDeleteCustomDomain),
	httpbase.AsFxRoute(newAllowCustomDomain),
)

type linkOutputDto struct {
	Id        string    `json:"id" example:"KKMEa"`
	Url       string    `json:"url" example:"https://example.com"`
	ShortUrl  string    `json:"short_url" example:"https://twir.app/s/KKMEa"`
	Views     int       `json:"views" example:"1"`
	CreatedAt time.Time `json:"created_at" format:"date-time" example:"2023-01-01T00:00:00Z"`
}
