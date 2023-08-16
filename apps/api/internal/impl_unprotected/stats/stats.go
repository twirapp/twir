package stats

import (
	"context"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/stats"
	"github.com/satont/twir/libs/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Stats struct {
	*impl_deps.Deps
}

type statsNResult struct {
	N int64
}

func (c *Stats) GetStats(ctx context.Context, _ *emptypb.Empty) (*stats.Response, error) {
	s := utils.NewSyncMap[int64]()

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			var count int64
			c.Db.WithContext(ctx).Model(&model.Users{}).Count(&count)
			s.Add("users", count)
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.WithContext(ctx).Model(&model.Channels{}).Where(`"isEnabled" = ?`, true).Count(&count)
			s.Add("channels", count)
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.WithContext(ctx).Model(&model.ChannelsCommands{}).
				Where("module = ?", "CUSTOM").
				Count(&count)

			s.Add("commands", count)
		},
	)

	wg.Go(
		func() {
			result := statsNResult{}
			c.Db.WithContext(ctx).Model(&model.UsersStats{}).
				Select("sum(messages) as n").
				Scan(&result)
			s.Add("messages", result.N)
		},
	)

	wg.Go(
		func() {
			var count int64
			c.Db.WithContext(ctx).Model(&model.ChannelEmoteUsage{}).Count(&count)
			s.Add("used_emotes", count)
		},
	)

	wg.Wait()

	return &stats.Response{
		Users:      s.Get("users"),
		Channels:   s.Get("channels"),
		Commands:   s.Get("commands"),
		Messages:   s.Get("messages"),
		UsedEmotes: s.Get("used_emotes"),
	}, nil
}
