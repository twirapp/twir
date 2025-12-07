package model

type Channel struct {
	ID             string
	IsEnabled      bool
	IsTwitchBanned bool
	IsBotMod       bool
	BotID          string

	isNil bool
}

func (c Channel) IsNil() bool {
	return c.isNil
}

var Nil = Channel{
	isNil: true,
}
