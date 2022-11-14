package handlers

import (
	"strings"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/bots/internal/bots/handlers/messages"
)

func (c *Handlers) OnPrivateMessage(msg irc.PrivateMessage) {
	userBadges := createUserBadges(msg.User.Badges)

	moderationResult := c.moderateMessage(msg, userBadges)
	if moderationResult {
		return
	}

	if strings.HasPrefix(msg.Message, "!") {
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
