package greetings

import (
	"errors"
	"fmt"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	uuid "github.com/satori/go.uuid"
)

func HandleGet(channelId string, services types.Services) []model.ChannelsGreetings {
	greetings := []model.ChannelsGreetings{}
	services.DB.Where(`"channelId" = ?`, channelId).Find(&greetings)

	return greetings
}

func HandlePost(
	channelId string,
	dto *greetingsDto,
	services types.Services,
) (*model.ChannelsGreetings, error) {
	twitchUser := getTwitchUserByName(dto.Username, services.Twitch)
	if twitchUser == nil {
		return nil, errors.New("cannot find twitch user")
	}

	existedGreeting := findGreetingByUser(twitchUser.ID, channelId, services.DB)
	if existedGreeting != nil {
		return nil, errors.New("greeting for that user already exists")
	}

	fmt.Println(*dto.Enabled, *dto.IsReply)
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
		fmt.Println(err)
		return nil, errors.New("cannot create greeting")
	}

	return greeting, nil
}

func HandleDelete(greetingId string, services types.Services) error {
	greeting := model.ChannelsGreetings{}
	err := services.DB.Where("id = ?", greetingId).First(&greeting).Error
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(404, "greeting not found")
	}
	err = services.DB.Where("id = ?", greetingId).Delete(&model.ChannelsGreetings{}).Error
	fmt.Println(err)
	if err != nil {
		return errors.New("cannot delete greeting")
	}

	return nil
}
