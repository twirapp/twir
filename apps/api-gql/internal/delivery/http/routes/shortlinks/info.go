package shortlinks

import (
	"context"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	humahelpers "github.com/twirapp/twir/apps/api-gql/internal/server/huma_helpers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type infoRequestDto struct {
	ShortId string `query:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
}

var _ httpbase.Route[*infoRequestDto, *httpbase.BaseOutputJson[linkOutputDto]] = (*info)(nil)

type InfoOpts struct {
	fx.In

	Service *shortenedurls.Service
	Config  config.Config
}

type info struct {
	service *shortenedurls.Service
	config  config.Config
}

func newInfo(opts InfoOpts) *info {
	return &info{
		service: opts.Service,
		config:  opts.Config,
	}
}

func (i *info) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-url-get-info",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/{shortId}/info",
		Tags:        []string{"Short links"},
		Summary:     "Get short url data",
	}
}

func (i *info) Register(api huma.API) {
	huma.Register(
		api,
		i.GetMeta(),
		i.Handler,
	)
}

func (i *info) Handler(
	ctx context.Context,
	input *infoRequestDto,
) (*httpbase.BaseOutputJson[linkOutputDto], error) {
	baseUrl, err := gincontext.GetBaseUrlFromContext(ctx, i.config.SiteBaseUrl)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get base URL", err)
	}

	var domain *string
	if host, err := humahelpers.GetHostFromCtx(ctx); err == nil && !isDefaultDomain(i.config.SiteBaseUrl, host) {
		domain = &host
	}

	link, err := i.service.GetByShortID(ctx, domain, input.ShortId)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
	}

	if link.IsNil() {
		return nil, huma.NewError(http.StatusNotFound, "Link not found")
	}

	var shortURL string
	if domain != nil {
		shortURL = "https://" + *domain + "/" + input.ShortId
	} else {
		parsedBaseUrl, _ := url.Parse(baseUrl)
		parsedBaseUrl.Path = "/s/" + input.ShortId
		shortURL = parsedBaseUrl.String()
	}

	return httpbase.CreateBaseOutputJson(
		linkOutputDto{
			Id:       link.ShortID,
			Url:      link.URL,
			ShortUrl: shortURL,
			Views:    link.Views,
		},
	), nil
}
