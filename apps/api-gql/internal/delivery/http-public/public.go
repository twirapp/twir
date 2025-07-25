package http_public

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	config "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Huma huma.API

	Gorm                   *gorm.DB
	CachedCommands         *generic_cacher.GenericCacher[[]model.ChannelsCommands]
	BadgesWithUsersService *badges_with_users.Service
	ChannelsService        *channels.Service
	Config                 config.Config
}

type Public struct {
	gorm                   *gorm.DB
	cachedCommands         *generic_cacher.GenericCacher[[]model.ChannelsCommands]
	badgesWithUsersService *badges_with_users.Service
	channelsService        *channels.Service
	config                 config.Config
}

func New(opts Opts) *Public {
	p := &Public{
		gorm:                   opts.Gorm,
		config:                 opts.Config,
		cachedCommands:         opts.CachedCommands,
		badgesWithUsersService: opts.BadgesWithUsersService,
		channelsService:        opts.ChannelsService,
	}

	huma.Register(
		opts.Huma,
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
		opts.Huma,
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
				ChannelId string `path:"channelId" maxLength:"36" minLength:"1" pattern:"^[0-9]+$" required:"true"`
			},
		) (*publicCommandsOutput, error) {
			return p.HandleChannelCommandsGet(ctx, input.ChannelId)
		},
	)

	return p
}
