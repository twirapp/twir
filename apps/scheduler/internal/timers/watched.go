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

const watchedChannelIDsQuery = `
	SELECT DISTINCT cp.channel_id AS id
	FROM channels_streams cs
	JOIN channel_platforms cp
		ON cp.platform = 'twitch'
		AND cp.platform_channel_id = cs."userId"
	WHERE cs.platform = 'twitch'
`

func buildWatchedChannelIDsQuery(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.WithContext(ctx).Raw(watchedChannelIDsQuery)
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
							if err := buildWatchedChannelIDsQuery(opts.Gorm, ctx).Scan(&channelIDs).Error; err != nil {
								opts.Logger.Error("cannot get stream channel IDs", logger.Error(err))
								continue
							}

							for _, ch := range channelIDs {
								if err := opts.Gorm.WithContext(ctx).Exec(`
									UPDATE users_stats
									SET watched = watched + $1
									WHERE channel_id = $2::uuid
									  AND user_id IN (
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
