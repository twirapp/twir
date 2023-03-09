package streams

import (
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (c *Streams) getService(channelId string) (*model.ChannelsStreams, error) {
	stream := model.ChannelsStreams{}
	err := c.services.Gorm.Where(`"userId" = ?`, channelId).First(&stream).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "stream not found")
	}
	return &stream, nil
}
