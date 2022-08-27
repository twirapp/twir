package types

type VariableHandlerParams struct {
	Key    string
	Params *string
}

type VariableHandlerResult struct {
	Result string
}

type VariableHandler func(data VariableHandlerParams) (*VariableHandlerResult, error)
type Variable struct {
	Name    string
	Handler VariableHandler
}

type Command struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	ChannelId    string   `json:"channel_id"`
	Aliases      []string `json:"aliases"`
	Responses    []string `json:"responses"`
	Permission   string   `json:"permission"`
	Description  *string  `json:"description"`
	Visible      bool     `json:"visible"`
	Module       string   `json:"module"`
	Enabled      bool     `json:"enabled"`
	Default      bool     `json:"default"`
	DefaultName  *bool    `json:"defaultName"`
	Cooldown     int      `json:"cooldown"`
	CooldownType string   `json:"cooldownType"`
}

type DefaultCommand struct {
	Command

	Handler func(data VariableHandlerParams)
}

var CommandPerms = []string{"BROADCASTER", "MODERATOR", "SUBSCRIBER", "VIP", "FOLLOWER", "VIEWER"}

type UserInfo struct {
	UserId          string   `json:"user_id"`
	UserName        *string  `json:"user_name"`
	UserDisplayName *string  `json:"user_display_name"`
	Badges          []string `json:"badges"`
}

type Channel struct {
	Id   string  `json:"id"`
	Name *string `json:"name"`
}

type ChatMessage struct {
	Channel Channel  `json:"channel"`
	Sender  UserInfo `json:"sender"`
	Text    string   `json:"text"`
}
