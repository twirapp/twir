package gamesAliases

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

func handleGetGamesAliases(channelId string, services types.Services) ([]model.ChannelCategoryAlias, error) {
	db := do.MustInvoke[*gorm.DB](di.Provider)

	var aliases []model.ChannelCategoryAlias
	err := db.Where(`"channelId" = ?`, channelId).Find(&aliases).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return aliases, nil
}
