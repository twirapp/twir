package parser

import (
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/generic"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type CommandParseResponse struct {
	Responses         []string
	IsReply           bool
	KeepOrder         bool
	SkipToxicityCheck bool
}

type ParseVariablesInTextRequest struct {
	ChannelID      uuid.UUID
	ChannelName    string
	Text           string
	UserID         string
	UserLogin      string
	UserName       string
	IsCommand      bool
	IsInCustomVar  bool
	Mentions       []generic.ChatMessageMessageFragmentMention
	PlatformSource *platformentity.Platform
}

type ParseVariablesInTextResponse struct {
	Text string
}

const DefaultCommandsSubject = "parser.getDefaultCommands"

type GetDefaultCommandsResponse struct {
	List []DefaultCommand
}

type DefaultCommand struct {
	Name               string
	Description        string
	Visible            bool
	RolesNames         []string
	Module             string
	IsReply            bool
	KeepResponsesOrder bool
	Aliases            []string
}
