package shortlinks

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	"go.uber.org/fx"
)

type deleteCustomDomain struct {
	customDomainsService *shortlinkscustomdomains.Service
	shortenedUrlsService *shortenedurls.Service
	sessions             *auth.Auth
}

type DeleteCustomDomainOpts struct {
	fx.In

	CustomDomainsService *shortlinkscustomdomains.Service
	ShortenedUrlsService *shortenedurls.Service
	Sessions             *auth.Auth
}

func newDeleteCustomDomain(opts DeleteCustomDomainOpts) *deleteCustomDomain {
	return &deleteCustomDomain{
		customDomainsService: opts.CustomDomainsService,
		shortenedUrlsService: opts.ShortenedUrlsService,
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

	customDomain, err := c.customDomainsService.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get custom domain", err)
	}
	if customDomain.IsNil() {
		return nil, huma.NewError(http.StatusNotFound, "Custom domain not found")
	}

	err = c.shortenedUrlsService.MoveLinksToDefaultDomain(ctx, user.ID, customDomain.Domain)
	if err != nil {
		if errors.Is(err, shortenedurlsrepository.ErrShortIDAlreadyExists) {
			return nil, huma.NewError(http.StatusConflict, "Short ID already exists on default domain", err)
		}
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot move links to default domain", err)
	}

	if err = c.customDomainsService.Delete(ctx, user.ID); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "Cannot delete custom domain", err)
	}

	return httpbase.CreateBaseOutputJson[any](map[string]string{"message": "Custom domain deleted"}), nil
}

func (c *deleteCustomDomain) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
