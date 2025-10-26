package shortlinks

import (
	"context"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
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
	link, err := i.service.GetByShortID(ctx, input.ShortId)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
	}

	if link == model.Nil {
		return nil, huma.NewError(http.StatusNotFound, "Link not found")
	}

	baseUrl, _ := url.Parse(i.config.SiteBaseUrl)
	baseUrl.Path = "/s/" + input.ShortId

	return httpbase.CreateBaseOutputJson(
		linkOutputDto{
			Id:       link.ShortID,
			Url:      link.URL,
			ShortUrl: baseUrl.String(),
			Views:    link.Views,
		},
	), nil
}
