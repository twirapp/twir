package shortlinks

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"github.com/twirapp/twir/libs/repositories/shortened_urls"
	"go.uber.org/fx"
)

type deleteRequestDto struct {
	ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
}

var _ httpbase.Route[*deleteRequestDto, *httpbase.BaseOutputJson[any]] = (*deleteRoute)(nil)

type DeleteOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

type deleteRoute struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

func newDelete(opts DeleteOpts) *deleteRoute {
	return &deleteRoute{
		service:  opts.Service,
		sessions: opts.Sessions,
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

	// Get the link to verify ownership
	link, err := d.service.GetByShortID(ctx, input.ShortId)
	if err != nil {
		if err == shortened_urls.ErrNotFound {
			return nil, huma.NewError(http.StatusNotFound, "Link not found")
		}
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get link", err)
	}

	// Check if user owns this link
	if link.CreatedByUserId == nil || *link.CreatedByUserId != user.ID {
		return nil, huma.NewError(http.StatusForbidden, "You don't have permission to delete this link")
	}

	// Delete the link
	err = d.service.Delete(ctx, input.ShortId)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete link", err)
	}

	return httpbase.CreateBaseOutputJson[any](map[string]string{"message": "Link deleted successfully"}), nil
}
