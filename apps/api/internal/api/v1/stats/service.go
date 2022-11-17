package stats

import (
	"sync"
	model "tsuwari/models"

	"github.com/satont/tsuwari/apps/api/internal/types"
)

type nResult struct {
	N int64
}

type statsItem struct {
	Count int64 `json:"count"`
}

type stats struct {
	Users     statsItem `json:"users"`
	Streamers statsItem `json:"streamers"`
	Commands  statsItem `json:"commands"`
	Messages  statsItem `json:"messages"`
}

func handleGet(services types.Services) (*stats, error) {
	wg := sync.WaitGroup{}
	statistic := stats{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		var count int64
		err := services.DB.Model(&model.Users{}).Count(&count).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
		} else {
			statistic.Users = statsItem{Count: count}
		}
	}()

	go func() {
		defer wg.Done()
		var count int64
		err := services.DB.Model(&model.Channels{}).Count(&count).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
		} else {
			statistic.Streamers = statsItem{Count: count}
		}
	}()

	go func() {
		defer wg.Done()
		var count int64
		err := services.DB.Model(&model.ChannelsCommands{}).Count(&count).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
		} else {
			statistic.Commands = statsItem{Count: count}
		}
	}()

	go func() {
		defer wg.Done()
		result := nResult{}
		err := services.DB.Model(&model.UsersStats{}).
			Select("sum(messages) as n").
			Scan(&result).
			Error
		if err != nil {
			services.Logger.Sugar().Error(err)
		} else {
			statistic.Messages = statsItem{Count: result.N}
		}
	}()

	wg.Wait()

	return &statistic, nil
}
