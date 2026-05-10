package emotes_store

import (
	"context"
	"log/slog"
	"sync"

	"github.com/twirapp/twir/apps/emotes-cacher/internal/emote"
	bttvfetcher "github.com/twirapp/twir/apps/emotes-cacher/internal/services/bttv/fetcher"
	ffzfetcher "github.com/twirapp/twir/apps/emotes-cacher/internal/services/ffz/fetcher"
	seventvfetcher "github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv/fetcher"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
)

type startupChannelData struct {
	ID               string  `gorm:"column:id"`
	TwitchPlatformID *string `gorm:"column:twitch_platform_id"`
	KickPlatformID   *string `gorm:"column:kick_platform_id"`
}

func (c *EmotesStore) fillChannels() {
	var (
		page int64 = 0
		size int64 = 1000
	)
	for {
		var channelsData []startupChannelData
		if err := c.gorm.
			Table("channels c").
			Select("c.id", "tu.platform_id as twitch_platform_id", "ku.platform_id as kick_platform_id").
			Joins("LEFT JOIN users tu ON tu.id = c.twitch_user_id").
			Joins("LEFT JOIN users ku ON ku.id = c.kick_user_id").
			Where(`c."isEnabled" = true`).
			Offset(int(page * size)).
			Limit(int(size)).
			Scan(&channelsData).Error; err != nil {
			c.logger.Error("failed to get channels", logger.Error(err))
			return
		}
		if len(channelsData) == 0 {
			break
		}

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, channelData := range channelsData {
			wg.Add(1)

			go func(channelData startupChannelData) {
				defer wg.Done()
				channelID := channelData.ID

				c.logger.Debug("fetching emotes for channel", slog.String("channel_id", channelID))
				if channelData.TwitchPlatformID != nil {
					sevenTvEmotes, err := seventvfetcher.GetChannelSevenTvEmotes(ctx, platform.PlatformTwitch, *channelData.TwitchPlatformID)
					if err == nil {
						result := make([]emote.Emote, 0)
						for _, e := range sevenTvEmotes {
							result = append(result, e)
						}
						c.AddEmotes(ChannelID(channelID), emotes_cacher.ServiceNameSevenTV, result...)
					} else {
						c.logger.Debug("failed to fetch 7tv emotes", slog.String("channel_id", channelID), logger.Error(err))
					}
				} else if channelData.KickPlatformID != nil {
					sevenTvEmotes, err := seventvfetcher.GetChannelSevenTvEmotes(ctx, platform.PlatformKick, *channelData.KickPlatformID)
					if err == nil {
						result := make([]emote.Emote, 0)
						for _, e := range sevenTvEmotes {
							result = append(result, e)
						}
						c.AddEmotes(ChannelID(channelID), emotes_cacher.ServiceNameSevenTV, result...)
					} else {
						c.logger.Debug("failed to fetch 7tv emotes", slog.String("channel_id", channelID), logger.Error(err))
					}
				}

				if channelData.TwitchPlatformID != nil {
					bttvEmotes, err := bttvfetcher.GetChannelBttvEmotes(ctx, *channelData.TwitchPlatformID)
					if err == nil {
						result := make([]emote.Emote, 0)
						for _, e := range bttvEmotes {
							result = append(result, e)
						}
						c.AddEmotes(
							ChannelID(channelID),
							emotes_cacher.ServiceNameBTTV,
							result...,
						)
					} else {
						c.logger.Debug("failed to fetch bttv emotes", slog.String("channel_id", channelID))
					}

					ffzEmotes, err := ffzfetcher.GetChannelFfzEmotes(ctx, *channelData.TwitchPlatformID)
					if err == nil {
						result := make([]emote.Emote, 0)
						for _, e := range ffzEmotes {
							result = append(result, e)
						}
						c.AddEmotes(
							ChannelID(channelID),
							emotes_cacher.ServiceNameFFZ,
							result...,
						)
					} else {
						c.logger.Debug("failed to fetch ffz emotes", slog.String("channel_id", channelID))
					}
				}
			}(channelData)
		}

		wg.Wait()

		page++
	}
}

func (c *EmotesStore) fillGlobal() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}

	wg.Add(3)

	go func() {
		defer wg.Done()
		em, err := seventvfetcher.GetGlobalSevenTvEmotes(ctx)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		c.AddEmotes(GlobalChannelID, emotes_cacher.ServiceNameSevenTV, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := bttvfetcher.GetGlobalBttvEmotes(ctx)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		c.AddEmotes(GlobalChannelID, emotes_cacher.ServiceNameBTTV, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := ffzfetcher.GetGlobalFfzEmotes(ctx)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		c.AddEmotes(GlobalChannelID, emotes_cacher.ServiceNameFFZ, em...)
	}()

	wg.Wait()
}
