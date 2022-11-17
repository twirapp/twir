package stats

import (
	"sync"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/tsuwari/apps/api/internal/types"
)

type nResult struct {
	N int64
}

type statsItem struct {
	Count int64  `json:"count"`
	Name  string `json:"name"`
}

func handleGet(services types.Services) ([]statsItem, error) {
	wg := sync.WaitGroup{}
	statistic := []statsItem{}
	mu := sync.Mutex{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		var count int64
		err := services.DB.Model(&model.Users{}).Count(&count).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
		} else {
			mu.Lock()
			defer mu.Unlock()
			statistic = append(statistic, statsItem{Count: count, Name: "users"})
		}
	}()

	go func() {
		defer wg.Done()
		var count int64
		err := services.DB.Model(&model.Channels{}).Count(&count).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
		} else {
			mu.Lock()
			defer mu.Unlock()
			statistic = append(statistic, statsItem{Count: count, Name: "channels"})
		}
	}()

	go func() {
		defer wg.Done()
		var count int64
		err := services.DB.Model(&model.ChannelsCommands{}).Count(&count).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
		} else {
			mu.Lock()
			defer mu.Unlock()
			statistic = append(statistic, statsItem{Name: "commands", Count: count})
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
			mu.Lock()
			defer mu.Unlock()
			statistic = append(statistic, statsItem{Name: "messages", Count: result.N})
		}
	}()

	wg.Wait()

	return statistic, nil
}
