package shortlinks

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type verifyCustomDomain struct {
	customDomainsService *shortlinkscustomdomains.Service
	sessions             *auth.Auth
	config               config.Config
}

type VerifyCustomDomainOpts struct {
	fx.In

	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
	Config               config.Config
}

func newVerifyCustomDomain(opts VerifyCustomDomainOpts) *verifyCustomDomain {
	return &verifyCustomDomain{
		customDomainsService: opts.CustomDomainsService,
		sessions:             opts.Sessions,
		config:               opts.Config,
	}
}

func (c *verifyCustomDomain) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-verify-custom-domain",
		Method:      http.MethodPost,
		Path:        "/v1/short-links/custom-domain/verify",
		Tags:        []string{"Short links"},
		Summary:     "Verify custom domain DNS configuration",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *verifyCustomDomain) Handler(
	ctx context.Context,
	input *struct{},
) (*httpbase.BaseOutputJson[customDomainOutputDto], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	err = c.customDomainsService.VerifyDomain(ctx, user.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "Domain verification failed", err)
	}

	customDomain, err := c.customDomainsService.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get custom domain", err)
	}

	return httpbase.CreateBaseOutputJson(
		mapCustomDomainOutput(
			customDomain,
			c.config.SiteBaseUrl,
		),
	), nil
}

func (c *verifyCustomDomain) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
