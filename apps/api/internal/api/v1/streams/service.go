package streams

import (
	"net/http"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"gorm.io/gorm"
)

func handleGet(channelId string, services types.Services) (*model.ChannelsStreams, error) {
	stream := model.ChannelsStreams{}
	err := services.DB.Where(`"userId" = ?`, channelId).First(&stream).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "stream not found")
	}
	return &stream, nil
}
