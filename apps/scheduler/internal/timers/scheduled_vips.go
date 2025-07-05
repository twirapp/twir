package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
	buscore "github.com/twirapp/twir/libs/bus-core"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	"go.uber.org/fx"
)

type ScheduledVipsOpts struct {
	fx.In
	LC fx.Lifecycle

	Config  config.Config
	Logger  logger.Logger
	TwirBus *buscore.Bus

	ScheduledVipsRepo scheduledvipsrepository.Repository
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
	logger            logger.Logger
	twirBus           *buscore.Bus
	scheduledVipsRepo scheduledvipsrepository.Repository
}

func (s *scheduledVips) process(ctx context.Context) {
	vips, err := s.scheduledVipsRepo.GetMany(
		ctx, scheduledvipsrepository.GetManyInput{
			Expired: lo.ToPtr(true),
		},
	)
	if err != nil {
		s.logger.Error("failed to get scheduled vips", slog.Any("err", err))
		return
	}

	cachedChannelsTwitchClients := make(map[string]*helix.Client)
	// create twitch clients for channels
	for _, vip := range vips {
		if cachedChannelsTwitchClients[vip.ChannelID] == nil {
			twitchClient, err := twitch.NewUserClientWithContext(
				ctx,
				vip.ChannelID,
				s.config,
				s.twirBus,
			)
			if err != nil {
				s.logger.Error("failed to create twitch client", slog.Any("err", err))
				continue
			}

			cachedChannelsTwitchClients[vip.ChannelID] = twitchClient
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
				BroadcasterID: vip.ChannelID,
				UserID:        vip.UserID,
			},
		)
		if err != nil {
			s.logger.Error("failed to remove vip", slog.Any("err", err))
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
