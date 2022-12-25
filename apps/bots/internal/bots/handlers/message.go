package handlers

import (
	"strings"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/bots/internal/bots/handlers/messages"
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

		go c.handleCommand(msg, userBadges)
	}

	go c.handleGreetings(msg, userBadges)
	go c.handleKeywords(msg, userBadges)

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
