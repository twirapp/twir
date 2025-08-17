package types

import (
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
)

type ParseContextSender struct {
	ID               string
	Name             string
	DisplayName      string
	Color            string
	Badges           []string
	Roles            []model.ChannelRole
	UserChannelStats *model.UsersStats
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
	Positions []*ParseContextEmotePosition
	Count     int64
}

type ParseContext struct {
	Cacher  DataCacher
	Channel *ParseContextChannel
	Sender  *ParseContextSender

	Text *string

	Command       *model.ChannelsCommands
	ChannelStream *streamsmodel.Stream

	Services *services.Services

	ArgsParser *command_arguments.Parser

	MessageId string
	RawText   string
	Emotes    []*ParseContextEmote
	Mentions  []twitch.ChatMessageMessageFragmentMention

	IsCommand     bool
	IsInCustomVar bool
}
