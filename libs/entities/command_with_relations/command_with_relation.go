package commandwithrelationentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/commandrolecooldownentity"
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
	EnabledCategories         []string
	RequiredWatchTime         int
	RequiredMessages          int
	RequiredUsedChannelPoints int
	GroupID                   *uuid.UUID
	ExpiresAt                 *time.Time
	ExpiresType               *CommandExpireType
	RolesCooldowns            []commandrolecooldownentity.CommandRoleCooldown

	isNil bool
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

	isNil bool
}

func (c CommandGroup) IsNil() bool {
	return c.isNil
}

var CommandGroupNil = CommandGroup{
	isNil: true,
}

type CommandResponse struct {
	ID                uuid.UUID
	Text              *string
	CommandID         uuid.UUID
	Order             int
	TwitchCategoryIDs []string
	OnlineOnly        bool
	OfflineOnly       bool

	isNil bool
}

func (c CommandResponse) IsNil() bool {
	return c.isNil
}

var CommandResponseNil = CommandResponse{
	isNil: true,
}

type CommandWithGroupAndResponses struct {
	Command        Command
	Group          *CommandGroup
	Responses      []CommandResponse
	RolesCooldowns []commandrolecooldownentity.CommandRoleCooldown

	isNil bool
}

func (c CommandWithGroupAndResponses) IsNil() bool {
	return c.isNil
}

var CommandWithGroupAndResponsesNil = CommandWithGroupAndResponses{
	isNil: true,
}
