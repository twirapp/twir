package timers

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	bustwitch "github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

type StreamOpts struct {
	fx.In
	Lc fx.Lifecycle

	Config config.Config
	Logger logger.Logger

	Gorm       *gorm.DB
	TokensGrpc tokens.TokensClient
	Bus        *buscore.Bus
}

type streams struct {
	config     config.Config
	logger     logger.Logger
	gorm       *gorm.DB
	tokensGrpc tokens.TokensClient
	bus        *buscore.Bus
}

func NewStreams(opts StreamOpts) {
	timeTick := lo.If(opts.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &streams{
		config:     opts.Config,
		logger:     opts.Logger,
		gorm:       opts.Gorm,
		tokensGrpc: opts.TokensGrpc,
		bus:        opts.Bus,
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
							if err := s.processStreams(ctx); err != nil {
								opts.Logger.Error("cannot process streams", slog.Any("err", err))
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

func (c *streams) processStreams(ctx context.Context) error {
	var channels []model.Channels
	err := c.gorm.
		WithContext(ctx).
		Where(`"channels"."isEnabled" = ? and "User"."is_banned" = ?`, true, false).
		Joins("User").
		Find(&channels).Error
	if err != nil {
		return fmt.Errorf("cannot get channels: %w", err)
	}

	usersIds := make([]string, len(channels))
	for i, channel := range channels {
		if !channel.IsEnabled && !channel.User.IsBanned {
			continue
		}

		usersIds[i] = channel.ID
	}

	discordIntegration := &model.Integrations{}
	err = c.gorm.
		WithContext(ctx).
		Where(`service = ?`, model.IntegrationServiceDiscord).
		Select("id").
		Find(discordIntegration).Error
	if err != nil {
		return fmt.Errorf("cannot get discord integration: %w", err)
	}

	var discordIntegrations []model.ChannelsIntegrations
	if discordIntegration.ID != "" {
		err = c.gorm.
			WithContext(ctx).
			Where(`"integrationId" = ?`, discordIntegration.ID).
			Select("id", `"integrationId"`, "data").
			Find(&discordIntegrations).Error
		if err != nil {
			return fmt.Errorf("cannot get discord integrations: %w", err)
		}

		for _, integration := range discordIntegrations {
			if integration.Data == nil ||
				integration.Data.Discord == nil ||
				len(integration.Data.Discord.Guilds) == 0 {
				continue
			}

			for _, guild := range integration.Data.Discord.Guilds {
				if !guild.LiveNotificationEnabled {
					continue
				}
				usersIds = append(usersIds, guild.AdditionalUsersIdsForLiveCheck...)
			}
		}
	}

	usersIds = lo.Uniq(usersIds)

	var existedStreams []model.ChannelsStreams
	err = c.gorm.WithContext(ctx).Select(
		"id",
		`"userId"`,
		`"parsedMessages"`,
	).Find(&existedStreams).Error
	if err != nil {
		return fmt.Errorf("cannot get existed streams: %w", err)
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.tokensGrpc)
	if err != nil {
		return fmt.Errorf("cannot create twitch client: %w", err)
	}

	chunks := lo.Chunk(usersIds, 100)
	wg := &sync.WaitGroup{}
	wg.Add(len(chunks))

	for _, chunk := range chunks {
		go func(chunk []string) {
			defer wg.Done()
			streams, err := twitchClient.GetStreams(
				&helix.StreamsParams{
					UserIDs: chunk,
				},
			)

			if err != nil || streams.ErrorMessage != "" {
				c.logger.Error("cannot get streams", slog.Any("err", err))
				return
			}

			for _, userId := range chunk {
				twitchStream, twitchStreamExists := lo.Find(
					streams.Data.Streams, func(stream helix.Stream) bool {
						return stream.UserID == userId
					},
				)
				dbStream, dbStreamExists := lo.Find(
					existedStreams, func(stream model.ChannelsStreams) bool {
						return stream.UserId == userId
					},
				)

				tags := &pq.StringArray{}
				for _, tag := range twitchStream.Tags {
					*tags = append(*tags, tag)
				}

				channelStream := &model.ChannelsStreams{
					ID:             twitchStream.ID,
					UserId:         userId,
					UserLogin:      twitchStream.UserLogin,
					UserName:       twitchStream.UserName,
					GameId:         twitchStream.GameID,
					GameName:       twitchStream.GameName,
					CommunityIds:   nil,
					Type:           twitchStream.Type,
					Title:          twitchStream.Title,
					ViewerCount:    twitchStream.ViewerCount,
					StartedAt:      twitchStream.StartedAt,
					Language:       twitchStream.Language,
					ThumbnailUrl:   twitchStream.ThumbnailURL,
					TagIds:         nil,
					Tags:           tags,
					IsMature:       twitchStream.IsMature,
					ParsedMessages: dbStream.ParsedMessages,
				}

				if twitchStreamExists && dbStreamExists {
					// stream still online, update
					if result := c.gorm.WithContext(ctx).Where(
						`"userId" = ?`,
						userId,
					).Save(channelStream); result.Error != nil {
						c.logger.Error("cannot update stream", slog.Any("err", result.Error))
						return
					}
				}

				if twitchStreamExists && !dbStreamExists {
					// stream online, create
					if result := c.gorm.WithContext(ctx).Where(
						`"userId" = ?`,
						userId,
					).Save(channelStream); result.Error != nil {
						c.logger.Error("cannot create stream", slog.Any("err", result.Error))
						return
					}

					c.bus.Channel.StreamOnline.Publish(
						bustwitch.StreamOnlineMessage{
							ChannelID: channelStream.UserId,
							StreamID:  channelStream.ID,
						},
					)
				}

				if !twitchStreamExists && dbStreamExists {
					// stream offline, delete
					err = c.gorm.WithContext(ctx).Where(
						`"userId" = ?`,
						userId,
					).Delete(&model.ChannelsStreams{}).Error
					if err != nil {
						c.logger.Error("cannot delete stream", slog.Any("err", err))
						return
					}

					c.bus.Channel.StreamOffline.Publish(
						bustwitch.StreamOfflineMessage{
							ChannelID: channelStream.UserId,
						},
					)
				}
			}
		}(chunk)
	}

	wg.Wait()

	return nil
}

// { streamId: stream.id, channelId: channel }
type streamOnlineMessage struct {
	StreamID  string `json:"streamId"`
	ChannelID string `json:"channelId"`
}

type streamOfflineMessage struct {
	ChannelID string `json:"channelId"`
}
