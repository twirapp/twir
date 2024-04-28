package types

import (
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types/services"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

type ParseContextSender struct {
	ID          string
	Name        string
	DisplayName string
	Color       string
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
	MessageId string
	Channel   *ParseContextChannel
	Sender    *ParseContextSender
	Emotes    []*ParseContextEmote
	Mentions  []twitch.ChatMessageMessageFragmentMention

	Text      *string
	RawText   string
	IsCommand bool

	Command *model.ChannelsCommands

	Services *services.Services

	ArgsParser *command_arguments.Parser

	Cacher DataCacher
}
