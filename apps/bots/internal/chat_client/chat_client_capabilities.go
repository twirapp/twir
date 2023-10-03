package chat_client

import (
	irc "github.com/gempir/go-twitch-irc/v3"
)

var capabilities = []string{
	irc.TagsCapability,
	irc.MembershipCapability,
	irc.CommandsCapability,
}
