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
	"gorm.io/gorm"
)

type startupChannelData struct {
	ID         string            `gorm:"column:id"`
	Platform   platform.Platform `gorm:"column:platform"`
	PlatformID string            `gorm:"column:platform_id"`
}

func buildStartupChannelsQuery(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.
		WithContext(ctx).
		Table("channel_platforms AS cp").
		Select("cp.channel_id AS id", "cp.platform", "cp.platform_channel_id AS platform_id").
		Joins("JOIN channels c ON c.id = cp.channel_id").
		Where(`c."isEnabled" = ?`, true).
		Where("cp.platform IN ?", []platform.Platform{
			platform.PlatformTwitch,
			platform.PlatformKick,
		})
}

func (c *EmotesStore) fillChannels() {
	var (
		page int64 = 0
		size int64 = 1000
	)
	for {
		var channelsData []startupChannelData
		if err := buildStartupChannelsQuery(c.gorm, context.Background()).
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
				if channelData.Platform != platform.PlatformTwitch && channelData.Platform != platform.PlatformKick {
					return
				}

				c.logger.Debug("fetching emotes for channel", slog.String("channel_id", channelID))
				channelKey := ChannelKey{Platform: channelData.Platform, ID: channelData.PlatformID}

				sevenTvEmotes, err := seventvfetcher.GetChannelSevenTvEmotes(ctx, channelData.Platform, channelData.PlatformID)
				if err == nil {
					result := make([]emote.Emote, 0)
					for _, e := range sevenTvEmotes {
						result = append(result, e)
					}
					c.AddEmotes(channelKey, emotes_cacher.ServiceNameSevenTV, result...)
				} else {
					c.logger.Debug("failed to fetch 7tv emotes", slog.String("channel_id", channelID), logger.Error(err))
				}

				if channelData.Platform == platform.PlatformTwitch {
					bttvEmotes, err := bttvfetcher.GetChannelBttvEmotes(ctx, channelData.PlatformID)
					if err == nil {
						result := make([]emote.Emote, 0)
						for _, e := range bttvEmotes {
							result = append(result, e)
						}
						c.AddEmotes(
							channelKey,
							emotes_cacher.ServiceNameBTTV,
							result...,
						)
					} else {
						c.logger.Debug("failed to fetch bttv emotes", slog.String("channel_id", channelID))
					}

					ffzEmotes, err := ffzfetcher.GetChannelFfzEmotes(ctx, channelData.PlatformID)
					if err == nil {
						result := make([]emote.Emote, 0)
						for _, e := range ffzEmotes {
							result = append(result, e)
						}
						c.AddEmotes(
							channelKey,
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

		c.AddEmotes(GlobalChannelKey, emotes_cacher.ServiceNameSevenTV, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := bttvfetcher.GetGlobalBttvEmotes(ctx)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		c.AddEmotes(GlobalChannelKey, emotes_cacher.ServiceNameBTTV, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := ffzfetcher.GetGlobalFfzEmotes(ctx)
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		c.AddEmotes(GlobalChannelKey, emotes_cacher.ServiceNameFFZ, em...)
	}()

	wg.Wait()
}
