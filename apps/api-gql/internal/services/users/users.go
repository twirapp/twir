package users

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	config "github.com/satont/twir/libs/config"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/repositories/users/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	UsersRepository users.Repository
	Gorm            *gorm.DB
	TokensService   tokens.TokensClient
	Config          config.Config
}

func New(opts Opts) *Service {
	return &Service{
		usersRepository: opts.UsersRepository,
		gorm:            opts.Gorm,
		tokensService:   opts.TokensService,
		config:          opts.Config,
	}
}

type Service struct {
	usersRepository users.Repository
	gorm            *gorm.DB
	tokensService   tokens.TokensClient
	config          config.Config
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
		ID:                m.ID,
		TokenID:           m.TokenID.Ptr(),
		IsBotAdmin:        m.IsBotAdmin,
		ApiKey:            m.ApiKey,
		IsBanned:          m.IsBanned,
		HideOnLandingPage: m.HideOnLandingPage,
	}
}

func (c *Service) Update(ctx context.Context, id string, input UpdateInput) (entity.User, error) {
	newUser, err := c.usersRepository.Update(
		ctx,
		id,
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

	return c.modelToEntity(newUser), nil
}

func (c *Service) GetByID(ctx context.Context, id string) (entity.User, error) {
	user, err := c.usersRepository.GetByID(ctx, id)
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
	dbUsers, err := c.usersRepository.GetManyByIDS(
		ctx,
		users.GetManyInput{
			Page:       input.Page,
			PerPage:    input.PerPage,
			IDs:        input.IDs,
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

	dbUserInfo := deprecatedgormmodel.Users{}
	if err := c.gorm.
		WithContext(ctx).
		Where("id = ?", input.UserID).
		Preload("Stats", `"channelId" = ? AND "userId" = ?`, input.ChannelID, input.UserID).
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
		input.ChannelID,
		c.config,
		c.tokensService,
	)
	if err != nil {
		return entity.ChannelUserInfo{}, fmt.Errorf("cannot create channel twitch client: %w", err)
	}

	follows, err := channelTwitchClient.GetChannelFollows(
		&helix.GetChannelFollowsParams{
			BroadcasterID: input.ChannelID,
			UserID:        input.UserID,
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
