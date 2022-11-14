package handlers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/nats/parser"
	"gorm.io/gorm"
)

func (c *Handlers) handleKeywords(
	nats *nats.Conn,
	db *gorm.DB,
	msg irc.PrivateMessage,
	userBadges []string,
) {
	keywords := []model.ChannelsKeywords{}
	err := db.Where(`"channelId" = ?`, msg.RoomID).Find(&keywords).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(keywords) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(keywords))

	message := strings.ToLower(msg.Message)
	for _, k := range keywords {
		go func(k model.ChannelsKeywords) {
			defer wg.Done()
			if !strings.Contains(message, strings.ToLower(k.Text)) {
				return
			}

			isOnCooldown := false
			if k.Cooldown.Valid && k.CooldownExpireAt.Valid {
				isOnCooldown = k.CooldownExpireAt.Time.After(time.Now())
			}
			if isOnCooldown {
				return
			}

			requestStruct := parser.ParseResponseRequest{
				Channel: &parser.Channel{
					Id:   msg.RoomID,
					Name: msg.Channel,
				},
				Sender: &parser.Sender{
					Id:          msg.User.ID,
					Name:        msg.User.Name,
					DisplayName: msg.User.DisplayName,
					Badges:      userBadges,
				},
				Message: &parser.Message{
					Text: k.Response,
					Id:   msg.ID,
				},
				ParseVariables: lo.ToPtr(true),
			}

			bytes, err := proto.Marshal(&requestStruct)
			if err != nil {
				fmt.Println(err)
				return
			}

			req, err := nats.Request("parser.parseTextResponse", bytes, 5*time.Second)
			if err != nil {
				fmt.Println(err)
				return
			}

			responseStruct := parser.ParseResponseResponse{}
			if err := proto.Unmarshal(req.Data, &responseStruct); err != nil {
				fmt.Println(err)
				return
			}

			for _, r := range responseStruct.Responses {
				c.BotClient.SayWithRateLimiting(
					msg.Channel,
					r,
					lo.If(k.IsReply, &msg.ID).Else(nil),
				)
			}

			db.Model(&k).Where("id = ?", k.ID).Select("*").Updates(map[string]any{
				"cooldownExpireAt": time.Now().Add(time.Duration(k.Cooldown.Int64) * time.Second),
			})
		}(k)
	}

	wg.Wait()
}
