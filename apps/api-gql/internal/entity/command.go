package entity

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	ID                        uuid.UUID
	Name                      string
	Cooldown                  *int
	CooldownType              string
	Enabled                   bool
	Aliases                   []string
	Description               *string
	Visible                   bool
	ChannelID                 string
	Default                   bool
	DefaultName               *string
	Module                    string
	IsReply                   bool
	KeepResponsesOrder        bool
	DeniedUsersIDS            []string
	AllowedUsersIDS           []string
	RolesIDS                  []string
	OnlineOnly                bool
	CooldownRolesIDs          []string
	EnabledCategories         []string
	RequiredWatchTime         int
	RequiredMessages          int
	RequiredUsedChannelPoints int
	GroupID                   *uuid.UUID
	ExpiresAt                 *time.Time
	ExpiresType               *CommandExpireType
}

type CommandExpireType string

const (
	CommandExpireTypeDisable CommandExpireType = "DISABLE"
	CommandExpireTypeDelete  CommandExpireType = "DELETE"
)

type CommandGroup struct {
	ID        uuid.UUID
	ChannelID string
	Name      string
	Color     string
}

type CommandResponse struct {
	ID                uuid.UUID
	Text              *string
	CommandID         uuid.UUID
	Order             int
	TwitchCategoryIDs []string
}

type CommandWithGroupAndResponses struct {
	Command   Command
	Group     *CommandGroup
	Responses []CommandResponse
}
