package handlers

import (
	"github.com/gempir/go-twitch-irc/v3"
)

func (c *Handlers) OnSelfJoin(msg twitch.UserJoinMessage) {
	c.logger.Sugar().
		Infow("Joined channel", "bot", msg.User, "channel", msg.Channel)
}
