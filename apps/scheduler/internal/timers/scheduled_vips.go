package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
)

type ScheduledVipsOpts struct {
	fx.In
	LC fx.Lifecycle

	Config  config.Config
	Logger  *slog.Logger
	TwirBus *buscore.Bus

	ScheduledVipsRepo scheduledvipsrepository.Repository
	UsersRepo         usersrepository.Repository
	ChannelsRepo      channelsrepository.Repository
}

func NewScheduledVips(opts ScheduledVipsOpts) {
	timeTick := 15 * time.Second
	if opts.Config.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &scheduledVips{
		config:            opts.Config,
		logger:            opts.Logger,
		twirBus:           opts.TwirBus,
		scheduledVipsRepo: opts.ScheduledVipsRepo,
		usersRepo:         opts.UsersRepo,
		channelsRepo:      opts.ChannelsRepo,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-ctx.Done():
							ticker.Stop()
							return
						case <-ticker.C:
							s.process(ctx)
						}
					}
				}()

				return nil
			},
			OnStop: func(_ context.Context) error {
				cancel()
				return nil
			},
		},
	)
}

type scheduledVips struct {
	config            config.Config
	logger            *slog.Logger
	twirBus           *buscore.Bus
	scheduledVipsRepo scheduledvipsrepository.Repository
	usersRepo         usersrepository.Repository
	channelsRepo      channelsrepository.Repository
}

func (s *scheduledVips) process(ctx context.Context) {
	vips, err := s.scheduledVipsRepo.GetMany(
		ctx, scheduledvipsrepository.GetManyInput{
			Expired: lo.ToPtr(true),
		},
	)
	if err != nil {
		s.logger.Error("failed to get scheduled vips", logger.Error(err))
		return
	}

	cachedChannelsTwitchClients := make(map[string]*helix.Client)
	channelPlatformIDs := make(map[string]string)
	userPlatformIDs := make(map[string]string)
	// create twitch clients for channels
	for _, vip := range vips {
		if cachedChannelsTwitchClients[vip.ChannelID] == nil {
			channelUUID, err := uuid.Parse(vip.ChannelID)
			if err != nil {
				s.logger.Error("failed to parse channel id", slog.String("channel_id", vip.ChannelID), logger.Error(err))
				continue
			}

			channel, err := s.channelsRepo.GetByID(ctx, channelUUID)
			if err != nil {
				s.logger.Error(
					"failed to get channel by id",
					slog.String("channel_id", vip.ChannelID),
					logger.Error(err),
				)
				continue
			}
			if channel.TwitchPlatformID == nil {
				s.logger.Warn("channel has no twitch platform id", slog.String("channel_id", vip.ChannelID))
				continue
			}

			channelPlatformIDs[vip.ChannelID] = *channel.TwitchPlatformID

			ownerUser, err := s.usersRepo.GetByID(ctx, *channel.TwitchUserID)
			if err != nil {
				s.logger.Error("failed to get owner user by id", slog.String("channel_id", vip.ChannelID), logger.Error(err))
				continue
			}

			twitchClient, err := twitch.NewUserClientWithContext(
				ctx,
				ownerUser.ID.String(),
				s.config,
				s.twirBus,
			)
			if err != nil {
				s.logger.Error("failed to create twitch client", logger.Error(err))
				continue
			}

			cachedChannelsTwitchClients[vip.ChannelID] = twitchClient
		}

		if _, ok := userPlatformIDs[vip.UserID]; !ok {
			userUUID, err := uuid.Parse(vip.UserID)
			if err != nil {
				s.logger.Error("failed to parse vip user id", slog.String("user_id", vip.UserID), logger.Error(err))
				continue
			}

			user, err := s.usersRepo.GetByID(ctx, userUUID)
			if err != nil {
				s.logger.Error("failed to get vip user by id", slog.String("user_id", vip.UserID), logger.Error(err))
				continue
			}

			userPlatformIDs[vip.UserID] = user.PlatformID
		}
	}

	// remove vips from channel
	// we'll delete row from db in eventsub service
	for _, vip := range vips {
		twitchClient, ok := cachedChannelsTwitchClients[vip.ChannelID]
		if !ok {
			s.logger.Warn("Twitch client not found", slog.String("channel_id", vip.ChannelID))
			continue
		}

			resp, err := twitchClient.RemoveChannelVip(
				&helix.RemoveChannelVipParams{
					BroadcasterID: channelPlatformIDs[vip.ChannelID],
					UserID:        userPlatformIDs[vip.UserID],
				},
			)
		if err != nil {
			s.logger.Error("failed to remove vip", logger.Error(err))
			continue
		}
		if resp.ErrorMessage != "" {
			s.logger.Error("failed to remove vip", slog.String("error", resp.ErrorMessage))
			continue
		}

		s.logger.Info(
			"vip removed",
			slog.String("user_id", vip.UserID),
			slog.String("channel_id", vip.ChannelID),
		)
	}
}
