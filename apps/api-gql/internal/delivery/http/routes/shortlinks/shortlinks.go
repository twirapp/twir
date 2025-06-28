package shortlinks

import (
	"context"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Api     huma.API
	Config  config.Config
	Service *shortenedurls.Service
}

func New(opts Opts) {
	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "short-url-create",
			Method:      http.MethodPost,
			Path:        "/v1/short-links",
			Tags:        []string{"Short links"},
			Summary:     "Create short url",
		}, func(
			ctx context.Context,
			input *createLinkInput,
		) (
			*createLinkOutput, error,
		) {
			existedLink, err := opts.Service.GetByUrl(ctx, input.Body.Url)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
			}

			if existedLink != model.Nil {
				baseUrl, _ := url.Parse(opts.Config.SiteBaseUrl)
				baseUrl.Path = "/s/" + existedLink.ShortID

				return &createLinkOutput{
					Body: linkOutputDto{
						Id:       existedLink.ShortID,
						Url:      existedLink.URL,
						ShortUrl: baseUrl.String(),
						Views:    existedLink.Views,
					},
				}, nil
			}

			link, err := opts.Service.Create(
				ctx, shortenedurls.CreateInput{
					URL: input.Body.Url,
				},
			)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot generate short id", err)
			}

			baseUrl, _ := url.Parse(opts.Config.SiteBaseUrl)
			baseUrl.Path = "/s/" + link.ShortID

			return &createLinkOutput{
				Body: linkOutputDto{
					Id:       link.ShortID,
					Url:      input.Body.Url,
					ShortUrl: baseUrl.String(),
					Views:    link.Views,
				},
			}, nil
		},
	)

	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "short-url-get-info",
			Method:      http.MethodGet,
			Path:        "/v1/short-links",
			Tags:        []string{"Short links"},
			Summary:     "Get short url data",
		}, func(
			ctx context.Context,
			input *struct {
				ShortId string `query:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
			},
		) (
			*createLinkOutput, error,
		) {
			link, err := opts.Service.GetByShortID(ctx, input.ShortId)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
			}

			if link == model.Nil {
				return nil, huma.NewError(http.StatusNotFound, "Link not found")
			}

			baseUrl, _ := url.Parse(opts.Config.SiteBaseUrl)
			baseUrl.Path = "/s/" + input.ShortId

			return &createLinkOutput{
				Body: linkOutputDto{
					Id:       link.ShortID,
					Url:      link.URL,
					ShortUrl: baseUrl.String(),
				},
			}, nil
		},
	)

	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID:   "short-url-redirect",
			Method:        http.MethodGet,
			Path:          "/v1/short-links/{shortId}",
			Tags:          []string{"Short links"},
			Summary:       "Redirect to url",
			DefaultStatus: 301,
		},
		func(
			ctx context.Context, input *struct {
				ShortId string `path:"shortId" maxLength:"5" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
			},
		) (
			*linkRedirectOutput, error,
		) {
			link, err := opts.Service.GetByShortID(ctx, input.ShortId)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
			}

			if link == model.Nil {
				return nil, huma.NewError(http.StatusNotFound, "Link not found")
			}

			if err := opts.Service.Update(
				ctx,
				link.ShortID,
				shortenedurls.UpdateInput{
					Views: &link.Views,
				},
			); err != nil {
				return nil, huma.NewError(http.StatusInternalServerError, "Cannot update link", err)
			}

			return &linkRedirectOutput{
				Status:   http.StatusPermanentRedirect,
				Location: link.URL,
			}, nil
		},
	)
}

type createLinkInput struct {
	Body createLinkInputDto
}

type createLinkInputDto struct {
	Url string `json:"url" required:"true" format:"uri" minLength:"1" maxLength:"2000" example:"https://example.com" pattern:"^https?://.*"`
}

type createLinkOutput struct {
	Body linkOutputDto
}

type linkOutputDto struct {
	Id       string `json:"id" example:"KKMEa"`
	Url      string `json:"url" example:"https://example.com"`
	ShortUrl string `json:"short_url" example:"https://twir.app/s/KKMEa"`
	Views    int    `json:"views" example:"1"`
}

type linkRedirectOutput struct {
	Status   int
	Location string `header:"Location"`
}
