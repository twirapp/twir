// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlmodel

import (
	"github.com/99designs/gqlgen/graphql"
)

type Command struct {
	ID                        string            `json:"id"`
	Name                      string            `json:"name"`
	Description               *string           `json:"description,omitempty"`
	Aliases                   []string          `json:"aliases,omitempty"`
	Responses                 []CommandResponse `json:"responses,omitempty"`
	Cooldown                  *int              `json:"cooldown,omitempty"`
	CooldownType              string            `json:"cooldownType"`
	Enabled                   bool              `json:"enabled"`
	Visible                   bool              `json:"visible"`
	Default                   bool              `json:"default"`
	DefaultName               *string           `json:"defaultName,omitempty"`
	Module                    string            `json:"module"`
	IsReply                   bool              `json:"isReply"`
	KeepResponsesOrder        bool              `json:"keepResponsesOrder"`
	DeniedUsersIds            []string          `json:"deniedUsersIds,omitempty"`
	AllowedUsersIds           []string          `json:"allowedUsersIds,omitempty"`
	RolesIds                  []string          `json:"rolesIds,omitempty"`
	OnlineOnly                bool              `json:"onlineOnly"`
	CooldownRolesIds          []string          `json:"cooldownRolesIds,omitempty"`
	EnabledCategories         []string          `json:"enabledCategories,omitempty"`
	RequiredWatchTime         int               `json:"requiredWatchTime"`
	RequiredMessages          int               `json:"requiredMessages"`
	RequiredUsedChannelPoints int               `json:"requiredUsedChannelPoints"`
}

type CommandResponse struct {
	ID        string `json:"id"`
	CommandID string `json:"commandId"`
	Text      string `json:"text"`
	Order     int    `json:"order"`
}

type CreateCommandInput struct {
	Name        string                                          `json:"name"`
	Description graphql.Omittable[*string]                      `json:"description,omitempty"`
	Aliases     graphql.Omittable[[]string]                     `json:"aliases,omitempty"`
	Responses   graphql.Omittable[[]CreateCommandResponseInput] `json:"responses,omitempty"`
}

type CreateCommandResponseInput struct {
	Text  string `json:"text"`
	Order int    `json:"order"`
}

type Mutation struct {
}

type Notification struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Text   string `json:"text"`
}

type Query struct {
}

type Subscription struct {
}

type UpdateCommandOpts struct {
	Name        graphql.Omittable[*string]  `json:"name,omitempty"`
	Description graphql.Omittable[*string]  `json:"description,omitempty"`
	Aliases     graphql.Omittable[[]string] `json:"aliases,omitempty"`
}

type User struct {
	ID                string       `json:"id"`
	IsBotAdmin        bool         `json:"isBotAdmin"`
	APIKey            string       `json:"apiKey"`
	IsBanned          bool         `json:"isBanned"`
	HideOnLandingPage bool         `json:"hideOnLandingPage"`
	Channel           *UserChannel `json:"channel"`
}

type UserChannel struct {
	IsEnabled      bool   `json:"isEnabled"`
	IsBotModerator bool   `json:"isBotModerator"`
	BotID          string `json:"botId"`
}