package api

import (
	"github.com/twirapp/twir/libs/types/types/api/bot"
	"github.com/twirapp/twir/libs/types/types/api/modules"
)

type Modules struct {
	OBS modules.OBS
	TTS modules.TTS
}

type Channels struct {
	MODULES Modules
	BOT     bot.Bot
}

type V1 struct {
	CHANNELS Channels
}
