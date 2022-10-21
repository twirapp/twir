package greetings

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	uuid "github.com/satori/go.uuid"
)

func handleGet(channelId string, services types.Services) []model.ChannelsGreetings {
	greetings := []model.ChannelsGreetings{}
	services.DB.Where(`"channelId" = ?`, channelId).Find(&greetings)

	return greetings
}

func handlePost(
	channelId string,
	dto *greetingsDto,
	services types.Services,
) (*model.ChannelsGreetings, error) {
	twitchUser := getTwitchUserByName(dto.Username, services.Twitch)
	if twitchUser == nil {
		return nil, fiber.NewError(404, "cannot find twitch user")
	}

	existedGreeting := findGreetingByUser(twitchUser.ID, channelId, services.DB)
	if existedGreeting != nil {
		return nil, fiber.NewError(400, "greeting for that user already exists")
	}

	greeting := &model.ChannelsGreetings{
		ID:        uuid.NewV4().String(),
		ChannelID: channelId,
		UserID:    twitchUser.ID,
		Enabled:   *dto.Enabled,
		Text:      dto.Text,
		IsReply:   *dto.IsReply,
		Processed: false,
	}

	err := services.DB.Save(greeting).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot create greeting")
	}

	return greeting, nil
}

func handleDelete(greetingId string, services types.Services) error {
	greeting := findGreetingById(greetingId, services.DB)
	if greeting == nil {
		return fiber.NewError(404, "greeting not found")
	}
	err := services.DB.Where("id = ?", greetingId).Delete(&model.ChannelsGreetings{}).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "cannot delete greeting")
	}

	return nil
}

func handleUpdate(
	greetingId string,
	dto *greetingsDto,
	services types.Services,
) (*model.ChannelsGreetings, error) {
	greeting := findGreetingById(greetingId, services.DB)
	if greeting == nil {
		return nil, fiber.NewError(404, "greeting not found")
	}

	newGreeting := &model.ChannelsGreetings{
		ID:        greeting.ID,
		ChannelID: greeting.ChannelID,
		UserID:    greeting.UserID,
		Enabled:   *dto.Enabled,
		Text:      dto.Text,
		IsReply:   *dto.IsReply,
		Processed: false,
	}

	twitchUser := getTwitchUserByName(dto.Username, services.Twitch)
	if twitchUser == nil {
		return nil, fiber.NewError(404, "cannot find twitch user")
	}

	newGreeting.UserID = twitchUser.ID

	err := services.DB.Model(greeting).Select("*").Updates(newGreeting).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot update greeting")
	}

	return newGreeting, nil
}
