package stats

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/libs/api/messages/stats"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Stats struct {
	*impl_deps.Deps

	statsCache     *stats.Response
	streamersCache *stats.GetTwirStreamersResponse
}

type statsNResult struct {
	N int64
}

func New(deps *impl_deps.Deps) *Stats {
	s := &Stats{
		Deps:           deps,
		statsCache:     &stats.Response{},
		streamersCache: &stats.GetTwirStreamersResponse{},
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

func (c *Stats) cacheCounts() {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			var count int64
			c.Db.Model(&model.Users{}).Count(&count)
			c.statsCache.Users = count
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.Model(&model.Channels{}).Where(
				`"channels"."isEnabled" = ? AND "channels"."isTwitchBanned" = ? AND "User"."is_banned" = ?`,
				true,
				false,
				false,
			).Joins("User").Count(&count)
			c.statsCache.Channels = count
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.Model(&model.ChannelsCommands{}).
				Where("module = ?", "CUSTOM").
				Count(&count)

			c.statsCache.Commands = count
		},
	)

	wg.Go(
		func() {
			result := statsNResult{}
			c.Db.Model(&model.UsersStats{}).
				Select("sum(messages) as n").
				Scan(&result)
			c.statsCache.Messages = result.N
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.Model(&model.ChannelEmoteUsage{}).Count(&count)
			c.statsCache.UsedEmotes = count
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.Model(&model.ChannelsCommandsUsages{}).Count(&count)
			c.statsCache.UsedCommands = count
		},
	)

	wg.Wait()
}

func (c *Stats) cacheStreamers() {
	c.Logger.Info("Updating streamers cache in stats")
	var streamers []model.Channels
	if err := c.Db.Model(&model.Channels{}).
		Where(
			`"isEnabled" = ? AND "isTwitchBanned" = ? AND "User"."is_banned" = ?`,
			true,
			false,
			false,
		).
		Joins("User").
		Find(&streamers).Error; err != nil {
		c.Logger.Error("cannot cache streamers", slog.Any("err", err))
		return
	}

	streamers = lo.Filter(
		streamers,
		func(item model.Channels, _ int) bool {
			return item.User != nil && !item.User.HideOnLandingPage
		},
	)

	helixUsersMu := sync.Mutex{}
	helixUsers := make([]helix.User, 0, len(streamers))

	usersWgGrp, ctx := errgroup.WithContext(context.Background())
	chunks := lo.Chunk(streamers, 100)

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.Config, c.Grpc.Tokens)
	if err != nil {
		c.Logger.Error("cannot create twitch client", slog.Any("err", err))
		return
	}

	for _, chunk := range chunks {
		chunk := chunk
		usersWgGrp.Go(
			func() error {
				usersIds := lo.Map(
					chunk, func(item model.Channels, _ int) string {
						return item.ID
					},
				)

				usersReq, usersErr := twitchClient.GetUsers(
					&helix.UsersParams{
						IDs: usersIds,
					},
				)
				if usersErr != nil {
					return usersErr
				}
				if usersReq.ErrorMessage != "" {
					return errors.New(usersReq.ErrorMessage)
				}

				helixUsersMu.Lock()
				defer helixUsersMu.Unlock()
				helixUsers = append(helixUsers, usersReq.Data.Users...)

				return nil
			},
		)
	}

	if err := usersWgGrp.Wait(); err != nil {
		c.Logger.Error("cannot get users", slog.Any("err", err))
		return
	}

	streamersFollowers := make(map[string]int)
	streamersFollowersMu := sync.Mutex{}
	streamersFollowersWg := utils.NewGoroutinesGroup()

	for _, user := range helixUsers {
		user := user
		streamersFollowersWg.Go(
			func() {
				userTwitchClientCtx, userTwitchClientCtxCancel := context.WithTimeout(
					context.Background(),
					60*time.Second,
				)
				defer userTwitchClientCtxCancel()
				userTwitchClient, err := twitch.NewUserClientWithContext(
					userTwitchClientCtx,
					user.ID,
					c.Config,
					c.Grpc.Tokens,
				)
				if err != nil {
					c.Logger.Error("cannot create twitch client", slog.Any("err", err))
					return
				}

				followersReq, followersErr := userTwitchClient.GetChannelFollows(
					&helix.GetChannelFollowsParams{
						BroadcasterID: user.ID,
					},
				)
				if followersErr != nil {
					c.Logger.Error(
						"cannot get followers",
						slog.Any("err", followersErr),
						slog.Group(
							"user",
							slog.String("id", user.ID),
							slog.String("login", user.Login),
						),
					)
					return
				}
				if followersReq.ErrorMessage != "" {
					c.Logger.Error(
						"cannot get followers",
						slog.Any("err", followersReq.ErrorMessage),
						slog.Group(
							"user",
							slog.String("id", user.ID),
							slog.String("login", user.Login),
						),
						slog.Int("status", followersReq.StatusCode),
					)
					return
				}

				streamersFollowersMu.Lock()
				defer streamersFollowersMu.Unlock()
				streamersFollowers[user.ID] = followersReq.Data.Total
			},
		)
	}

	streamersFollowersWg.Wait()

	streamersWithFollowers := make(
		[]*stats.GetTwirStreamersResponse_Streamer,
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
		if err := c.Db.Where(`"userId" = ?`, streamer.ID).Find(&stream).Error; err != nil {
			c.Logger.Error(
				"cannot get stream",
				slog.Any("err", err),
				slog.String("channelId", streamer.ID),
			)
			continue
		}

		streamersWithFollowers = append(
			streamersWithFollowers,
			&stats.GetTwirStreamersResponse_Streamer{
				UserId:          streamer.ID,
				UserLogin:       streamer.Login,
				UserDisplayName: streamer.DisplayName,
				Avatar:          streamer.ProfileImageURL,
				FollowersCount:  int32(followers),
				IsLive:          stream.ID != "",
				IsPartner:       streamer.BroadcasterType == "partner",
			},
		)
	}

	c.Logger.Info("Cache streamers updated", slog.Int("count", len(streamersWithFollowers)))
	c.streamersCache.Streamers = streamersWithFollowers
}

func (c *Stats) GetStats(_ context.Context, _ *emptypb.Empty) (*stats.Response, error) {
	return c.statsCache, nil
}

func (c *Stats) GetStatsTwirStreamers(
	_ context.Context,
	_ *emptypb.Empty,
) (*stats.GetTwirStreamersResponse, error) {
	return c.streamersCache, nil
}
