package shortlinks

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	config "github.com/twirapp/twir/libs/config"
	shortlinkscustomdomainsrepo "github.com/twirapp/twir/libs/repositories/short_links_custom_domains"
	"go.uber.org/fx"
)

type createCustomDomain struct {
	customDomainsService *shortlinkscustomdomains.Service
	sessions             *auth.Auth
	config               config.Config
}

type CreateCustomDomainOpts struct {
	fx.In

	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
	Config               config.Config
}

func newCreateCustomDomain(opts CreateCustomDomainOpts) *createCustomDomain {
	return &createCustomDomain{
		customDomainsService: opts.CustomDomainsService,
		sessions:             opts.Sessions,
		config:               opts.Config,
	}
}

type createCustomDomainInput struct {
	Body struct {
		Domain string `json:"domain" required:"true" minLength:"3" maxLength:"255" example:"links.example.com"`
	}
}

func (c *createCustomDomain) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-create-custom-domain",
		Method:      http.MethodPost,
		Path:        "/v1/short-links/custom-domain",
		Tags:        []string{"Short links"},
		Summary:     "Configure custom domain",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *createCustomDomain) Handler(
	ctx context.Context,
	input *createCustomDomainInput,
) (*httpbase.BaseOutputJson[customDomainOutputDto], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	customDomain, err := c.customDomainsService.Create(
		ctx, shortlinkscustomdomains.CreateInput{
			UserID: user.ID,
			Domain: input.Body.Domain,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, shortlinkscustomdomainsrepo.ErrUserAlreadyHasDomain):
			return nil, huma.NewError(http.StatusConflict, "Custom domain already configured", err)
		case errors.Is(err, shortlinkscustomdomainsrepo.ErrDomainAlreadyExists):
			return nil, huma.NewError(http.StatusConflict, "Domain is already in use", err)
		default:
			return nil, huma.NewError(http.StatusBadRequest, "Cannot create custom domain", err)
		}
	}

	return httpbase.CreateBaseOutputJson(
		mapCustomDomainOutput(
			customDomain,
			c.config.SiteBaseUrl,
		),
	), nil
}

func (c *createCustomDomain) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
