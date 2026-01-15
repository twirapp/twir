package scheduled_vips

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/scheduledvips"
	"go.uber.org/fx"
)

type deleteRequestDto struct {
	ID string `path:"id" required:"true" minLength:"36" maxLength:"36" pattern:"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UUID of the scheduled VIP"`
}

type deleteResponseDto struct {
	Status int `json:"status" example:"204"`
}

var _ httpbase.Route[*deleteRequestDto, *deleteResponseDto] = (*deleteRoute)(nil)

type DeleteOpts struct {
	fx.In

	Service  *scheduledvips.Service
	Sessions *auth.Auth
}

func newDelete(opts DeleteOpts) *deleteRoute {
	return &deleteRoute{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type deleteRoute struct {
	service  *scheduledvips.Service
	sessions *auth.Auth
}

func (d *deleteRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "scheduled-vips-delete",
		Method:      http.MethodDelete,
		Path:        "/v1/scheduled-vips/{id}",
		Tags:        []string{"Scheduled VIPs"},
		Summary:     "Delete scheduled VIP",
		Description: "Remove a scheduled VIP. Note: This only removes the schedule, not the VIP status on Twitch.",
		Security: []map[string][]string{
			{"api-key": {}},
		},
	}
}

func (d *deleteRoute) Handler(
	ctx context.Context,
	input *deleteRequestDto,
) (*deleteResponseDto, error) {
	dashboardID, err := d.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Cannot get selected dashboard", err)
	}

	err = d.service.Remove(
		ctx,
		scheduledvips.RemoveInput{
			ID:        input.ID,
			ChannelID: dashboardID,
		},
	)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete scheduled VIP", err)
	}

	return &deleteResponseDto{
		Status: http.StatusNoContent,
	}, nil
}

func (d *deleteRoute) Register(api huma.API) {
	huma.Register(api, d.GetMeta(), d.Handler)
}
