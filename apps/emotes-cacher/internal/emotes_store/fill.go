package emotes_store

import (
	"context"
	"log/slog"
	"sync"

	"github.com/satont/twir/apps/emotes-cacher/internal/emote"
	bttvfetcher "github.com/satont/twir/apps/emotes-cacher/internal/services/bttv/fetcher"
	ffzfetcher "github.com/satont/twir/apps/emotes-cacher/internal/services/ffz/fetcher"
	seventvfetcher "github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/fetcher"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/repositories/channels/model"
)

func (c *EmotesStore) fillChannels() {
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

				c.logger.Debug("fetching emotes for channel", slog.String("channel_id", channelID))
				sevenTvEmotes, err := seventvfetcher.GetChannelSevenTvEmotes(ctx, channelID)
				if err == nil {
					result := make([]emote.Emote, 0)
					for _, e := range sevenTvEmotes {
						result = append(result, e)
					}
					c.AddEmotes(
						ChannelID(channelID),
						emotes_cacher.ServiceNameSevenTV,
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
						emotes_cacher.ServiceNameBTTV,
						result...,
					)
				} else {
					c.logger.Error("failed to fetch bttv emotes", slog.String("channel_id", channelID))
				}

				ffzEmotes, err := ffzfetcher.GetChannelFfzEmotes(ctx, channelID)
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
					c.logger.Error("failed to fetch ffz emotes", slog.String("channel_id", channelID))
				}
			}()
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
