package stats

import (
	"context"
	"time"

	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/stats"
	"github.com/satont/twir/libs/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Stats struct {
	*impl_deps.Deps

	cache *stats.Response
}

type statsNResult struct {
	N int64
}

func New(deps *impl_deps.Deps) *Stats {
	s := &Stats{
		Deps:  deps,
		cache: &stats.Response{},
	}

	go s.cacheCounts()
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			s.cacheCounts()
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
			c.cache.Users = count
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.Model(&model.Channels{}).Where(
				`"isEnabled" = ? AND "isTwitchBanned" = ? AND "isBanned" = ?`,
				true,
				false,
				false,
			).Count(&count)
			c.cache.Channels = count
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.Model(&model.ChannelsCommands{}).
				Where("module = ?", "CUSTOM").
				Count(&count)

			c.cache.Commands = count
		},
	)

	wg.Go(
		func() {
			result := statsNResult{}
			c.Db.Model(&model.UsersStats{}).
				Select("sum(messages) as n").
				Scan(&result)
			c.cache.Messages = result.N
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.Model(&model.ChannelEmoteUsage{}).Count(&count)
			c.cache.UsedEmotes = count
		},
	)

	wg.Wait()
}

func (c *Stats) GetStats(_ context.Context, _ *emptypb.Empty) (*stats.Response, error) {
	return c.cache, nil
}
