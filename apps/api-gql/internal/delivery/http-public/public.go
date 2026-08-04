package http_public

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"gorm.io/gorm"
)

type Public struct {
	gorm                   *gorm.DB
	cachedCommands         *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	badgesWithUsersService *badges_with_users.Service
	channelsService        *channels.Service
	config                 config.Config
}

func New(humaAPI huma.API, gorm *gorm.DB, cachedCommands *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses], badgesWithUsersService *badges_with_users.Service, channelsService *channels.Service, config config.Config) *Public {
	p := &Public{
		gorm:                   gorm,
		config:                 config,
		cachedCommands:         cachedCommands,
		badgesWithUsersService: badgesWithUsersService,
		channelsService:        channelsService,
	}

	huma.Register(
		humaAPI,
		huma.Operation{
			OperationID: "public-twir-badges",
			Method:      http.MethodGet,
			Path:        "/v1/public/badges",
			Summary:     "Get badges",
			Description: "Get created badges for twitch chat",
			Tags:        []string{"Public"},
		},
		func(
			ctx context.Context,
			_ *struct{},
		) (*badgesOutput, error) {
			return p.HandleBadgesGet(ctx)
		},
	)

	huma.Register(
		humaAPI,
		huma.Operation{
			OperationID: "public-channel-public-commands",
			Method:      http.MethodGet,
			Path:        "/v1/public/channels/{channelId}/commands",
			Summary:     "Get channel commands",
			Description: "Get channel commands filtered by enabled and visible",
			Tags:        []string{"Public"},
		},
		func(
			ctx context.Context,
			input *struct {
				ChannelId string `path:"channelId" maxLength:"36" minLength:"1" required:"true"`
			},
		) (*publicCommandsOutput, error) {
			return p.HandleChannelCommandsGet(ctx, input.ChannelId)
		},
	)

	return p
}
