package stats

import (
	"context"
	"sync"

	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/stats"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Stats struct {
	*impl_deps.Deps
}

type statsItem struct {
	Count int64  `json:"count"`
	Name  string `json:"name"`
}

type statsNResult struct {
	N int64
}

func (c *Stats) GetStats(ctx context.Context, _ *emptypb.Empty) (*stats.Response, error) {
	wg := sync.WaitGroup{}
	statistic := []statsItem{
		{Name: "users"},
		{Name: "channels"},
		{Name: "commands"},
		{Name: "messages"},
		{Name: "used_emotes"},
	}

	wg.Add(5)

	go func() {
		defer wg.Done()
		var count int64
		c.Db.WithContext(ctx).Model(&model.Users{}).Count(&count)
		statistic[0].Count = count
	}()

	go func() {
		defer wg.Done()
		var count int64
		c.Db.WithContext(ctx).Model(&model.Channels{}).Where(`"isEnabled" = ?`, true).Count(&count)
		statistic[1].Count = count
	}()

	go func() {
		defer wg.Done()
		var count int64
		c.Db.WithContext(ctx).Model(&model.ChannelsCommands{}).
			Where("module = ?", "CUSTOM").
			Count(&count)

		statistic[2].Count = count
	}()

	go func() {
		defer wg.Done()
		result := statsNResult{}
		c.Db.WithContext(ctx).Model(&model.UsersStats{}).
			Select("sum(messages) as n").
			Scan(&result)
		statistic[3].Count = result.N
	}()

	go func() {
		defer wg.Done()
		var count int64
		c.Db.WithContext(ctx).Model(&model.ChannelEmoteUsage{}).Count(&count)
		statistic[3].Count = count
	}()

	wg.Wait()

	return &stats.Response{
		Users:      statistic[0].Count,
		Channels:   statistic[1].Count,
		Commands:   statistic[2].Count,
		Messages:   statistic[3].Count,
		UsedEmotes: statistic[4].Count,
	}, nil
}
