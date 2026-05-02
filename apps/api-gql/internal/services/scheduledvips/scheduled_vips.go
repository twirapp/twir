package scheduledvips

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ScheduledVipsRepository scheduledvipsrepository.Repository
	ChannelsRepository      channelsrepository.Repository
	UsersRepository         usersrepository.Repository
	Config                  config.Config
	Bus                     *buscore.Bus
	Logger                  *slog.Logger
}

func New(opts Opts) *Service {
	return &Service{
		repo:         opts.ScheduledVipsRepository,
		channelsRepo: opts.ChannelsRepository,
		usersRepo:    opts.UsersRepository,
		config:       opts.Config,
		bus:          opts.Bus,
		logger:       opts.Logger,
	}
}

type Service struct {
	repo         scheduledvipsrepository.Repository
	channelsRepo channelsrepository.Repository
	usersRepo    usersrepository.Repository
	config       config.Config
	bus          *buscore.Bus
	logger       *slog.Logger
}

func (c *Service) GetUserByPlatformID(ctx context.Context, platformUserID string) (usersmodel.User, error) {
	return c.usersRepo.GetByPlatformID(ctx, platformentity.PlatformTwitch, platformUserID)
}

func (c *Service) GetScheduledVips(ctx context.Context, channelID string) (
	[]scheduledvipsentity.ScheduledVip,
	error,
) {
	scheduledVips, err := c.repo.GetMany(
		ctx,
		scheduledvipsrepository.GetManyInput{
			ChannelID: &channelID,
		},
	)
	if err != nil {
		return nil, err
	}

	return scheduledVips, nil
}

type CreateInput struct {
	UserID     string
	ChannelID  string
	RemoveAt   *time.Time
	RemoveType *scheduledvipsentity.RemoveType
}

func (c *Service) Create(ctx context.Context, input CreateInput) error {
	return c.repo.Create(
		ctx,
		scheduledvipsrepository.CreateInput{
			ChannelID:  input.ChannelID,
			UserID:     input.UserID,
			RemoveAt:   input.RemoveAt,
			RemoveType: input.RemoveType,
		},
	)
}

type RemoveInput struct {
	ID        string
	ChannelID string
	KeepVip   *bool
}

func (c *Service) Remove(ctx context.Context, input RemoveInput) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return err
	}

	vip, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if vip.ChannelID != input.ChannelID {
		return fmt.Errorf("vip does not belong to the channel")
	}

	parsedChannelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelsRepo.GetByID(ctx, parsedChannelID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() || !channel.TwitchConnected() {
		return fmt.Errorf("channel not found or twitch not connected")
	}

	vipUserUUID, err := uuid.Parse(vip.UserID)
	if err != nil {
		return fmt.Errorf("invalid vip user id: %w", err)
	}

	vipUser, err := c.usersRepo.GetByID(ctx, vipUserUUID)
	if err != nil {
		return fmt.Errorf("get vip user: %w", err)
	}

	twitchClient, err := twitch.NewUserClient(
		*channel.TwitchPlatformID,
		c.config,
		c.bus,
	)
	if err != nil {
		return fmt.Errorf("cannot create twitch client: %w", err)
	}

	err = c.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if input.KeepVip != nil && *input.KeepVip {
		return nil
	}

	vipResp, err := twitchClient.RemoveChannelVip(
		&helix.RemoveChannelVipParams{
			BroadcasterID: *channel.TwitchPlatformID,
			UserID:        vipUser.PlatformID,
		},
	)
	if err != nil {
		c.logger.Error("Cannot remove VIP on Twitch", logger.Error(err))
	}
	if vipResp.ErrorMessage != "" {
		c.logger.Error("Twitch error", logger.Error(errors.New(vipResp.ErrorMessage)))
	}

	return nil
}

func (c *Service) Update(ctx context.Context, id, channelID string, removeAt *time.Time, removeType *scheduledvipsentity.RemoveType) error {
	vipID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	vip, err := c.repo.GetByID(ctx, vipID)
	if err != nil {
		return err
	}
	if vip.ChannelID != channelID {
		return fmt.Errorf("vip does not belong to the channel")
	}

	return c.repo.Update(
		ctx,
		vipID,
		scheduledvipsrepository.UpdateInput{
			RemoveAt:   removeAt,
			RemoveType: removeType,
		},
	)
}

type CreateWithTwitchVipInput struct {
	UserID     string
	ChannelID  string
	RemoveAt   *time.Time
	RemoveType *scheduledvipsentity.RemoveType
}

func (c *Service) CreateWithTwitchVip(ctx context.Context, input CreateWithTwitchVipInput) error {
	parsedChannelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelsRepo.GetByID(ctx, parsedChannelID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() || !channel.TwitchConnected() {
		return fmt.Errorf("channel not found or twitch not connected")
	}

	targetDbUser, err := c.usersRepo.GetByPlatformID(ctx, platformentity.PlatformTwitch, input.UserID)
	if err != nil {
		return fmt.Errorf("cannot get user by platform id: %w", err)
	}

	twitchClient, err := twitch.NewUserClient(
		*channel.TwitchPlatformID,
		c.config,
		c.bus,
	)
	if err != nil {
		return fmt.Errorf("cannot create twitch client: %w", err)
	}

	vipResp, err := twitchClient.AddChannelVip(
		&helix.AddChannelVipParams{
			BroadcasterID: *channel.TwitchPlatformID,
			UserID:        input.UserID,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot add vip on twitch: %w", err)
	}
	if vipResp.ErrorMessage != "" {
		return fmt.Errorf("twitch error: %s", vipResp.ErrorMessage)
	}

	err = c.repo.Create(
		ctx,
		scheduledvipsrepository.CreateInput{
			ChannelID:  input.ChannelID,
			UserID:     targetDbUser.ID.String(),
			RemoveAt:   input.RemoveAt,
			RemoveType: input.RemoveType,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot create scheduled vip in database: %w", err)
	}

	return nil
}
