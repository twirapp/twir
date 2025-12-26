package brb

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/be_right_back"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	"go.uber.org/fx"
)

var _ httpbase.Route[*stopRequestDto, *httpbase.BaseOutputJson[stopResponseDto]] = (*stopEndpoint)(nil)

type StopOpts struct {
	fx.In

	Service        *be_right_back.Service
	TwirBus        *buscore.Bus
	ChannelService *channels.Service
	Middlewares    *middlewares.Middlewares
	Sessions       *auth.Auth
}

type stopEndpoint struct {
	service        *be_right_back.Service
	channelService *channels.Service
	twirBus        *buscore.Bus
	middlewares    *middlewares.Middlewares
	sessions       *auth.Auth
}

type stopResponseDto struct {
	Success bool `json:"success"`
}

func newStop(opts StopOpts) *stopEndpoint {
	return &stopEndpoint{
		service:        opts.Service,
		channelService: opts.ChannelService,
		twirBus:        opts.TwirBus,
		middlewares:    opts.Middlewares,
		sessions:       opts.Sessions,
	}
}

type stopRequestDto struct {
}

func (se *stopEndpoint) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "overlays-brb-stop",
		Method:      http.MethodPut,
		Path:        "/v1/channels/overlays/brb/stop",
		Tags:        []string{"Overlays/BRB"},
		Summary:     "Stop BRB overlay",
		Middlewares: huma.Middlewares{se.middlewares.RateLimit("brb-stop", 10, 1*time.Minute)},
		Security: []map[string][]string{
			{"api-key": {}},
		},
	}
}

func (se *stopEndpoint) Handler(
	ctx context.Context,
	input *stopRequestDto,
) (*httpbase.BaseOutputJson[stopResponseDto], error) {
	selectedDashboardId, err := se.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
	}

	if err := se.twirBus.Api.TriggerBrbStop.Publish(
		ctx,
		api.TriggerBrbStop{ChannelId: selectedDashboardId},
	); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to stop BRB overlay", err)
	}

	return httpbase.CreateBaseOutputJson(stopResponseDto{Success: true}), nil
}

func (se *stopEndpoint) Register(api huma.API) {
	huma.Register(api, se.GetMeta(), se.Handler)
}
