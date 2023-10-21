package chat_client

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/satont/twir/libs/grpc/generated/events"

	"github.com/samber/lo"
)

func (c *ChatClient) onMessage(msg Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 510*time.Second)
	defer cancel()

	stream := model.ChannelsStreams{}
	if err := c.services.DB.Where(`"userId" = ?`, msg.Channel.ID).Find(&stream).Error; err != nil {
		c.services.Logger.Error(
			"cannot get stream",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
		)
		return
	}
	msg.DbStream = stream

	userBadges := createUserBadges(msg.User.Badges)
	// this need to be first because if we have no user in db it will produce many bugs
	user, err := c.updateUserStats(msg, userBadges)
	msg.DbUser = user
	if err != nil {
		c.services.Logger.Error(
			"cannot update user stats",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
			slog.String("userId", msg.User.ID),
		)
		return
	}

	var dbChannel model.Channels
	if err := c.services.DB.Where("id = ?", msg.Channel.ID).Find(&dbChannel).Error; err != nil {
		c.services.Logger.Error(
			"cannot get channel",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
		)
		return
	}
	msg.DbChannel = dbChannel

	if dbChannel.IsBotMod {
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
				c.handleModeration(ctx, msg)
			},
		)
	}

	c.workersPool.Submit(
		func() {
			c.handleTts(msg, userBadges)
		},
	)

	c.workersPool.Submit(
		func() {
			c.storeMessage(
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
			c.incrementStreamParsedMessages(msg.Channel.ID)
		},
	)

	c.workersPool.Submit(
		func() {
			c.removeUserFromLurkers(msg.User.ID)
		},
	)

	if msg.Tags["first-msg"] == "1" {
		c.workersPool.Submit(
			func() {
				if err := c.services.DB.Create(
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
					c.services.Logger.Error(
						"cannot save first user message",
						slog.Any("err", err),
						slog.String("channelId", msg.Channel.ID),
						slog.String("userId", msg.User.ID),
					)
				}
				_, err := c.services.EventsGrpc.FirstUserMessage(
					context.Background(), &events.FirstUserMessageMessage{
						BaseInfo:        &events.BaseInfo{ChannelId: msg.Channel.ID},
						UserId:          msg.User.ID,
						UserName:        msg.User.Name,
						UserDisplayName: msg.User.DisplayName,
					},
				)
				if err != nil {
					c.services.Logger.Error(
						"cannot fire first user message",
						slog.Any("err", err),
						slog.String("channelId", msg.Channel.ID),
						slog.String("userId", msg.User.ID),
					)
					return
				}
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
