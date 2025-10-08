package shortlinks

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"go.uber.org/fx"
)

type redirectRequestDto struct {
	ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
}

type redirectResponseDto struct {
	Status   int
	Location string `header:"Location"`
}

var _ httpbase.Route[*redirectRequestDto, *redirectResponseDto] = (*redirect)(nil)

type RedirectOpts struct {
	fx.In

	Service *shortenedurls.Service
	Config  config.Config
}

func newRedirect(opts RedirectOpts) *redirect {
	return &redirect{
		service: opts.Service,
		config:  opts.Config,
	}
}

type redirect struct {
	service *shortenedurls.Service
	config  config.Config
}

func (r *redirect) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "short-url-redirect",
		Method:        http.MethodGet,
		Path:          "/v1/short-links/{shortId}",
		Tags:          []string{"Short links"},
		Summary:       "Redirect to url",
		DefaultStatus: 301,
	}
}

func (r *redirect) Handler(ctx context.Context, input *redirectRequestDto) (
	*redirectResponseDto,
	error,
) {
	link, err := r.service.GetByShortID(ctx, input.ShortId)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
	}

	if link == model.Nil {
		return nil, huma.NewError(http.StatusNotFound, "Link not found")
	}

	newViews := link.Views + 1

	if err := r.service.Update(
		ctx,
		link.ShortID,
		shortenedurls.UpdateInput{
			Views: &newViews,
		},
	); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot update link", err)
	}

	return &redirectResponseDto{
		Status:   http.StatusPermanentRedirect,
		Location: link.URL,
	}, nil
}

func (r *redirect) Register(api huma.API) {
	huma.Register(
		api,
		r.GetMeta(),
		r.Handler,
	)
}
