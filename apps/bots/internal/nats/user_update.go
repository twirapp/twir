package nats_handlers

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	model "github.com/satont/tsuwari/libs/gomodels"
)

type user struct {
	UserID        string `json:"user_id"`
	UserLogin     string `json:"user_login"`
	UserName      string `json:"user_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Description   string `json:"description"`
}

type incoming struct {
	Pattern string `json:"pattern"`
	User    user   `json:"data"`
}

func (c *NatsHandlers) UserUpdate(m *nats.Msg) {
	data := incoming{}
	if err := json.Unmarshal(m.Data, &data); err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	channel := model.Channels{}
	if err := c.db.Where("id = ?", data.User.UserID).Find(&channel).Error; err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	if channel.ID == "" {
		return
	}

	bot, isBotFound := c.botsService.Instances[channel.BotID]
	if !isBotFound {
		return
	}

	if channel.IsEnabled {
		bot.Join(data.User.UserName)
	} else {
		bot.Depart(data.User.UserName)
	}

	fmt.Printf("%+v\n", data)
}
