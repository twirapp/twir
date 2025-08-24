package stream

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/libs/repositories/streams"
	streammodel "github.com/twirapp/twir/libs/repositories/streams/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	StreamsRepository streams.Repository
	Api               huma.API
	Sessions          *auth.Auth
}

func New(opts Opts) {
	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "channels-streams-current",
			Method:      http.MethodGet,
			Path:        "/v1/channels/streams/current",
			Tags:        []string{"Streams"},
			Summary:     "Get current stream",
			Description: "Get current stream",
			Responses: map[string]*huma.Response{
				"404": {
					Description: "No current stream",
				},
			},
			Security: []map[string][]string{
				{"api-key": {}},
			},
		},
		func(ctx context.Context, i *struct{}) (
			*httpdelivery.BaseOutputJson[streammodel.Stream],
			error,
		) {
			user, err := opts.Sessions.GetAuthenticatedUser(ctx)
			if user == nil || err != nil {
				return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
			}

			selectedDashboardID, err := opts.Sessions.GetSelectedDashboard(ctx)
			if err != nil {
				return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
			}

			stream, err := opts.StreamsRepository.GetByChannelID(ctx, selectedDashboardID)
			if err != nil {
				return nil, huma.NewError(http.StatusInternalServerError, "Cannot get stream", err)
			}

			if stream.IsNil() {
				return nil, huma.NewError(http.StatusNotFound, "No current stream")
			}

			return httpdelivery.CreateBaseOutputJson(stream), nil
		},
	)
}
