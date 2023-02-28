package types

import (
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
)

type VariableHandlerParams struct {
	Key    string
	Params *string
}

type VariableHandlerResult struct {
	Result string
}

type CommandsHandlerResult struct {
	Result  []string
	IsReply *bool
}

type (
	VariableHandler func(ctx *variables_cache.VariablesCacheService, data VariableHandlerParams) (*VariableHandlerResult, error)
	Variable        struct {
		Name         string
		Handler      VariableHandler
		Description  *string
		Example      *string
		CommandsOnly *bool
		Visible      *bool
	}
)

type DefaultCommand struct {
	*model.ChannelsCommands

	Handler func(ctx *variables_cache.ExecutionContext) *CommandsHandlerResult
}

type Sender struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Badges      []string `json:"badges"`
}

type Channel struct {
	Id   string  `json:"id"`
	Name *string `json:"name"`
}

type Message struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type HandleProcessCommandData struct {
	Channel Channel `json:"channel"`
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
}
