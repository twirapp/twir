package v2

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/entities/platform"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Huma            huma.API
	CachedCommands  *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	ChannelsService *channels.Service
}

type Public struct {
	cachedCommands  *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	channelsService *channels.Service
}

func New(opts Opts) *Public {
	p := &Public{
		cachedCommands:  opts.CachedCommands,
		channelsService: opts.ChannelsService,
	}

	huma.Register(
		opts.Huma,
		huma.Operation{
			OperationID: "public-v2-channel-commands-by-platform-id",
			Method:      http.MethodGet,
			Path:        "/v2/public/channels/{platform}/{channelId}/commands",
			Summary:     "Get channel commands by platform platform id",
			Description: "Get channel commands filtered by enabled and visible, looked up by Twitch platform id",
			Tags:        []string{"Public V2"},
		},
		func(
			ctx context.Context,
			input *struct {
				ChannelId string            `path:"channelId" minLength:"1" required:"true"`
				Platform  platform.Platform `path:"platform" required:"true"`
			},
		) (*publicCommandsOutput, error) {
			return p.handleGetChannelByPlatformCommandsGet(ctx, input.Platform, input.ChannelId)
		},
	)

	huma.Register(
		opts.Huma,
		huma.Operation{
			OperationID: "public-v2-channel-commands-by-uuid",
			Method:      http.MethodGet,
			Path:        "/v2/channels/{channelUuid}/commands",
			Summary:     "Get channel commands by internal channel uuid",
			Description: "Get channel commands filtered by enabled and visible, looked up by internal channel uuid",
			Tags:        []string{"Public V2"},
		},
		func(
			ctx context.Context,
			input *struct {
				ChannelUuid uuid.UUID `path:"channelUuid" format:"uuid" required:"true"`
			},
		) (*publicCommandsOutput, error) {
			return p.getChannelCommands(ctx, input.ChannelUuid)
		},
	)

	return p
}
