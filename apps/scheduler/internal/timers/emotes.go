package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/emotes_cacher"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type EmotesOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Config config.Config

	EmotesService emotes_cacher.EmotesCacherClient
	Gorm          *gorm.DB
}

func NewEmotes(opts EmotesOpts) {
	timeTick := lo.If(opts.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)

	channelsTicker := time.NewTicker(timeTick)
	globalTicker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-ctx.Done():
							channelsTicker.Stop()
							break
						case <-channelsTicker.C:
							var channels []model.Channels
							err := opts.Gorm.
								WithContext(ctx).
								Where(`"isEnabled" = ? and "isBanned" = ?`, true, false).
								Select("isEnabled", "id").
								Find(&channels).
								Error
							if err != nil {
								opts.Logger.Error("error while getting channels", slog.Any("err", err))
							} else {
								for _, channel := range channels {
									_, err = opts.EmotesService.CacheChannelEmotes(
										ctx,
										&emotes_cacher.Request{
											ChannelId: channel.ID,
										},
									)
									if err != nil {
										opts.Logger.Error("error while caching channel emotes", slog.Any("err", err))
									}
								}
							}
						}
					}
				}()

				go func() {
					for {
						select {
						case <-ctx.Done():
							globalTicker.Stop()
							break
						case <-globalTicker.C:
							_, err := opts.EmotesService.CacheGlobalEmotes(ctx, &emptypb.Empty{})
							if err != nil {
								opts.Logger.Error("error while caching global emotes", slog.Any("err", err))
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
