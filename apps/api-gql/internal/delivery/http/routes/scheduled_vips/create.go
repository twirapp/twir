package scheduled_vips

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/scheduledvips"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
	"go.uber.org/fx"
)

type createRequestDto struct {
	Body struct {
		UserID     string     `json:"user_id" required:"true" minLength:"1" maxLength:"100" example:"123456789" doc:"Twitch user ID"`
		RemoveAt   *time.Time `json:"remove_at,omitempty" format:"date-time" nullable:"true" doc:"When to remove VIP (for time-based removal)"`
		RemoveType string     `json:"remove_type" required:"true" enum:"time,stream_end" example:"time" doc:"Type of removal: 'time' or 'stream_end'"`
	}
}

var _ httpbase.Route[*createRequestDto, *httpbase.BaseOutputJson[scheduledVipOutputDto]] = (*create)(nil)

type CreateOpts struct {
	fx.In

	Service  *scheduledvips.Service
	Sessions *auth.Auth
}

func newCreate(opts CreateOpts) *create {
	return &create{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type create struct {
	service  *scheduledvips.Service
	sessions *auth.Auth
}

func (c *create) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "scheduled-vips-create",
		Method:      http.MethodPost,
		Path:        "/v1/scheduled-vips",
		Tags:        []string{"Scheduled VIPs"},
		Summary:     "Create scheduled VIP",
		Description: "Add a user as VIP on Twitch and schedule their removal",
		Security: []map[string][]string{
			{"api-key": {}},
		},
	}
}

func (c *create) Handler(
	ctx context.Context,
	input *createRequestDto,
) (*httpbase.BaseOutputJson[scheduledVipOutputDto], error) {
	dashboardID, err := c.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Cannot get selected dashboard", err)
	}

	// Validate remove type
	var removeType *scheduledvipsentity.RemoveType
	switch input.Body.RemoveType {
	case "time":
		if input.Body.RemoveAt == nil {
			return nil, huma.NewError(
				http.StatusBadRequest,
				"remove_at is required when remove_type is 'time'",
			)
		}
		if input.Body.RemoveAt.Before(time.Now()) {
			return nil, huma.NewError(
				http.StatusBadRequest,
				"remove_at must be in the future",
			)
		}
		removeType = lo.ToPtr(scheduledvipsentity.RemoveTypeTime)
	case "stream_end":
		removeType = lo.ToPtr(scheduledvipsentity.RemoveTypeStreamEnd)
	default:
		return nil, huma.NewError(
			http.StatusBadRequest,
			"remove_type must be 'time' or 'stream_end'",
		)
	}

	// Create scheduled VIP with Twitch VIP addition
	err = c.service.CreateWithTwitchVip(
		ctx,
		scheduledvips.CreateWithTwitchVipInput{
			UserID:     input.Body.UserID,
			ChannelID:  dashboardID,
			RemoveAt:   input.Body.RemoveAt,
			RemoveType: removeType,
		},
	)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot create scheduled VIP", err)
	}

	// Return the created VIP (fetch it back to get the ID and CreatedAt)
	vips, err := c.service.GetScheduledVips(ctx, dashboardID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "VIP created but cannot fetch details", err)
	}

	// Find the created VIP
	var createdVip *scheduledvipsentity.ScheduledVip
	for i := range vips {
		if vips[i].UserID == input.Body.UserID {
			createdVip = &vips[i]
			break
		}
	}

	if createdVip == nil {
		return nil, huma.NewError(http.StatusInternalServerError, "VIP created but cannot find it")
	}

	var removeTypeStr *string
	if createdVip.RemoveType != nil {
		removeTypeStr = lo.ToPtr(string(*createdVip.RemoveType))
	}

	return httpbase.CreateBaseOutputJson(
		scheduledVipOutputDto{
			ID:         createdVip.ID.String(),
			UserID:     createdVip.UserID,
			ChannelID:  createdVip.ChannelID,
			CreatedAt:  createdVip.CreatedAt,
			RemoveAt:   createdVip.RemoveAt,
			RemoveType: removeTypeStr,
		},
	), nil
}

func (c *create) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
