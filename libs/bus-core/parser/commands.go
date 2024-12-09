package parser

type CommandParseResponse struct {
	Responses []string
	IsReply   bool
	KeepOrder bool
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
