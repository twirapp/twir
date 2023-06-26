package processor

import (
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func (c *Processor) CreateGreeting() error {
	if c.data.RewardInput == nil {
		return InternalError
	}

	user, err := c.streamerApiClient.GetUsers(
		&helix.UsersParams{
			Logins: []string{c.data.UserName},
		},
	)
	if err != nil {
		return err
	}

	if len(user.Data.Users) == 0 {
		return InternalError
	}

	newGreeting := model.ChannelsGreetings{
		ID:        uuid.NewV4().String(),
		ChannelID: c.channelId,
		UserID:    user.Data.Users[0].ID,
		Enabled:   true,
		Text:      *c.data.RewardInput,
		IsReply:   true,
		Processed: false,
	}
	err = c.services.DB.Create(&newGreeting).Error
	if err != nil {
		return err
	}

	return nil
}
