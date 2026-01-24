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

type deleteCustomDomain struct {
	customDomainsService *shortlinkscustomdomains.Service
	sessions             *auth.Auth
}

type DeleteCustomDomainOpts struct {
	fx.In

	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
}

func newDeleteCustomDomain(opts DeleteCustomDomainOpts) *deleteCustomDomain {
	return &deleteCustomDomain{
		customDomainsService: opts.CustomDomainsService,
		sessions:             opts.Sessions,
	}
}

func (c *deleteCustomDomain) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-delete-custom-domain",
		Method:      http.MethodDelete,
		Path:        "/v1/short-links/custom-domain",
		Tags:        []string{"Short links"},
		Summary:     "Delete custom domain configuration",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *deleteCustomDomain) Handler(
	ctx context.Context,
	input *struct{},
) (*httpbase.BaseOutputJson[any], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	err = c.customDomainsService.Delete(ctx, user.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "Cannot delete custom domain", err)
	}

	return httpbase.CreateBaseOutputJson[any](map[string]string{"message": "Custom domain deleted"}), nil
}

func (c *deleteCustomDomain) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
