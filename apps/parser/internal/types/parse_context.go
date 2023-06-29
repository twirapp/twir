package types

import (
	"github.com/satont/twir/apps/parser/internal/types/services"
	model "github.com/satont/twir/libs/gomodels"
)

type ParseContextSender struct {
	ID          string
	Name        string
	DisplayName string
	Badges      []string
}

type ParseContextChannel struct {
	ID   string
	Name string
}

type ParseContextEmotePosition struct {
	Start int64
	End   int64
}

type ParseContextEmote struct {
	Name      string
	ID        string
	Count     int64
	Positions []*ParseContextEmotePosition
}

type ParseContext struct {
	Channel *ParseContextChannel
	Sender  *ParseContextSender
	Emotes  []*ParseContextEmote

	Text      *string
	IsCommand bool

	Command *model.ChannelsCommands

	Services *services.Services

	Cacher DataCacher
}
