package greetings

import (
	"net/http"
	"sync"

	"github.com/satont/tsuwari/libs/twitch"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	uuid "github.com/satori/go.uuid"
)

type Greeting struct {
	model.ChannelsGreetings
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
}

func (c *Greetings) getService(channelId string) []Greeting {
	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		return nil
	}

	greetings := []model.ChannelsGreetings{}
	c.services.Gorm.Where(`"channelId" = ?`, channelId).Find(&greetings)
	users := []Greeting{}

	chunks := lo.Chunk(greetings, 100)
	greetingsWg := sync.WaitGroup{}
	greetingsWg.Add(len(chunks))
	for _, chunk := range chunks {
		go func(c []model.ChannelsGreetings) {
			defer greetingsWg.Done()
			ids := lo.Map(chunk, func(g model.ChannelsGreetings, _ int) string {
				return g.UserID
			})
			twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
				IDs: ids,
			})
			if err != nil {
				return
			}
			for _, u := range twitchUsers.Data.Users {
				user, ok := lo.Find(greetings, func(g model.ChannelsGreetings) bool {
					return g.UserID == u.ID
				})
				if ok {
					users = append(users, Greeting{
						ChannelsGreetings: user,
						UserName:          u.Login,
						Avatar:            u.ProfileImageURL,
					})
				}
			}
		}(chunk)
	}

	greetingsWg.Wait()

	return users
}

func (c *Greetings) postService(
	channelId string,
	dto *greetingsDto,
) (*Greeting, error) {
	twitchUser := c.getTwitchUserByName(dto.Username)
	if twitchUser == nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find twitch user")
	}

	existedGreeting := c.findGreetingByUser(twitchUser.ID, channelId)
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

	err := c.services.Gorm.Save(greeting).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create greeting")
	}

	return &Greeting{
		ChannelsGreetings: *greeting,
		UserName:          twitchUser.Login,
		Avatar:            twitchUser.ProfileImageURL,
	}, nil
}

func (c *Greetings) deleteService(greetingId string) error {
	greeting := c.findGreetingById(greetingId)
	if greeting == nil {
		return fiber.NewError(http.StatusNotFound, "greeting not found")
	}
	err := c.services.Gorm.Where("id = ?", greetingId).Delete(&model.ChannelsGreetings{}).Error
	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete greeting")
	}

	return nil
}

func (c *Greetings) putService(
	greetingId string,
	dto *greetingsDto,
) (*Greeting, error) {
	greeting := c.findGreetingById(greetingId)
	if greeting == nil {
		return nil, fiber.NewError(http.StatusNotFound, "greeting not found")
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

	twitchUser := c.getTwitchUserByName(dto.Username)
	if twitchUser == nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find twitch user")
	}

	newGreeting.UserID = twitchUser.ID

	err := c.services.Gorm.Model(greeting).Select("*").Updates(newGreeting).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update greeting")
	}

	return &Greeting{
		ChannelsGreetings: *newGreeting,
		UserName:          twitchUser.Login,
		Avatar:            twitchUser.ProfileImageURL,
	}, nil
}

func (c *Greetings) patchService(
	channelId,
	greetingId string,
	dto *greetingsPatchDto,
) (*Greeting, error) {
	greeting := c.findGreetingById(greetingId)
	if greeting == nil {
		return nil, fiber.NewError(http.StatusNotFound, "greeting not found")
	}

	twitchUser := c.getTwitchUserById(greeting.UserID)
	if twitchUser == nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find twitch user")
	}

	greeting.Enabled = *dto.Enabled
	err := c.services.Gorm.Model(greeting).Select("*").Updates(greeting).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update greeting")
	}

	return &Greeting{
		ChannelsGreetings: *greeting,
		UserName:          twitchUser.Login,
		Avatar:            twitchUser.ProfileImageURL,
	}, nil
}
