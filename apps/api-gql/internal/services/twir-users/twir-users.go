package twir_users

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	twitch_channels "github.com/twirapp/twir/apps/api-gql/internal/services/twitch-channels"
	users_with_channel "github.com/twirapp/twir/libs/repositories/users-with-channel"
	userswithchannelmodel "github.com/twirapp/twir/libs/repositories/users-with-channel/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TwitchChannelsService       *twitch_channels.Service
	UsersWithChannelsRepository users_with_channel.Repository
}

func New(opts Opts) *Service {
	return &Service{
		twitchChannelsService:       opts.TwitchChannelsService,
		usersWithChannelsRepository: opts.UsersWithChannelsRepository,
	}
}

type Service struct {
	twitchChannelsService       *twitch_channels.Service
	usersWithChannelsRepository users_with_channel.Repository
}

func (c *Service) modelToEntity(m userswithchannelmodel.UserWithChannel) entity.UserWithChannel {
	e := entity.UserWithChannel{
		User: entity.User{
			ID:                m.User.ID,
			TokenID:           m.User.TokenID.Ptr(),
			IsTester:          m.User.IsTester,
			IsBotAdmin:        m.User.IsBotAdmin,
			ApiKey:            m.User.ApiKey,
			IsBanned:          m.User.IsBanned,
			HideOnLandingPage: m.User.HideOnLandingPage,
		},
	}

	if m.Channel != nil {
		e.Channel = &entity.Channel{
			ID:             m.Channel.ID,
			IsEnabled:      m.Channel.IsEnabled,
			IsTwitchBanned: m.Channel.IsTwitchBanned,
			IsBotMod:       m.Channel.IsBotMod,
			BotID:          m.Channel.BotID,
		}
	}

	return e
}

type GetManyInput struct {
	SearchQuery       string
	Page              int
	PerPage           int
	ChannelIsEnabled  *bool
	ChannelIsBotAdmin *bool
	UserIsBanned      *bool
	HasBadges         []string
}

type GetManyOutput struct {
	Users []entity.UserWithChannel
	Total int
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) (
	GetManyOutput,
	error,
) {
	twitchUsers, err := c.twitchChannelsService.SearchByName(ctx, input.SearchQuery)
	if err != nil {
		return GetManyOutput{}, err
	}

	var page int
	perPage := 10

	if input.Page != 0 {
		page = input.Page
	}

	if input.PerPage != 0 {
		perPage = input.PerPage
	}

	usersInput := users_with_channel.GetManyInput{
		Page:              page,
		PerPage:           perPage,
		IDs:               nil,
		ChannelEnabled:    input.ChannelIsEnabled,
		ChannelIsBotAdmin: input.ChannelIsBotAdmin,
		IsBanned:          input.UserIsBanned,
		HasBadgesIDS:      input.HasBadges,
	}

	twitchUserIDs := make([]string, 0, len(twitchUsers))
	for _, user := range twitchUsers {
		twitchUserIDs = append(twitchUserIDs, user.ID)
	}
	usersInput.IDs = twitchUserIDs

	dbUsers, err := c.usersWithChannelsRepository.GetManyByIDS(ctx, usersInput)
	if err != nil {
		return GetManyOutput{}, err
	}

	entities := make([]entity.UserWithChannel, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		entities = append(entities, c.modelToEntity(dbUser))
	}

	total, err := c.usersWithChannelsRepository.GetManyCount(ctx, usersInput)
	if err != nil {
		return GetManyOutput{}, err
	}

	return GetManyOutput{
		Users: entities,
		Total: total,
	}, nil
}
