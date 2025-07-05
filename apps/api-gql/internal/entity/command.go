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
	RolesIDS                  []uuid.UUID
	OnlineOnly                bool
	OfflineOnly               bool
	CooldownRolesIDs          []string
	EnabledCategories         []string
	RequiredWatchTime         int
	RequiredMessages          int
	RequiredUsedChannelPoints int
	GroupID                   *uuid.UUID
	ExpiresAt                 *time.Time
	ExpiresType               *CommandExpireType
}

var CommandNil = Command{}

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

var CommandGroupNil = CommandGroup{}

type CommandResponse struct {
	ID                uuid.UUID
	Text              *string
	CommandID         uuid.UUID
	Order             int
	TwitchCategoryIDs []string
	OnlineOnly        bool
	OfflineOnly       bool
}

var CommandResponseNil = CommandResponse{}

type CommandWithGroupAndResponses struct {
	Command   Command
	Group     *CommandGroup
	Responses []CommandResponse
}

var CommandWithGroupAndResponsesNil = CommandWithGroupAndResponses{}
