package variables

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func handleGet(channelId string, services types.Services) ([]model.ChannelsCustomvars, error) {
	variables := []model.ChannelsCustomvars{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Find(&variables).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot get variables")
	}

	return variables, nil
}
