package chat_client

import (
	irc "github.com/gempir/go-twitch-irc/v3"
	model "github.com/satont/twir/libs/gomodels"
)

type MessageUser struct {
	ID          string
	Name        string
	DisplayName string
	Badges      map[string]int
}

type MessageChannel struct {
	ID   string
	Name string
}

type EmotePosition struct {
	Start int
	End   int
}

type MessageEmote struct {
	Name      string
	ID        string
	Count     int
	Positions []EmotePosition
}

type Message struct {
	ID      string
	Channel MessageChannel
	User    MessageUser
	Message string
	Emotes  []*irc.Emote
	Tags    map[string]string

	DbChannel model.Channels
	DbUser    model.Users
	DbStream  model.ChannelsStreams
}
