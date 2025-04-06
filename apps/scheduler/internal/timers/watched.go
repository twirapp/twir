package timers

import (
	"context"
	"log/slog"
	"time"

	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type WatchedOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Config config.Config

	Gorm       *gorm.DB
	TokensGrpc tokens.TokensClient
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
								opts.Logger.Error("cannot get streams", slog.Any("err", err))
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
									opts.Logger.Error("cannot update watched", slog.Any("err", err))
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
