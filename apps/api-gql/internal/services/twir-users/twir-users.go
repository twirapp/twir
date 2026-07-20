package twir_users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/users_with_channel"
	userswithchannelmodel "github.com/twirapp/twir/libs/repositories/users_with_channel/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	UsersWithChannelsRepository users_with_channel.Repository
}

func New(opts Opts) *Service {
	return &Service{
		usersWithChannelsRepository: opts.UsersWithChannelsRepository,
	}
}

type Service struct {
	usersWithChannelsRepository users_with_channel.Repository
}

func (c *Service) modelToEntity(m userswithchannelmodel.UserWithChannel) entity.UserWithChannel {
	e := entity.UserWithChannel{
		User: entity.User{
			ID:                m.User.ID.String(),
			Platform:          m.User.Platform,
			PlatformID:        m.User.PlatformID,
			TokenID:           m.User.TokenID.Ptr(),
			IsBotAdmin:        m.User.IsBotAdmin,
			ApiKey:            m.User.ApiKey,
			IsBanned:          m.User.IsBanned,
			HideOnLandingPage: m.User.HideOnLandingPage,
			Login:             m.User.Login,
			DisplayName:       m.User.DisplayName,
			Avatar:            m.User.Avatar,
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
	Platforms         []platformentity.Platform
}

type GetManyOutput struct {
	Users []entity.UserWithChannel
	Total int
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) (
	GetManyOutput,
	error,
) {
	var page int
	perPage := 10

	if input.Page != 0 {
		page = input.Page
	}

	if input.PerPage != 0 {
		perPage = input.PerPage
	}

	badgeIDs := make([]uuid.UUID, 0, len(input.HasBadges))
	for _, badgeID := range input.HasBadges {
		parsedBadgeID, err := uuid.Parse(badgeID)
		if err != nil {
			return GetManyOutput{}, fmt.Errorf("parse badge ID: %w", err)
		}

		badgeIDs = append(badgeIDs, parsedBadgeID)
	}

	usersInput := users_with_channel.GetManyInput{
		Page:              page,
		PerPage:           perPage,
		IDs:               nil,
		SearchQuery:       input.SearchQuery,
		Platforms:         input.Platforms,
		ChannelEnabled:    input.ChannelIsEnabled,
		ChannelIsBotAdmin: input.ChannelIsBotAdmin,
		IsBanned:          input.UserIsBanned,
		HasBadgesIDS:      badgeIDs,
	}

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
