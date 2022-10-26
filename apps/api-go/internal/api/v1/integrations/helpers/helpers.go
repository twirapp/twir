package helpers

import (
	"fmt"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetIntegration(
	channelId string,
	service string,
	db *gorm.DB,
) (*model.ChannelsIntegrations, error) {
	integration := model.ChannelsIntegrations{}

	err := db.
		Preload("Integration").
		Joins(`JOIN integrations i on i.id = channels_integrations."integrationId"`).
		Where(`"channels_integrations"."channelId" = ? AND i.service = ?`, channelId, service).
		First(&integration).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		fmt.Println(err)
		return nil, fiber.NewError(404, "internal error")
	}

	return &integration, nil
}
