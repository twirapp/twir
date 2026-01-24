package shortlinks

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	"go.uber.org/fx"
)

type deleteRequestDto struct {
	ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
}

var _ httpbase.Route[*deleteRequestDto, *httpbase.BaseOutputJson[any]] = (*deleteRoute)(nil)

type DeleteOpts struct {
	fx.In

	Service              *shortenedurls.Service
	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
}

type deleteRoute struct {
	service              *shortenedurls.Service
	customDomainsService *shortlinkscustomdomains.Service
	sessions             *auth.Auth
}

func newDelete(opts DeleteOpts) *deleteRoute {
	return &deleteRoute{
		service:              opts.Service,
		customDomainsService: opts.CustomDomainsService,
		sessions:             opts.Sessions,
	}
}

func (d *deleteRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-url-delete",
		Method:      http.MethodDelete,
		Path:        "/v1/short-links/{shortId}",
		Tags:        []string{"Short links"},
		Summary:     "Delete short url",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (d *deleteRoute) Register(api huma.API) {
	huma.Register(
		api,
		d.GetMeta(),
		d.Handler,
	)
}

func (d *deleteRoute) Handler(
	ctx context.Context,
	input *deleteRequestDto,
) (*httpbase.BaseOutputJson[any], error) {
	user, err := d.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized", err)
	}

	var domain *string
	if userDomain, err := d.customDomainsService.GetByUserID(ctx, user.ID); err == nil && !userDomain.IsNil() && userDomain.Verified {
		domain = &userDomain.Domain
	}

	// Get the link to verify ownership
	link, err := d.service.GetByShortID(ctx, domain, input.ShortId)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get link", err)
	}
	if link.IsNil() && domain != nil {
		domain = nil
		link, err = d.service.GetByShortID(ctx, domain, input.ShortId)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot get link", err)
		}
	}
	if link.IsNil() {
		return nil, huma.NewError(http.StatusNotFound, "Link not found")
	}

	// Check if user owns this link
	if link.CreatedByUserId == nil || *link.CreatedByUserId != user.ID {
		return nil, huma.NewError(http.StatusForbidden, "You don't have permission to delete this link")
	}

	// Delete the link
	err = d.service.Delete(ctx, domain, input.ShortId)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete link", err)
	}

	return httpbase.CreateBaseOutputJson[any](map[string]string{"message": "Link deleted successfully"}), nil
}
