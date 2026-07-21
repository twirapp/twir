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
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/scorfly/gokick"
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

	StreamsRepo    streamsrepository.Repository
	ChannelService *channelservice.ChannelService
}

type streams struct {
	config  config.Config
	logger  *slog.Logger
	gorm    *gorm.DB
	twirBus *buscore.Bus

	streamsRepo    streamsrepository.Repository
	channelService *channelservice.ChannelService
}

func NewStreams(opts StreamOpts) {
	timeTick := 15 * time.Second
	if opts.Config.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &streams{
		config:         opts.Config,
		logger:         opts.Logger,
		gorm:           opts.Gorm,
		twirBus:        opts.TwirBus,
		streamsRepo:    opts.StreamsRepo,
		channelService: opts.ChannelService,
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
	var channels []twitchStreamChannelRow
	err := buildTwitchChannelsQuery(c.gorm, ctx).Scan(&channels).Error
	if err != nil {
		return fmt.Errorf("cannot get channels: %w", err)
	}

	usersIds := make([]string, 0, len(channels))
	for _, channel := range channels {
		if channel.TwitchPlatformID != nil && *channel.TwitchPlatformID != "" {
			usersIds = append(usersIds, *channel.TwitchPlatformID)
		}
	}

	discordIntegration := &model.Integrations{}
	err = c.gorm.
		WithContext(ctx).
		Where(`service = ?`, model.IntegrationServiceDiscord).
		Select("id").
		Find(discordIntegration).Error
	if err == nil {
		var discordIntegrations []model.ChannelsIntegrations
		if discordIntegration.ID != "" {
			err = c.gorm.
				WithContext(ctx).
				Where(`"integrationId" = ?`, discordIntegration.ID).
				Select("id", `"integrationId"`, "data").
				Find(&discordIntegrations).Error

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
	} else {
		c.logger.InfoContext(ctx, "Cannot get discord integration", logger.Error(err))
	}

	usersIds = lo.Uniq(usersIds)

	existedStreams, err := c.streamsRepo.GetList(ctx)
	if err != nil {
		return fmt.Errorf("cannot get existed streams: %w", err)
	}

	existedTwitchStreams := lo.Filter(
		existedStreams, func(stream streamsmodel.Stream, _ int) bool {
			return stream.Platform == platformentity.PlatformTwitch
		},
	)

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
					existedTwitchStreams, func(stream streamsmodel.Stream) bool {
						return stream.UserId == userId
					},
				)

				if twitchStreamExists {
					channel, err := c.channelService.GetChannelByPlatformUserID(
						ctx,
						userId,
						platformentity.PlatformTwitch,
					)
					if err != nil || channel.IsNil() {
						c.logger.Error(
							"cannot resolve channel for twitch user",
							slog.String("twitch_user_id", userId),
							logger.Error(err),
						)
						continue
					}

					if err := c.streamsRepo.Save(
						ctx,
						streamsrepository.SaveInput{
							ID:           twitchStream.ID,
							ChannelID:    channel.ID,
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
							Tags:         twitchStream.Tags,
							IsMature:     twitchStream.IsMature,
							Platform:     platformentity.PlatformTwitch,
						},
					); err != nil {
						c.logger.Error("cannot save stream", slog.Any("err", err))
						continue
					}

					if err := c.channelService.InvalidateOnlineCache(ctx, channel.ID); err != nil {
						c.logger.Error("cannot invalidate online cache", logger.Error(err))
					}

					if !dbStreamExists {
						c.twirBus.Channel.StreamOnline.Publish(
							ctx,
							bustwitch.StreamOnlineMessage{
								ChannelID:    userId,
								StreamID:     twitchStream.ID,
								CategoryName: twitchStream.GameName,
								CategoryID:   twitchStream.GameID,
								Title:        twitchStream.Title,
								Viewers:      twitchStream.ViewerCount,
								StartedAt:    twitchStream.StartedAt,
							},
						)
					}

					continue
				}

				if dbStreamExists {
					// stream offline, delete
					if err := c.streamsRepo.DeleteByChannelID(
						ctx,
						dbStream.ChannelID,
						platformentity.PlatformTwitch,
					); err != nil {
						c.logger.Error("cannot delete stream", logger.Error(err))
						continue
					}

					if err := c.channelService.InvalidateOnlineCache(ctx, dbStream.ChannelID); err != nil {
						c.logger.Error("cannot invalidate online cache", logger.Error(err))
					}

					c.twirBus.Channel.StreamOffline.Publish(
						ctx,
						bustwitch.StreamOfflineMessage{
							ChannelID: userId,
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

type twitchStreamChannelRow struct {
	ID               string  `gorm:"column:id"`
	TwitchPlatformID *string `gorm:"column:twitch_platform_id"`
}

const (
	twitchChannelsSelectClause        = `channels.id, users.platform_id AS twitch_platform_id`
	twitchChannelsJoinClause          = `LEFT JOIN users ON users.id = channels.twitch_user_id`
	twitchChannelsPlatformIDIsNotNull = `users.platform_id IS NOT NULL`
)

func buildTwitchChannelsQuery(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.
		WithContext(ctx).
		Table("channels").
		Select(twitchChannelsSelectClause).
		Joins(twitchChannelsJoinClause).
		Where(`channels.twitch_bot_enabled IS TRUE`).
		Where(twitchChannelsPlatformIDIsNotNull).
		Where(`COALESCE(users.is_banned, false) = false`)
}

type kickChannelRow struct {
	ID             string  `gorm:"column:id"`
	KickPlatformID *string `gorm:"column:kick_platform_id"`
}

const (
	kickChannelsSelectClause        = `channels.id, users.platform_id AS kick_platform_id`
	kickChannelsJoinClause          = `LEFT JOIN users ON users.id = channels.kick_user_id`
	kickChannelsPlatformIDIsNotNull = `users.platform_id IS NOT NULL`
)

func buildKickChannelsQuery(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.
		WithContext(ctx).
		Table("channels").
		Select(kickChannelsSelectClause).
		Joins(kickChannelsJoinClause).
		Where(`channels."isEnabled" = ?`, true).
		Where(kickChannelsPlatformIDIsNotNull).
		Where(`COALESCE(users.is_banned, false) = false`)
}

func (c *streams) processKickStreams(ctx context.Context, existedStreams []streamsmodel.Stream) error {
	var channels []kickChannelRow
	err := buildKickChannelsQuery(c.gorm, ctx).Scan(&channels).Error
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

	existedKickStreams := lo.Filter(
		existedStreams, func(stream streamsmodel.Stream, _ int) bool {
			return stream.Platform == platformentity.PlatformKick
		},
	)

	validChannels := lo.Filter(channels, func(channel kickChannelRow, _ int) bool {
		return channel.KickPlatformID != nil && *channel.KickPlatformID != ""
	})

	chunks := lo.Chunk(validChannels, 50)
	for _, chunk := range chunks {
		platformIDs := make([]int, 0, len(chunk))
		for _, channel := range chunk {
			platformIDInt, convErr := strconv.Atoi(*channel.KickPlatformID)
			if convErr != nil {
				c.logger.Error("cannot parse kick platform id", slog.Any("err", convErr), slog.String("channel_id", channel.ID))
				continue
			}
			platformIDs = append(platformIDs, platformIDInt)
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
			dbStream, dbStreamExists := lo.Find(existedKickStreams, func(stream streamsmodel.Stream) bool {
				return stream.UserId == platformID
			})

			channelUUID, uuidErr := uuid.Parse(channel.ID)
			if uuidErr != nil {
				c.logger.Error("cannot parse channel id", slog.Any("err", uuidErr), slog.String("channel_id", channel.ID))
				continue
			}

			if exists && kickChannel.Stream.IsLive {
				startedAt := time.Now().UTC()
				if kickChannel.Stream.StartTime != "" {
					if parsedStartedAt, parseErr := time.Parse(time.RFC3339, kickChannel.Stream.StartTime); parseErr == nil {
						startedAt = parsedStartedAt
					}
				}

				if err := c.streamsRepo.Save(
					ctx,
					streamsrepository.SaveInput{
						ID:           "",
						ChannelID:    channelUUID,
						UserId:       platformID,
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
						Tags:         kickChannel.Stream.CustomTags,
						IsMature:     kickChannel.Stream.IsMature,
						Platform:     platformentity.PlatformKick,
					},
				); err != nil {
					c.logger.Error("cannot save kick stream", slog.Any("err", err))
					continue
				}

				if err := c.channelService.InvalidateOnlineCache(ctx, channelUUID); err != nil {
					c.logger.Error("cannot invalidate online cache", logger.Error(err))
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
				if err := c.streamsRepo.DeleteByChannelID(
					ctx,
					channelUUID,
					platformentity.PlatformKick,
				); err != nil {
					c.logger.Error("cannot delete kick stream", logger.Error(err))
					continue
				}

				if err := c.channelService.InvalidateOnlineCache(ctx, channelUUID); err != nil {
					c.logger.Error("cannot invalidate online cache", logger.Error(err))
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
