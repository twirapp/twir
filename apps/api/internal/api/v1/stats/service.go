package stats

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"gorm.io/gorm"
	"sync"

	model "github.com/satont/tsuwari/libs/gomodels"
)

type nResult struct {
	N int64
}

type statsItem struct {
	Count int64  `json:"count"`
	Name  string `json:"name"  enums:"users,channels,commands,messages"`
}

func handleGet() ([]statsItem, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)
	db := do.MustInvoke[*gorm.DB](di.Injector)

	wg := sync.WaitGroup{}
	statistic := []statsItem{
		{Name: "users"},
		{Name: "channels"},
		{Name: "commands"},
		{Name: "messages"},
	}

	wg.Add(4)

	go func() {
		defer wg.Done()
		var count int64
		err := db.Model(&model.Users{}).Count(&count).Error
		if err != nil {
			logger.Error(err)
		} else {
			statistic[0].Count = count
		}
	}()

	go func() {
		defer wg.Done()
		var count int64
		err := db.Model(&model.Channels{}).Count(&count).Error
		if err != nil {
			logger.Error(err)
		} else {
			statistic[1].Count = count
		}
	}()

	go func() {
		defer wg.Done()
		var count int64
		err := db.Model(&model.ChannelsCommands{}).
			Where("module = ?", "CUSTOM").
			Count(&count).
			Error
		if err != nil {
			logger.Error(err)
		} else {
			statistic[2].Count = count
		}
	}()

	go func() {
		defer wg.Done()
		result := nResult{}
		err := db.Model(&model.UsersStats{}).
			Select("sum(messages) as n").
			Scan(&result).
			Error
		if err != nil {
			logger.Error(err)
		} else {
			statistic[3].Count = result.N
		}
	}()

	wg.Wait()

	return statistic, nil
}
