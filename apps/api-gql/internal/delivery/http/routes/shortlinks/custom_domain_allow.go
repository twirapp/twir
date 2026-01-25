package shortlinks

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	"go.uber.org/fx"
)

type allowCustomDomain struct {
	customDomainsService *shortlinkscustomdomains.Service
}

type AllowCustomDomainOpts struct {
	fx.In

	CustomDomainsService *shortlinkscustomdomains.Service
}

func newAllowCustomDomain(opts AllowCustomDomainOpts) *allowCustomDomain {
	return &allowCustomDomain{
		customDomainsService: opts.CustomDomainsService,
	}
}

type allowCustomDomainInput struct {
	Domain string `query:"domain" required:"true" minLength:"1"`
}

type allowCustomDomainOutput struct {
	Allowed bool `json:"allowed"`
}

func (a *allowCustomDomain) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-custom-domain-allow",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/custom-domain/allow",
		Tags:        []string{"Short links"},
		Summary:     "Check if custom domain is allowed for TLS",
	}
}

func (a *allowCustomDomain) Handler(
	ctx context.Context,
	input *allowCustomDomainInput,
) (*httpbase.BaseOutputJson[allowCustomDomainOutput], error) {
	allowed, err := a.customDomainsService.IsDomainAllowed(ctx, input.Domain)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot check custom domain", err)
	}

	if !allowed {
		return nil, huma.NewError(http.StatusNotFound, "Custom domain not allowed")
	}

	return httpbase.CreateBaseOutputJson(allowCustomDomainOutput{Allowed: true}), nil
}

func (a *allowCustomDomain) Register(api huma.API) {
	huma.Register(api, a.GetMeta(), a.Handler)
}
