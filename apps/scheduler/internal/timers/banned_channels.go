package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type BannedChannelsOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Config config.Config

	TokensGrpc tokens.TokensClient
	Gorm       *gorm.DB
}

type bannedChannels struct {
	cfg    config.Config
	tokens tokens.TokensClient
	db     *gorm.DB
	logger logger.Logger
}

func NewBannedChannels(opts BannedChannelsOpts) {
	timeTick := lo.If(opts.Config.AppEnv != "production", 15*time.Second).Else(30 * time.Minute)
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &bannedChannels{
		cfg:    opts.Config,
		tokens: opts.TokensGrpc,
		db:     opts.Gorm,
		logger: opts.Logger,
	}

	opts.Lc.Append(
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

func (c *bannedChannels) process(ctx context.Context) {
	var channels []model.Channels
	if err := c.db.WithContext(ctx).Find(&channels).Error; err != nil {
		c.logger.Error("failed to get channels", slog.Any("err", err))
		return
	}

	chunks := lo.Chunk(channels, 100)

	wg := utils.NewGoroutinesGroup()

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		c.cfg,
		c.tokens,
	)
	if err != nil {
		c.logger.Error("failed to create twitch client", slog.Any("err", err))
		return
	}

	for _, chunk := range chunks {
		chunk := chunk

		wg.Go(
			func() {
				channelsReq, err := twitchClient.GetUsers(
					&helix.UsersParams{
						IDs: lo.Map(chunk, func(channel model.Channels, _ int) string { return channel.ID }),
					},
				)
				if err != nil {
					c.logger.Error("failed to get users", slog.Any("err", err))
					return
				}
				if channelsReq.ErrorMessage != "" {
					c.logger.Error("failed to get users", slog.Any("err", channelsReq.ErrorMessage))
					return
				}

				for _, channel := range chunk {
					exists := lo.SomeBy(
						channelsReq.Data.Users, func(user helix.User) bool {
							return user.ID == channel.ID
						},
					)

					err := c.db.WithContext(ctx).Model(&channel).Update(
						`"isTwitchBanned"`,
						!exists,
					).Error
					if err != nil {
						c.logger.Error("failed to update channel", slog.Any("err", err))
						continue
					}
				}
			},
		)
	}

	wg.Wait()
}
