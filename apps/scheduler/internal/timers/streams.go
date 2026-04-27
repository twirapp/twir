package timers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	buscore "github.com/twirapp/twir/libs/bus-core"
	buskick "github.com/twirapp/twir/libs/bus-core/kick"
	bustokens "github.com/twirapp/twir/libs/bus-core/tokens"
	bustwitch "github.com/twirapp/twir/libs/bus-core/twitch"
	config "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/scorfly/gokick"
	"github.com/samber/lo"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
)

type StreamOpts struct {
	fx.In
	Lc fx.Lifecycle

	Config config.Config
	Logger *slog.Logger

	Gorm    *gorm.DB
	TwirBus *buscore.Bus
}

type streams struct {
	config  config.Config
	logger  *slog.Logger
	gorm    *gorm.DB
	twirBus *buscore.Bus
}

func NewStreams(opts StreamOpts) {
	timeTick := 15 * time.Second
	if opts.Config.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &streams{
		config:  opts.Config,
		logger:  opts.Logger,
		gorm:    opts.Gorm,
		twirBus: opts.TwirBus,
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
								opts.Logger.Error("cannot process streams", logger.Error(err))
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

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
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
				c.logger.Error("cannot get streams", logger.Error(err))
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
					ID:           twitchStream.ID,
					UserId:       userId,
					UserLogin:    twitchStream.UserLogin,
					UserName:     twitchStream.UserName,
					GameId:       twitchStream.GameID,
					GameName:     twitchStream.GameName,
					CommunityIds: nil,
					Type:         twitchStream.Type,
					Title:        twitchStream.Title,
					ViewerCount:  twitchStream.ViewerCount,
					StartedAt:    twitchStream.StartedAt,
					Language:     twitchStream.Language,
					ThumbnailUrl: twitchStream.ThumbnailURL,
					TagIds:       nil,
					Tags:         tags,
					IsMature:     twitchStream.IsMature,
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

					c.twirBus.Channel.StreamOnline.Publish(
						ctx,
						bustwitch.StreamOnlineMessage{
							ChannelID:    channelStream.UserId,
							StreamID:     channelStream.ID,
							CategoryName: channelStream.GameName,
							CategoryID:   channelStream.GameId,
							Title:        channelStream.Title,
							Viewers:      channelStream.ViewerCount,
							StartedAt:    channelStream.StartedAt,
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
						c.logger.Error("cannot delete stream", logger.Error(err))
						return
					}

					c.twirBus.Channel.StreamOffline.Publish(
						ctx,
						bustwitch.StreamOfflineMessage{
							ChannelID: channelStream.UserId,
							StartedAt: dbStream.StartedAt,
						},
					)
				}
			}
		}(chunk)
	}

	wg.Wait()

	if err := c.processKickStreams(ctx, existedStreams); err != nil {
		return fmt.Errorf("cannot process kick streams: %w", err)
	}

	return nil
}

type kickChannelRow struct {
	ID             string  `gorm:"column:id"`
	KickPlatformID *string `gorm:"column:kick_platform_id"`
}

func (c *streams) processKickStreams(ctx context.Context, existedStreams []model.ChannelsStreams) error {
	var channels []kickChannelRow
	err := c.gorm.
		WithContext(ctx).
		Table("channels").
		Select(`channels.id, channels.kick_platform_id`).
		Joins(`LEFT JOIN users ON users.id = channels.kick_user_id`).
		Where(`channels."isEnabled" = ?`, true).
		Where(`channels.kick_platform_id IS NOT NULL`).
		Where(`COALESCE(users.is_banned, false) = false`).
		Scan(&channels).Error
	if err != nil {
		return err
	}

	if len(channels) == 0 {
		return nil
	}

	appToken, err := c.twirBus.Tokens.RequestAppToken.Request(ctx, bustokens.GetAppTokenRequest{Platform: platformentity.PlatformKick})
	if err != nil {
		return fmt.Errorf("request kick app token: %w", err)
	}

	kickClient, err := gokick.NewClient(&gokick.ClientOptions{AppAccessToken: appToken.Data.AccessToken})
	if err != nil {
		return fmt.Errorf("create kick client: %w", err)
	}

	validChannels := lo.Filter(channels, func(channel kickChannelRow, _ int) bool {
		return channel.KickPlatformID != nil && *channel.KickPlatformID != ""
	})

	chunks := lo.Chunk(validChannels, 50)
	for _, chunk := range chunks {
		platformIDs := make([]int, 0, len(chunk))
		byPlatformID := make(map[string]kickChannelRow, len(chunk))
		for _, channel := range chunk {
			platformIDInt, convErr := strconv.Atoi(*channel.KickPlatformID)
			if convErr != nil {
				c.logger.Error("cannot parse kick platform id", slog.Any("err", convErr), slog.String("channel_id", channel.ID))
				continue
			}
			platformIDs = append(platformIDs, platformIDInt)
			byPlatformID[*channel.KickPlatformID] = channel
		}

		if len(platformIDs) == 0 {
			continue
		}

		resp, reqErr := kickClient.GetChannels(ctx, gokick.NewChannelListFilter().SetBroadcasterUserIDs(platformIDs))
		if reqErr != nil {
			return fmt.Errorf("get kick channels: %w", reqErr)
		}

		channelsByPlatformID := make(map[string]gokick.ChannelResponse, len(resp.Result))
		for _, item := range resp.Result {
			channelsByPlatformID[strconv.Itoa(item.BroadcasterUserID)] = item
		}

		for _, channel := range chunk {
			platformID := *channel.KickPlatformID
			kickChannel, exists := channelsByPlatformID[platformID]
			dbStream, dbStreamExists := lo.Find(existedStreams, func(stream model.ChannelsStreams) bool {
				return stream.UserId == channel.ID
			})

			if exists && kickChannel.Stream.IsLive {
				startedAt := time.Now().UTC()
				if kickChannel.Stream.StartTime != "" {
					if parsedStartedAt, parseErr := time.Parse(time.RFC3339, kickChannel.Stream.StartTime); parseErr == nil {
						startedAt = parsedStartedAt
					}
				}

				tags := &pq.StringArray{}
				for _, tag := range kickChannel.Stream.CustomTags {
					*tags = append(*tags, tag)
				}

				channelStream := &model.ChannelsStreams{
					ID:           channel.ID,
					UserId:       channel.ID,
					UserLogin:    kickChannel.Slug,
					UserName:     kickChannel.Slug,
					GameId:       strconv.Itoa(kickChannel.Category.ID),
					GameName:     kickChannel.Category.Name,
					CommunityIds: nil,
					Type:         "live",
					Title:        kickChannel.StreamTitle,
					ViewerCount:  kickChannel.Stream.ViewerCount,
					StartedAt:    startedAt,
					Language:     kickChannel.Stream.Language,
					ThumbnailUrl: kickChannel.Stream.Thumbnail,
					TagIds:       nil,
					Tags:         tags,
					IsMature:     kickChannel.Stream.IsMature,
				}

				if result := c.gorm.WithContext(ctx).Where(`"userId" = ?`, channel.ID).Save(channelStream); result.Error != nil {
					c.logger.Error("cannot save kick stream", slog.Any("err", result.Error))
					continue
				}

				if !dbStreamExists {
					c.twirBus.KickStreamOnline.Publish(ctx, buskick.KickStreamOnline{
						BroadcasterUserID:    platformID,
						BroadcasterUserLogin: kickChannel.Slug,
					})
				}
				continue
			}

			if dbStreamExists {
				if err := c.gorm.WithContext(ctx).Where(`"userId" = ?`, channel.ID).Delete(&model.ChannelsStreams{}).Error; err != nil {
					c.logger.Error("cannot delete kick stream", logger.Error(err))
					continue
				}

				c.twirBus.KickStreamOffline.Publish(ctx, buskick.KickStreamOffline{
					BroadcasterUserID:    platformID,
					BroadcasterUserLogin: dbStream.UserLogin,
				})
			}
		}
	}

	return nil
}
