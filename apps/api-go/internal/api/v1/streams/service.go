package streams

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func handleGet(channelId string, services types.Services) (*model.ChannelsStreams, error) {
	stream := model.ChannelsStreams{}
	err := services.DB.Where(`"user_id" = ?`, channelId).First(&stream).Error
	if err != nil {
		return nil, fiber.NewError(404, "stream not found")
	}
	return &stream, nil
}
