package types

import variablescache "tsuwari/parser/internal/variablescache"

type VariableHandlerParams struct {
	Key    string
	Params *string
}

type VariableHandlerResult struct {
	Result string
}

type VariableHandler func(ctx *variablescache.VariablesCacheService, data VariableHandlerParams) (*VariableHandlerResult, error)
type Variable struct {
	Name    string
	Handler VariableHandler
}

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
}

type DefaultCommand struct {
	Command

	Handler func(data VariableHandlerParams) []string
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
