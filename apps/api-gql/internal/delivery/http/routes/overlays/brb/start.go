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

var _ httpbase.Route[*startRequestDto, *httpbase.BaseOutputJson[startResponseDto]] = (*startEndpoint)(nil)

type StartOpts struct {
	fx.In

	Service        *be_right_back.Service
	TwirBus        *buscore.Bus
	ChannelService *channels.Service
	Middlewares    *middlewares.Middlewares
	Sessions       *auth.Auth
}

type startEndpoint struct {
	service        *be_right_back.Service
	channelService *channels.Service
	twirBus        *buscore.Bus
	middlewares    *middlewares.Middlewares
	sessions       *auth.Auth
}

type startResponseDto struct {
	Success bool `json:"success"`
}

func newStart(opts StartOpts) *startEndpoint {
	return &startEndpoint{
		service:        opts.Service,
		channelService: opts.ChannelService,
		twirBus:        opts.TwirBus,
		middlewares:    opts.Middlewares,
		sessions:       opts.Sessions,
	}
}

type startRequestDto struct {
	Body struct {
		Time int32   `json:"time" description:"Duration in minutes for the BRB overlay"`
		Text *string `json:"text" description:"Custom text to display on the BRB overlay"`
	}
}

func (se *startEndpoint) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "overlays-brb-start",
		Method:      http.MethodPut,
		Path:        "/v1/channels/overlays/brb/start",
		Tags:        []string{"Overlays/BRB"},
		Summary:     "Start BRB overlay",
		Middlewares: huma.Middlewares{se.middlewares.RateLimit("brb-start", 10, 1*time.Minute)},
		Security: []map[string][]string{
			{"api-key": {}},
		},
	}
}

func (se *startEndpoint) Handler(
	ctx context.Context,
	input *startRequestDto,
) (*httpbase.BaseOutputJson[startResponseDto], error) {
	selectedDashboardId, err := se.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
	}

	if err := se.twirBus.Api.TriggerBrbStart.Publish(
		ctx,
		api.TriggerBrbStart{
			ChannelId: selectedDashboardId,
			Minutes:   input.Body.Time,
			Text:      input.Body.Text,
		},
	); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to start BRB overlay", err)
	}

	return httpbase.CreateBaseOutputJson(startResponseDto{Success: true}), nil
}

func (se *startEndpoint) Register(api huma.API) {
	huma.Register(api, se.GetMeta(), se.Handler)
}
