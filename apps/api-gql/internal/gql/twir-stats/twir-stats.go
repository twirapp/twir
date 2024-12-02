package twir_stats

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm               *gorm.DB
	Logger             logger.Logger
	Config             config.Config
	GrpcTokensClient   tokens.TokensClient
	CachedTwitchClient *twitchcache.CachedTwitchClient
}

type TwirStats struct {
	gorm *gorm.DB

	cachedResponse     *gqlmodel.TwirStats
	logger             logger.Logger
	config             config.Config
	grpcTokensClients  tokens.TokensClient
	cachedTwitchClient *twitchcache.CachedTwitchClient
}

func New(opts Opts) *TwirStats {
	s := &TwirStats{
		gorm:               opts.Gorm,
		cachedResponse:     &gqlmodel.TwirStats{},
		logger:             opts.Logger,
		config:             opts.Config,
		grpcTokensClients:  opts.GrpcTokensClient,
		cachedTwitchClient: opts.CachedTwitchClient,
	}

	go s.cacheCounts()
	go s.cacheStreamers()

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			s.cacheCounts()
			s.cacheStreamers()
		}
	}()

	return s
}

func (c *TwirStats) GetCachedData() *gqlmodel.TwirStats {
	return c.cachedResponse
}

func (c *TwirStats) cacheCounts() {
	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		var count int64
		c.gorm.Model(&model.Users{}).Count(&count)
		c.cachedResponse.Viewers = int(count)
	}()

	go func() {
		defer wg.Done()

		var count int64
		c.gorm.Model(&model.Channels{}).Where(
			`"channels"."isEnabled" = ? AND "channels"."isTwitchBanned" = ? AND "User"."is_banned" = ?`,
			true,
			false,
			false,
		).Joins("User").Count(&count)
		c.cachedResponse.Channels = int(count)
	}()

	go func() {
		var count int64
		c.gorm.Model(&model.ChannelsCommands{}).
			Where("module = ?", "CUSTOM").
			Count(&count)
		c.cachedResponse.CreatedCommands = int(count)
	}()

	go func() {
		defer wg.Done()
		result := statsNResult{}
		c.gorm.Model(&model.UsersStats{}).
			Select("sum(messages) as n").
			Scan(&result)
		c.cachedResponse.Messages = int(result.N)
	}()

	go func() {
		defer wg.Done()
		var count int64
		c.gorm.Model(&model.ChannelEmoteUsage{}).Count(&count)
		c.cachedResponse.UsedEmotes = int(count)
	}()

	go func() {
		defer wg.Done()

		var count int64
		c.gorm.Model(&model.ChannelsCommandsUsages{}).Count(&count)
		c.cachedResponse.UsedCommands = int(count)
	}()

	wg.Wait()
}

type statsNResult struct {
	N int64
}

func (c *TwirStats) cacheStreamers() {
	var streamers []model.Channels
	if err := c.gorm.Model(&model.Channels{}).
		Where(
			`"isEnabled" = ? AND "isTwitchBanned" = ? AND "User"."is_banned" = ?`,
			true,
			false,
			false,
		).
		Joins("User").
		Find(&streamers).Error; err != nil {
		c.logger.Error("cannot cache streamers", slog.Any("err", err))
		return
	}

	streamers = lo.Filter(
		streamers,
		func(item model.Channels, _ int) bool {
			return item.User != nil && !item.User.HideOnLandingPage
		},
	)

	helixUsers := make([]helix.User, 0, len(streamers))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	usersIds := make([]string, 0, len(streamers))
	for _, streamer := range streamers {
		usersIds = append(usersIds, streamer.ID)
	}

	usersReq, usersErr := c.cachedTwitchClient.GetUsersByIds(ctx, usersIds)
	if usersErr != nil {
		c.logger.Error("cannot get users", slog.Any("err", usersErr))
		return
	}
	for _, user := range usersReq {
		if user.NotFound {
			continue
		}

		helixUsers = append(helixUsers, user.User)
	}

	streamersFollowers := make(map[string]int)
	streamersFollowersMu := sync.Mutex{}
	streamersFollowersWg := utils.NewGoroutinesGroup()

	// fetch channel followers
	for _, user := range helixUsers {
		user := user
		streamersFollowersWg.Go(
			func() {
				followers, err := c.cachedTwitchClient.GetChannelFollowersCountByChannelId(ctx, user.ID)
				if err != nil {
					c.logger.Error(
						"cannot get followers",
						slog.Any("err", err),
						slog.String("userId", user.ID),
					)
					return
				}

				streamersFollowersMu.Lock()
				streamersFollowers[user.ID] = followers
				streamersFollowersMu.Unlock()
			},
		)
	}

	streamersFollowersWg.Wait()

	streamersWithFollowers := make(
		[]gqlmodel.TwirStatsStreamer,
		0,
		len(streamersFollowers),
	)

	for userId, followers := range streamersFollowers {
		streamer, ok := lo.Find(
			helixUsers, func(item helix.User) bool {
				return item.ID == userId
			},
		)
		if !ok {
			continue
		}

		stream := model.ChannelsStreams{}
		if err := c.gorm.Where(`"userId" = ?`, streamer.ID).Find(&stream).Error; err != nil {
			c.logger.Error(
				"cannot get stream",
				slog.Any("err", err),
				slog.String("channelId", streamer.ID),
			)
			continue
		}

		streamersWithFollowers = append(
			streamersWithFollowers,
			gqlmodel.TwirStatsStreamer{
				ID: streamer.ID,
				TwitchProfile: &gqlmodel.TwirUserTwitchInfo{
					ID:              streamer.ID,
					Login:           streamer.Login,
					DisplayName:     streamer.DisplayName,
					ProfileImageURL: streamer.ProfileImageURL,
					Description:     streamer.Description,
				},
				IsLive:         stream.ID != "",
				IsPartner:      streamer.BroadcasterType == "partner",
				FollowersCount: followers,
			},
		)
	}

	c.logger.Info("Cache streamers updated", slog.Int("count", len(streamersWithFollowers)))
	c.cachedResponse.Streamers = streamersWithFollowers
}
