package handlers

import (
	"context"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/satont/twir/libs/grpc/generated/events"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/bots/handlers/messages"
)

func (c *Handlers) OnMessage(msg *Message) {
	userBadges := createUserBadges(msg.User.Badges)
	// this need to be first because if we have no user in db it will produce many bugs
	messages.IncrementUserMessages(c.db, msg.User.ID, msg.Channel.ID)

	c.workersPool.Submit(
		func() {
			splittedMsg := strings.Split(msg.Message, " ")
			isReplyCommand := len(splittedMsg) >= 2 && strings.HasPrefix(splittedMsg[0], "@") &&
				strings.HasPrefix(splittedMsg[1], "!")

			if strings.HasPrefix(msg.Message, "!") || isReplyCommand {
				if isReplyCommand {
					msg.Message = strings.Join(splittedMsg[1:], " ")
				}

				c.handleCommand(msg, userBadges)
			}
		},
	)

	c.workersPool.Submit(
		func() {
			c.handleGreetings(msg, userBadges)
		},
	)

	c.workersPool.Submit(
		func() {
			c.handleKeywords(msg, userBadges)
		},
	)

	c.workersPool.Submit(
		func() {
			c.handleEmotes(msg)
		},
	)

	c.workersPool.Submit(
		func() {
			c.handleTts(msg, userBadges)
		},
	)

	c.workersPool.Submit(
		func() {
			messages.StoreMessage(
				c.db,
				msg.ID,
				msg.Channel.ID,
				msg.User.ID,
				msg.User.Name,
				msg.Message,
				!lo.Some(
					userBadges,
					[]string{"BROADCASTER", "MODERATOR", "SUBSCRIBER", "VIP"},
				),
			)
		},
	)

	c.workersPool.Submit(
		func() {
			messages.IncrementStreamParsedMessages(c.db, msg.Channel.ID)
		},
	)

	c.workersPool.Submit(
		func() {
			messages.RemoveUserFromLurkers(c.db, msg.User.ID)
		},
	)

	if msg.Tags["first-msg"] == "1" {
		c.workersPool.Submit(
			func() {
				if err := c.db.Create(
					&model.ChannelsEventsListItem{
						ID:        uuid.New().String(),
						ChannelID: msg.Channel.ID,
						UserID:    msg.User.ID,
						Type:      model.ChannelEventListItemTypeFirstUserMessage,
						CreatedAt: time.Now().UTC(),
						Data: &model.ChannelsEventsListItemData{
							FirstUserMessageUserName:        msg.User.Name,
							FirstUserMessageUserDisplayName: msg.User.DisplayName,
							FirstUserMessageMessage:         msg.Message,
						},
					},
				).Error; err != nil {
					zap.S().Error(err)
				}
				c.eventsGrpc.FirstUserMessage(
					context.Background(), &events.FirstUserMessageMessage{
						BaseInfo:        &events.BaseInfo{ChannelId: msg.Channel.ID},
						UserId:          msg.User.ID,
						UserName:        msg.User.Name,
						UserDisplayName: msg.User.DisplayName,
					},
				)
			},
		)
	}
}

func createUserBadges(badges map[string]int) []string {
	userBadges := lo.MapToSlice(
		badges, func(k string, _ int) string {
			return strings.ToUpper(k)
		},
	)

	if _, ok := badges["founder"]; ok {
		userBadges = append(userBadges, "SUBSCRIBER")
	}

	return userBadges
}
