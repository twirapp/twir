package entity

type Channel struct {
	ID             string
	IsEnabled      bool
	IsTwitchBanned bool
	IsBotMod       bool
	BotID          string
}

var ChannelNil = Channel{}
