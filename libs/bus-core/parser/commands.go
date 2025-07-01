package parser

type CommandParseResponse struct {
	Responses         []string
	IsReply           bool
	KeepOrder         bool
	SkipToxicityCheck bool
}

type ParseVariablesInTextRequest struct {
	ChannelID     string
	ChannelName   string
	Text          string
	UserID        string
	UserLogin     string
	UserName      string
	IsCommand     bool
	IsInCustomVar bool
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
