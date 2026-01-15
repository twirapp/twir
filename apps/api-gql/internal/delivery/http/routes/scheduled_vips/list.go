package scheduled_vips

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/scheduledvips"
	"go.uber.org/fx"
)

type listRequestDto struct{}

type listResponseDto struct {
	Body struct {
		Data []scheduledVipOutputDto `json:"data"`
	}
}

var _ httpbase.Route[*listRequestDto, *listResponseDto] = (*list)(nil)

type ListOpts struct {
	fx.In

	Service  *scheduledvips.Service
	Sessions *auth.Auth
}

func newList(opts ListOpts) *list {
	return &list{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type list struct {
	service  *scheduledvips.Service
	sessions *auth.Auth
}

func (l *list) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "scheduled-vips-list",
		Method:      http.MethodGet,
		Path:        "/v1/scheduled-vips",
		Tags:        []string{"Scheduled VIPs"},
		Summary:     "List scheduled VIPs",
		Description: "Get all scheduled VIPs for the selected dashboard",
		Security: []map[string][]string{
			{"api-key": {}},
		},
	}
}

func (l *list) Handler(
	ctx context.Context,
	input *listRequestDto,
) (*listResponseDto, error) {
	dashboardID, err := l.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Cannot get selected dashboard", err)
	}

	vips, err := l.service.GetScheduledVips(ctx, dashboardID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get scheduled VIPs", err)
	}

	result := make([]scheduledVipOutputDto, len(vips))
	for i, vip := range vips {
		var removeTypeStr *string
		if vip.RemoveType != nil {
			removeTypeStr = lo.ToPtr(string(*vip.RemoveType))
		}

		result[i] = scheduledVipOutputDto{
			ID:         vip.ID.String(),
			UserID:     vip.UserID,
			ChannelID:  vip.ChannelID,
			CreatedAt:  vip.CreatedAt,
			RemoveAt:   vip.RemoveAt,
			RemoveType: removeTypeStr,
		}
	}

	return &listResponseDto{
		Body: struct {
			Data []scheduledVipOutputDto `json:"data"`
		}{
			Data: result,
		},
	}, nil
}

func (l *list) Register(api huma.API) {
	huma.Register(api, l.GetMeta(), l.Handler)
}
