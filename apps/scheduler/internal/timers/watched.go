package timers

import (
	"context"
	"log/slog"
	"time"

	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type WatchedOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger *slog.Logger
	Config config.Config

	Gorm    *gorm.DB
	TwirBus *buscore.Bus
}

func NewWatched(opts WatchedOpts) {
	timeTick := 15 * time.Second
	if opts.Config.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

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
							var streams []model.ChannelsStreams
							if err := opts.Gorm.WithContext(ctx).Select(`"userId"`).Find(&streams).Error; err != nil {
								opts.Logger.Error("cannot get streams", logger.Error(err))
								continue
							}

							for _, s := range streams {
								err := opts.Gorm.WithContext(ctx).Model(&model.UsersStats{}).
									WithContext(ctx).
									Where(
										`"channelId" = ? AND "userId" IN (?)`,
										s.UserId,
										opts.Gorm.Table("users_online").Where(
											`"channelId" = ?`,
											s.UserId,
										).Select(`"userId"`),
									).
									Update("watched", gorm.Expr("watched + ?", timeTick.Milliseconds())).Error
								if err != nil {
									opts.Logger.Error("cannot update watched", logger.Error(err))
								}
							}
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
