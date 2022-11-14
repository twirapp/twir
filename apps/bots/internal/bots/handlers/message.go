package handlers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/bots/internal/bots/handlers/messages"
	"github.com/satont/tsuwari/libs/nats/parser"
	"gorm.io/gorm"
)

func (c *Handlers) OnPrivateMessage(msg irc.PrivateMessage) {
	userBadges := createUserBadges(msg.User.Badges)

	splittedMsg := strings.Split(msg.Message, " ")
	isReplyCommand := len(splittedMsg) >= 2 && strings.HasPrefix(splittedMsg[0], "@") &&
		strings.HasPrefix(splittedMsg[1], "!")

	if strings.HasPrefix(msg.Message, "!") || isReplyCommand {
		if isReplyCommand {
			msg.Message = strings.Join(splittedMsg[1:], " ")
		}

		go c.handleCommand(c.nats, msg, userBadges)
	}

	go c.handleGreetings(c.nats, c.db, msg, userBadges)
	go c.handleKeywords(c.nats, c.db, msg, userBadges)

	go func() {
		messages.IncrementUserMessages(c.db, msg.User.ID, msg.RoomID)
		messages.StoreMessage(
			c.db,
			msg.ID,
			msg.RoomID,
			msg.User.ID,
			msg.User.Name,
			msg.Message,
			!lo.Some(
				userBadges,
				[]string{"BROADCASTER", "MODERATOR", "SUBSCRIBER", "VIP"},
			),
		)
	}()
	go messages.IncrementStreamParsedMessages(c.db, msg.RoomID)
}

func (c *Handlers) handleCommand(nats *nats.Conn, msg irc.PrivateMessage, userBadges []string) {
	requestStruct := parser.Request{
		Sender: &parser.Sender{
			Id:          msg.User.ID,
			Name:        msg.User.Name,
			DisplayName: msg.User.DisplayName,
			Badges:      userBadges,
		},
		Channel: &parser.Channel{
			Id:   msg.RoomID,
			Name: msg.Channel,
		},
		Message: &parser.Message{
			Id:   msg.ID,
			Text: msg.Message,
		},
	}
	bytes, err := proto.Marshal(&requestStruct)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := nats.Request("parser.handleProcessCommand", bytes, 5*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseStruct := parser.Response{}
	if err = proto.Unmarshal(res.Data, &responseStruct); err != nil {
		return
	}

	if responseStruct.KeepOrder != nil && *responseStruct.KeepOrder {
		for _, v := range responseStruct.Responses {
			r := v
			go c.BotClient.SayWithRateLimiting(
				msg.Channel,
				r,
				lo.If(responseStruct.IsReply, lo.ToPtr(msg.ID)).Else(nil),
			)
		}
	} else {
		for _, r := range responseStruct.Responses {
			c.BotClient.SayWithRateLimiting(
				msg.Channel,
				r,
				lo.If(responseStruct.IsReply, lo.ToPtr(msg.ID)).Else(nil),
			)
		}
	}

	return
}

func createUserBadges(badges map[string]int) []string {
	userBadges := lo.MapToSlice(badges, func(k string, _ int) string {
		return strings.ToUpper(k)
	})

	if _, ok := badges["founder"]; ok {
		userBadges = append(userBadges, "SUBSCRIBER")
	}

	userBadges = append(userBadges, "VIEWER")

	return userBadges
}

func (c *Handlers) handleGreetings(
	nats *nats.Conn,
	db *gorm.DB,
	msg irc.PrivateMessage,
	userBadges []string,
) {
	entity := model.ChannelsGreetings{}
	err := db.
		Where(`"channelId" = ? AND "userId" = ? AND processed = ?`, msg.RoomID, msg.User.ID, false).
		First(&entity).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	requestStruct := parser.ParseResponseRequest{
		Channel: &parser.Channel{
			Id:   msg.RoomID,
			Name: msg.Channel,
		},
		Message: &parser.Message{
			Text: entity.Text,
			Id:   msg.ID,
		},
		Sender: &parser.Sender{
			Id:          msg.User.ID,
			Name:        msg.User.Name,
			DisplayName: msg.User.DisplayName,
			Badges:      userBadges,
		},
		ParseVariables: lo.ToPtr(true),
	}
	bytes, err := proto.Marshal(&requestStruct)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseStruct := parser.ParseResponseResponse{}
	res, err := nats.Request("parser.parseTextResponse", bytes, 5*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := proto.Unmarshal(res.Data, &responseStruct); err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range responseStruct.Responses {
		c.BotClient.SayWithRateLimiting(
			msg.Channel,
			r,
			lo.If(entity.IsReply, &msg.ID).Else(nil),
		)
	}

	db.Model(&entity).Where("id = ?", entity.ID).Select("*").Updates(map[string]any{
		"processed": true,
	})
}

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
