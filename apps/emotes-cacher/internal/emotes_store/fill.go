package emotes_store

import (
	"context"
	"log/slog"
	"sync"

	"github.com/satont/twir/apps/emotes-cacher/internal/emote"
	"github.com/satont/twir/apps/emotes-cacher/internal/services"
	bttvfetcher "github.com/satont/twir/apps/emotes-cacher/internal/services/bttv/fetcher"
	seventvfetcher "github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/fetcher"
	"github.com/twirapp/twir/libs/repositories/channels/model"
)

func (c *EmotesStore) fillChannels() {
	if len(c.channels) > 0 {
		return
	}

	var (
		page int64 = 0
		size int64 = 1000
	)
	for {
		var channelsIDs []string
		if err := c.gorm.
			Model(&model.Channel{}).
			Select("id").
			Where(`"isEnabled" = true`).
			Offset(int(page * size)).
			Limit(int(size)).
			Find(&channelsIDs).Error; err != nil {
			c.logger.Error("failed to get channels", slog.Any("error", err))
			return
		}
		if len(channelsIDs) == 0 {
			break
		}

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, channelID := range channelsIDs {
			wg.Add(1)

			go func() {
				defer wg.Done()

				c.logger.Info("fetching emotes for channel", slog.String("channel_id", channelID))
				sevenTvEmotes, err := seventvfetcher.GetChannelSevenTvEmotes(ctx, channelID)
				if err == nil {
					result := make([]emote.Emote, 0)
					for _, e := range sevenTvEmotes {
						result = append(result, e)
					}
					c.AddEmotes(
						ChannelID(channelID),
						services.ServiceSevenTV,
						result...,
					)
				} else {
					c.logger.Error(
						"failed to fetch 7tv emotes",
						slog.String("channel_id", channelID),
						slog.Any("error", err),
					)
				}

				bttvEmotes, err := bttvfetcher.GetChannelBttvEmotes(ctx, channelID)
				if err == nil {
					result := make([]emote.Emote, 0)
					for _, e := range bttvEmotes {
						result = append(result, e)
					}
					c.AddEmotes(
						ChannelID(channelID),
						services.ServiceBttv,
						result...,
					)
				} else {
					c.logger.Error("failed to fetch bttv emotes", slog.String("channel_id", channelID))
				}
			}()
		}

		wg.Wait()

		page++
	}
}
