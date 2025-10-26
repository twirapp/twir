package pastebins

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"go.uber.org/fx"
)

type deleteRequestDto struct {
	ID string `path:"id" maxLength:"5" minLength:"1" pattern:"^[-_a-zA-Z0-9]+$" required:"true"`
}

type deleteResponseDto struct {
	Status int `json:"status" example:"ok"`
}

var _ httpbase.Route[*deleteRequestDto, *deleteResponseDto] = (*deleteRoute)(nil)

type DeleteOpts struct {
	fx.In

	Service  *pastebins.Service
	Sessions *auth.Auth
}

func newDelete(opts CreateOpts) *deleteRoute {
	return &deleteRoute{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type deleteRoute struct {
	service  *pastebins.Service
	sessions *auth.Auth
}

func (d *deleteRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "pastebin-delete",
		Method:      http.MethodDelete,
		Path:        "/v1/pastebin/{id}",
		Tags:        []string{"Pastebin"},
		Summary:     "Delete pastebin",
		Security: []map[string][]string{
			{"api-key": {}},
		},
	}
}

func (d *deleteRoute) Handler(ctx context.Context, input *deleteRequestDto) (
	*deleteResponseDto,
	error,
) {
	paste, err := d.service.GetByID(ctx, input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Cannot get pastebin", err)
	}

	user, err := d.sessions.GetAuthenticatedUserModel(ctx)
	if user == nil || err != nil || paste.OwnerUserID == nil || *paste.OwnerUserID != user.ID {
		return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
	}

	if err := d.service.Delete(ctx, input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete pastebin", err)
	}

	return &deleteResponseDto{
		Status: http.StatusNoContent,
	}, nil
}

func (d *deleteRoute) Register(api huma.API) {
	huma.Register(api, d.GetMeta(), d.Handler)
}
