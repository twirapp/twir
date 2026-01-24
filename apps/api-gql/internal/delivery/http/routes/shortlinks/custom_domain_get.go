package shortlinks

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	"go.uber.org/fx"
)

type getCustomDomain struct {
	customDomainsService *shortlinkscustomdomains.Service
	sessions             *auth.Auth
}

type GetCustomDomainOpts struct {
	fx.In

	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
}

func newGetCustomDomain(opts GetCustomDomainOpts) *getCustomDomain {
	return &getCustomDomain{
		customDomainsService: opts.CustomDomainsService,
		sessions:             opts.Sessions,
	}
}

func (c *getCustomDomain) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-get-custom-domain",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/custom-domain",
		Tags:        []string{"Short links"},
		Summary:     "Get custom domain configuration",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *getCustomDomain) Handler(
	ctx context.Context,
	input *struct{},
) (*httpbase.BaseOutputJson[customDomainOutputDto], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	customDomain, err := c.customDomainsService.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get custom domain", err)
	}

	if customDomain.IsNil() {
		return nil, huma.NewError(http.StatusNotFound, "Custom domain not configured")
	}

	return httpbase.CreateBaseOutputJson(mapCustomDomainOutput(customDomain)), nil
}

func (c *getCustomDomain) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
