package model

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
	ExpiresType               *ExpireType
}

var Nil = Command{}

type ExpireType string

const (
	ExpireTypeDisable ExpireType = "DISABLE"
	ExpireTypeDelete  ExpireType = "DELETE"
)
