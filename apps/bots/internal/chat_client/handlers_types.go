package chat_client

import "github.com/gempir/go-twitch-irc/v3"

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

type Message struct {
	ID      string
	Channel MessageChannel
	User    MessageUser
	Message string
	Emotes  []*twitch.Emote
	Tags    map[string]string
}
