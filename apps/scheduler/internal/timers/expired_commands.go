package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ExpiredCommandsOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Config config.Config

	Gorm        *gorm.DB
	RedisClient *redis.Client
}

type expiredCommands struct {
	config        config.Config
	logger        logger.Logger
	db            *gorm.DB
	commandsCache *generic_cacher.GenericCacher[[]model.ChannelsCommands]
}

func NewExpiredCommands(opts ExpiredCommandsOpts) *expiredCommands {
	timeTick := lo.If(opts.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &expiredCommands{
		config:        opts.Config,
		logger:        opts.Logger,
		db:            opts.Gorm,
		commandsCache: commandscache.New(opts.Gorm, opts.RedisClient),
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
							s.checkForExpiredCommands(ctx)
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

	return s
}

func (s *expiredCommands) checkForExpiredCommands(ctx context.Context) {
	var commands []model.ChannelsCommands
	if err := s.db.WithContext(ctx).Preload("Channel").Where(
		`"expires_at" < ? AND "expired" = ?`,
		time.Now().UTC(),
		false,
	).Find(&commands).Error; err != nil {
		s.logger.Error("failed to get commands", slog.Any("err", err))
		return
	}

	for _, c := range commands {
		s.logger.Info("Command expired", slog.Any("command", c))
		c.Expired = true

		err := s.db.WithContext(ctx).Updates(
			&c,
		).Error
		if err != nil {
			s.logger.Error("failed to update command", slog.Any("err", err))
		}

		err = s.commandsCache.Invalidate(
			ctx,
			c.Channel.ID)
		if err != nil {
			s.logger.Error("failed to invalidate commands cache", slog.Any("err", err))
			return
		}

	}
}
