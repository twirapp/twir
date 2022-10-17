package types

import (
	usersauth "tsuwari/parser/internal/twitch/user"
	variables_cache "tsuwari/parser/internal/variablescache"
)

type VariableHandlerParamsServices struct {
	UsersAuth *usersauth.UsersTokensService
}
type VariableHandlerParams struct {
	Key      string
	Params   *string
	Services VariableHandlerParamsServices
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
		Name        string
		Handler     VariableHandler
		Description *string
		Example     *string
	}
)

type Command struct {
	Id           *string  `json:"id"`
	Name         string   `json:"name"`
	ChannelId    string   `json:"channel_id"`
	Aliases      []string `json:"aliases"`
	Responses    []string `json:"responses"`
	Permission   string   `json:"permission"`
	Description  *string  `json:"description"`
	Visible      bool     `json:"visible"`
	Module       *string  `json:"module"`
	Enabled      bool     `json:"enabled"`
	Default      bool     `json:"default"`
	DefaultName  *string  `json:"defaultName"`
	Cooldown     int      `json:"cooldown"`
	CooldownType string   `json:"cooldownType"`
	IsReply      bool     `json:"isReply"`
	KeepOrder    *int     `json:"keepOrder"`
}

type DefaultCommand struct {
	Command

	Handler func(ctx variables_cache.ExecutionContext) *CommandsHandlerResult
	IsReply *bool
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
