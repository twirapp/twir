package parser

const PARSER_COMMANDS_QUEUE = "parser.commands_queue"

type CommandParseResponse struct {
	Responses []string
	IsReply   bool
	KeepOrder bool
}
