package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	config "github.com/twirapp/twir/libs/config"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/repositories/users/model"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	UsersRepository    users.Repository
	ChannelsRepository channelsrepository.Repository
	Gorm               *gorm.DB
	Config             config.Config
	TwirBus            *buscore.Bus
}

func New(opts Opts) *Service {
	return &Service{
		usersRepository:    opts.UsersRepository,
		channelsRepository: opts.ChannelsRepository,
		gorm:               opts.Gorm,
		config:             opts.Config,
		twirBus:            opts.TwirBus,
	}
}

type Service struct {
	usersRepository    users.Repository
	channelsRepository channelsrepository.Repository
	gorm               *gorm.DB
	config             config.Config
	twirBus            *buscore.Bus
}

type UpdateInput struct {
	IsBotAdmin        *bool
	ApiKey            *string
	IsBanned          *bool
	HideOnLandingPage *bool
	TokenID           *string
}

func (c *Service) modelToEntity(m model.User) entity.User {
	return entity.User{
		ID:                m.ID.String(),
		TokenID:           m.TokenID.Ptr(),
		IsBotAdmin:        m.IsBotAdmin,
		ApiKey:            m.ApiKey,
		IsBanned:          m.IsBanned,
		HideOnLandingPage: m.HideOnLandingPage,
	}
}

func (c *Service) Update(ctx context.Context, id string, input UpdateInput) (entity.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return entity.UserNil, fmt.Errorf("parse user id: %w", err)
	}

	newUser, err := c.usersRepository.Update(
		ctx,
		parsedID,
		users.UpdateInput{
			IsBanned:          input.IsBanned,
			IsBotAdmin:        input.IsBotAdmin,
			ApiKey:            input.ApiKey,
			HideOnLandingPage: input.HideOnLandingPage,
			TokenID:           input.TokenID,
		},
	)
	if err != nil {
		return entity.UserNil, err
	}

	if input.IsBanned != nil {
		if *input.IsBanned {
			c.twirBus.EventSub.Unsubscribe.Publish(ctx, id)
		} else {
			c.twirBus.EventSub.SubscribeToAllEvents.Publish(
				ctx,
				eventsub.EventsubSubscribeToAllEventsRequest{ChannelID: id},
			)
		}
	}

	return c.modelToEntity(newUser), nil
}

func (c *Service) GetByID(ctx context.Context, id string) (entity.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return entity.UserNil, fmt.Errorf("parse user id: %w", err)
	}

	user, err := c.usersRepository.GetByID(ctx, parsedID)
	if err != nil {
		return entity.UserNil, err
	}

	return c.modelToEntity(user), nil
}

type GetManyInput struct {
	Page       int
	PerPage    int
	IDs        []string
	IsBotAdmin *bool
	IsBanned   *bool
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) ([]entity.User, error) {
	parsedIDs := make([]uuid.UUID, 0, len(input.IDs))
	for _, id := range input.IDs {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("parse user id: %w", err)
		}

		parsedIDs = append(parsedIDs, parsedID)
	}

	dbUsers, err := c.usersRepository.GetManyByIDS(
		ctx,
		users.GetManyInput{
			Page:       input.Page,
			PerPage:    input.PerPage,
			IDs:        parsedIDs,
			IsBotAdmin: input.IsBotAdmin,
			IsBanned:   input.IsBanned,
		},
	)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.User, 0, len(dbUsers))
	for _, user := range dbUsers {
		entities = append(entities, c.modelToEntity(user))
	}

	return entities, nil
}

type ChannelUserInfoInput struct {
	ChannelID string
	UserID    string
}

func (c *Service) GetChannelUserInfo(ctx context.Context, input ChannelUserInfoInput) (
	entity.ChannelUserInfo,
	error,
) {
	if input.ChannelID == "" || input.UserID == "" {
		return entity.ChannelUserInfo{}, fmt.Errorf("channel_id and user_id are required")
	}

	parsedChannelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return entity.ChannelUserInfo{}, fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelsRepository.GetByID(ctx, parsedChannelID)
	if err != nil {
		return entity.ChannelUserInfo{}, fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() || !channel.TwitchConnected() {
		return entity.ChannelUserInfo{}, fmt.Errorf("channel not found or twitch not connected")
	}

	parsedUserID, err := uuid.Parse(input.UserID)
	if err != nil {
		return entity.ChannelUserInfo{}, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := c.usersRepository.GetByID(ctx, parsedUserID)
	if err != nil {
		return entity.ChannelUserInfo{}, fmt.Errorf("get user: %w", err)
	}

	dbUserInfo := deprecatedgormmodel.Users{}
	if err := c.gorm.
		WithContext(ctx).
		Where("id = ?::uuid", input.UserID).
		Preload("Stats", `"channelId" = ?::uuid AND "userId" = ?::uuid`, input.ChannelID, input.UserID).
		First(&dbUserInfo).
		Error; err != nil {
		return entity.ChannelUserInfo{}, err
	}

	info := entity.ChannelUserInfo{
		ID:                dbUserInfo.ID,
		Messages:          0,
		Watched:           0,
		UsedEmotes:        0,
		UsedChannelPoints: 0,
		IsMod:             false,
		IsVip:             false,
		IsSubscriber:      false,
		FollowerSince:     nil,
	}

	if dbUserInfo.Stats != nil {
		info.Messages = int(dbUserInfo.Stats.Messages)
		info.Watched = int(dbUserInfo.Stats.Watched)
		info.UsedEmotes = dbUserInfo.Stats.Emotes
		info.UsedChannelPoints = int(dbUserInfo.Stats.UsedChannelPoints)
		info.IsMod = dbUserInfo.Stats.IsMod
		info.IsVip = dbUserInfo.Stats.IsVip
		info.IsSubscriber = dbUserInfo.Stats.IsSubscriber
	}

	channelTwitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		*channel.TwitchPlatformID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return entity.ChannelUserInfo{}, fmt.Errorf("cannot create channel twitch client: %w", err)
	}

	follows, err := channelTwitchClient.GetChannelFollows(
		&helix.GetChannelFollowsParams{
			BroadcasterID: *channel.TwitchPlatformID,
			UserID:        user.PlatformID,
		},
	)
	if err != nil {
		return entity.ChannelUserInfo{}, fmt.Errorf("cannot get channel follows: %w", err)
	}
	if follows.ErrorMessage != "" {
		return entity.ChannelUserInfo{}, fmt.Errorf(
			"cannot get channel follows: %s",
			follows.ErrorMessage,
		)
	}

	if len(follows.Data.Channels) != 0 {
		info.FollowerSince = &follows.Data.Channels[0].Followed.Time
	}

	return info, nil
}

func (c *Service) GetByApiKey(ctx context.Context, apiKey string) (entity.User, error) {
	user, err := c.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return entity.User{}, err
	}

	return c.modelToEntity(user), nil
}
