package shortlinks

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	humahelpers "github.com/twirapp/twir/apps/api-gql/internal/server/huma_helpers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	config "github.com/twirapp/twir/libs/config"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	"go.uber.org/fx"
)

var _ httpbase.Route[*createLinkInput, *httpbase.BaseOutputJson[linkOutputDto]] = (*create)(nil)

type create struct {
	config               config.Config
	service              *shortenedurls.Service
	customDomainsService *shortlinkscustomdomains.Service
	sessions             *auth.Auth
	logger               *slog.Logger
	middlewares          *middlewares.Middlewares
}

type CreateOpts struct {
	fx.In

	Config               config.Config
	Service              *shortenedurls.Service
	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
	Logger               *slog.Logger
	Middlewares          *middlewares.Middlewares
}

func newCreate(opts CreateOpts) *create {
	return &create{
		config:               opts.Config,
		service:              opts.Service,
		customDomainsService: opts.CustomDomainsService,
		sessions:             opts.Sessions,
		logger:               opts.Logger,
		middlewares:          opts.Middlewares,
	}
}

func (c *create) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-url-create",
		Method:      http.MethodPost,
		Path:        "/v1/short-links",
		Tags:        []string{"Short links"},
		Summary:     "Create short url",
		Middlewares: huma.Middlewares{c.middlewares.RateLimit("short-url-create", 10, 1*time.Minute)},
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
	Url             string `json:"url"   required:"true"  format:"uri" minLength:"1" maxLength:"2000" example:"https://example.com" pattern:"^https?://.*"`
	Alias           string `json:"alias" required:"false"              minLength:"3" maxLength:"30"   example:"stream"              pattern:"^[a-zA-Z0-9]+$"`
	UseCustomDomain *bool  `json:"use_custom_domain,omitempty"`
}

func (c *create) Handler(
	ctx context.Context,
	input *createLinkInput,
) (*httpbase.BaseOutputJson[linkOutputDto], error) {
	baseUrl, err := gincontext.GetBaseUrlFromContext(ctx, c.config.SiteBaseUrl)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get base URL", err)
	}

	var customDomain *string
	var createdByUserID *string
	useCustomDomain := input.Body.UseCustomDomain != nil && *input.Body.UseCustomDomain
	user, _ := c.sessions.GetAuthenticatedUserModel(ctx)
	if user != nil {
		createdByUserID = &user.ID

		userDomain, err := c.customDomainsService.GetByUserID(ctx, user.ID)
		if err == nil && !userDomain.IsNil() && userDomain.Verified {
			if input.Body.UseCustomDomain == nil || useCustomDomain {
				customDomain = &userDomain.Domain
			}
		} else if useCustomDomain {
			return nil, huma.NewError(http.StatusBadRequest, "Custom domain is not available")
		}
	} else if useCustomDomain {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	if input.Body.Alias == "" {
		existedLink, err := c.service.GetByUrl(ctx, customDomain, input.Body.Url)
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
		}

		if !existedLink.IsNil() {
			var shortURL string
			if customDomain != nil {
				shortURL = "https://" + *customDomain + "/" + existedLink.ShortID
			} else {
				parsedBaseUrl, _ := url.Parse(baseUrl)
				parsedBaseUrl.Path = "/s/" + existedLink.ShortID
				shortURL = parsedBaseUrl.String()
			}

			sessionKey := shortenedurls.EncodeShortLinkKey(customDomain, existedLink.ShortID)
			if err := c.sessions.AddLatestShortenerUrlsId(ctx, sessionKey); err != nil {
				c.logger.Warn("Cannot save latest short links ids to session: " + err.Error())
			}

			return httpbase.CreateBaseOutputJson(
				linkOutputDto{
					Id:        existedLink.ShortID,
					Url:       existedLink.URL,
					ShortUrl:  shortURL,
					Views:     existedLink.Views,
					CreatedAt: existedLink.CreatedAt,
				},
			), nil
		}
	}

	if input.Body.Alias != "" {
		existedLink, err := c.service.GetByShortID(ctx, customDomain, input.Body.Alias)
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
		}

		if !existedLink.IsNil() {
			return nil, huma.NewError(http.StatusConflict, "Alias already in use")
		}
	}

	clientIp, err := humahelpers.GetClientIpFromCtx(ctx)
	if err != nil {
		return nil, huma.NewError(
			http.StatusInternalServerError,
			"Internal error on getting your information",
			err,
		)
	}

	clientAgent, err := humahelpers.GetClientUserAgentFromCtx(ctx)
	if err != nil {
		return nil, huma.NewError(
			http.StatusInternalServerError,
			"Internal error on getting your information",
			err,
		)
	}

	link, err := c.service.Create(
		ctx, shortenedurls.CreateInput{
			CreatedByUserID: createdByUserID,
			ShortID:         input.Body.Alias,
			URL:             input.Body.Url,
			UserIp:          &clientIp,
			UserAgent:       &clientAgent,
			Domain:          customDomain,
		},
	)
	if err != nil {
		if errors.Is(err, shortenedurlsrepository.ErrShortIDAlreadyExists) {
			return nil, huma.NewError(http.StatusConflict, "Alias already in use", err)
		}
		return nil, huma.NewError(http.StatusNotFound, "Cannot generate short id", err)
	}

	var shortURL string
	if customDomain != nil {
		shortURL = "https://" + *customDomain + "/" + link.ShortID
	} else {
		parsedBaseUrl, _ := url.Parse(baseUrl)
		parsedBaseUrl.Path = "/s/" + link.ShortID
		shortURL = parsedBaseUrl.String()
	}

	sessionKey := shortenedurls.EncodeShortLinkKey(customDomain, link.ShortID)
	if err := c.sessions.AddLatestShortenerUrlsId(ctx, sessionKey); err != nil {
		c.logger.Warn("Cannot save latest short links ids to session: " + err.Error())
	}

	return httpbase.CreateBaseOutputJson(
		linkOutputDto{
			Id:        link.ShortID,
			Url:       link.URL,
			ShortUrl:  shortURL,
			Views:     link.Views,
			CreatedAt: link.CreatedAt,
		},
	), nil
}
