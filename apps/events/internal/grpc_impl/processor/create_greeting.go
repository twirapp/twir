package processor

import (
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func (c *Processor) CreateGreeting() {
	if c.data.RewardInput == nil {
		return
	}

	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{c.data.UserName},
	})

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			c.services.Logger.Sugar().Error(err)
		}
		return
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
	c.services.DB.Create(&newGreeting)
}
