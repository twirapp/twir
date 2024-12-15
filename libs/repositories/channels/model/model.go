package model

type Channel struct {
	ID             string
	IsEnabled      bool
	IsTwitchBanned bool
	IsBotMod       bool
	BotID          string
}

var Nil = Channel{}
