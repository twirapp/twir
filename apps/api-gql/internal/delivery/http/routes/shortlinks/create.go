package shortlinks

import (
	"context"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"go.uber.org/fx"
)

var _ httpbase.Route[*createLinkInput, *httpbase.BaseOutputJson[linkOutputDto]] = (*create)(nil)

type create struct {
	config   config.Config
	service  *shortenedurls.Service
	sessions *auth.Auth
	logger   logger.Logger
}

type CreateOpts struct {
	fx.In

	Config   config.Config
	Service  *shortenedurls.Service
	Sessions *auth.Auth
	Logger   logger.Logger
}

func newCreate(opts CreateOpts) *create {
	return &create{
		config:   opts.Config,
		service:  opts.Service,
		sessions: opts.Sessions,
		logger:   opts.Logger,
	}
}

func (c *create) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-url-create",
		Method:      http.MethodPost,
		Path:        "/v1/short-links",
		Tags:        []string{"Short links"},
		Summary:     "Create short url",
	}
}

func (c *create) Register(api huma.API) {
	huma.Register(
		api,
		c.GetMeta(),
		c.Handler,
	)
}

type createLinkInput struct {
	Body createLinkInputDto
}

type createLinkInputDto struct {
	Url   string `json:"url" required:"true" format:"uri" minLength:"1" maxLength:"2000" example:"https://example.com" pattern:"^https?://.*"`
	Alias string `json:"alias" required:"false" minLength:"3" maxLength:"30" example:"stream" pattern:"^[a-zA-Z0-9]+$"`
}

func (c *create) Handler(
	ctx context.Context,
	input *createLinkInput,
) (*httpbase.BaseOutputJson[linkOutputDto], error) {
	if input.Body.Alias == "" {
		existedLink, err := c.service.GetByUrl(ctx, input.Body.Url)
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
		}

		if existedLink != model.Nil {
			baseUrl, _ := url.Parse(c.config.SiteBaseUrl)
			baseUrl.Path = "/s/" + existedLink.ShortID

			return httpbase.CreateBaseOutputJson(
				linkOutputDto{
					Id:        existedLink.ShortID,
					Url:       existedLink.URL,
					ShortUrl:  baseUrl.String(),
					Views:     existedLink.Views,
					CreatedAt: existedLink.CreatedAt,
				},
			), nil
		}
	}

	if input.Body.Alias != "" {
		existedLink, err := c.service.GetByShortID(ctx, input.Body.Alias)
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
		}

		if existedLink != model.Nil {
			return nil, huma.NewError(http.StatusConflict, "Alias already in use")
		}
	}

	var createdByUserID *string
	user, _ := c.sessions.GetAuthenticatedUserModel(ctx)
	if user != nil {
		createdByUserID = &user.ID
	}

	link, err := c.service.Create(
		ctx, shortenedurls.CreateInput{
			URL:             input.Body.Url,
			ShortID:         input.Body.Alias,
			CreatedByUserID: createdByUserID,
		},
	)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Cannot generate short id", err)
	}

	baseUrl, _ := url.Parse(c.config.SiteBaseUrl)
	baseUrl.Path = "/s/" + link.ShortID

	if err := c.sessions.AddLatestShortenerUrlsId(ctx, link.ShortID); err != nil {
		c.logger.Warn("Cannot save latest short links ids to session: " + err.Error())
	}

	return httpbase.CreateBaseOutputJson(
		linkOutputDto{
			Id:        link.ShortID,
			Url:       link.URL,
			ShortUrl:  baseUrl.String(),
			Views:     link.Views,
			CreatedAt: link.CreatedAt,
		},
	), nil
}
