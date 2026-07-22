package twir_users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apiChannelbinding "github.com/twirapp/twir/apps/api-gql/internal/channelbinding"
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

func (c *Service) modelToEntity(m userswithchannelmodel.UserWithChannel) (entity.UserWithChannel, error) {
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

	if m.Channel == nil {
		return e, nil
	}

	e.Channel = &entity.Channel{ID: m.Channel.ID}
	binding, found := apiChannelbinding.Find(*m.Channel, m.User.Platform)
	if !found {
		return e, nil
	}

	e.Channel.IsEnabled = binding.Enabled
	if binding.Platform != platformentity.PlatformTwitch {
		if binding.BotUserID != nil {
			e.Channel.BotID = binding.BotUserID.String()
		}

		return e, nil
	}

	_, botConfig, _, err := apiChannelbinding.FindTwitch(*m.Channel)
	if err != nil {
		return entity.UserWithChannel{}, fmt.Errorf("parse Twitch channel bot config: %w", err)
	}

	e.Channel.IsTwitchBanned = botConfig.IsTwitchBanned
	e.Channel.IsBotMod = botConfig.IsBotMod
	e.Channel.BotID = botConfig.BotID

	return e, nil
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
		user, err := c.modelToEntity(dbUser)
		if err != nil {
			return GetManyOutput{}, fmt.Errorf("map user with channel: %w", err)
		}

		entities = append(entities, user)
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
