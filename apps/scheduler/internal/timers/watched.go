package timers

import (
	"context"
	"log/slog"
	"time"

	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
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
							var channelIDs []struct {
								ID string `gorm:"column:id"`
							}
							if err := opts.Gorm.WithContext(ctx).Raw(`
								SELECT DISTINCT c.id
								FROM channels_streams cs
								JOIN users u ON u.platform_id = cs."userId" AND u.platform = 'twitch'
								JOIN channels c ON c.twitch_user_id = u.id
							`).Scan(&channelIDs).Error; err != nil {
								opts.Logger.Error("cannot get stream channel IDs", logger.Error(err))
								continue
							}

							for _, ch := range channelIDs {
								if err := opts.Gorm.WithContext(ctx).Exec(`
									UPDATE users_stats
									SET watched = watched + $1
									WHERE "channelId" = $2::uuid
									  AND "userId" IN (
									      SELECT "userId" FROM users_online WHERE "channelId" = $2::uuid
									  )
								`, timeTick.Milliseconds(), ch.ID).Error; err != nil {
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
