package commands

import (
	"time"

	"github.com/google/uuid"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"go.uber.org/fx"
)

var FxModule = fx.Provide(
	httpbase.AsFxRoute(newListById),
)

type commandResponseDto struct {
	ID                        uuid.UUID                        `json:"id"`
	Name                      string                           `json:"name" description:"The name of the command, without the prefix"`
	Cooldown                  *int                             `json:"cooldown" description:"The cooldown time in seconds"`
	CooldownType              string                           `json:"cooldown_type" enum:"GLOBAL,PER_USER"`
	Enabled                   bool                             `json:"enabled"`
	Aliases                   []string                         `json:"aliases" description:"List of alternative names for the command"`
	Description               *string                          `json:"description" description:"A brief description of what the command does"`
	Visible                   bool                             `json:"visible" description:"Whether the command is visible in the command list"`
	IsDefault                 bool                             `json:"is_default" description:"Whether the command is a default command provided by the system"`
	DefaultName               *string                          `json:"default_name" description:"The original name of the default command, if applicable"`
	Module                    string                           `json:"module" description:"The module or category the command belongs to"`
	IsReply                   bool                             `json:"is_reply" description:"Whether the command responses are replies to the user"`
	KeepResponsesOrder        bool                             `json:"keep_responses_order" description:"Whether to keep the order of responses as defined"`
	DeniedUsersIDS            []string                         `json:"denied_users_ids" description:"List of user IDs who are denied access to the command"`
	AllowedUsersIDS           []string                         `json:"allowed_users_ids" description:"List of user IDs who are allowed access to the command"`
	RolesIDS                  []uuid.UUID                      `json:"roles_ids" format:"uuid" description:"List of role IDs that have access to the command"`
	OnlineOnly                bool                             `json:"online_only" description:"Whether the command is only usable when the stream is online"`
	OfflineOnly               bool                             `json:"offline_only" description:"Whether the command is only usable when the stream is offline"`
	EnabledCategories         []string                         `json:"enabled_categories" description:"List of categories where the command is enabled"`
	RequiredWatchTime         int                              `json:"required_watch_time" description:"The required watch time in minutes to use the command"`
	RequiredMessages          int                              `json:"required_messages" description:"The required number of messages sent to use the command"`
	RequiredUsedChannelPoints int                              `json:"required_used_channel_points" description:"The required amount of channel points used to use the command"`
	Expire                    *Expire                          `json:"expire" description:"Information about when the command expires and what happens upon expiration"`
	Responses                 []commandResponsesResponseDto    `json:"responses" json:"Responses" description:"List of possible responses for the command"`
	Group                     *commandGroupResponseDto         `json:"group" description:"The group this command belongs to, if any"`
	RolesCooldowns            []commandRoleCooldownResponseDto `json:"roles_cooldowns" description:"List of role-specific cooldowns for the command"`
}

type commandGroupResponseDto struct {
	ID    uuid.UUID `json:"id" format:"uuid"`
	Name  string    `json:"name"`
	Color string    `json:"color"`
}

type commandResponsesResponseDto struct {
	ID                uuid.UUID `json:"id" format:"uuid"`
	Text              string    `json:"text"`
	Order             int       `json:"order" description:"The order of the response in the list of responses"`
	TwitchCategoryIDs []string  `json:"twitch_category_id" description:"List of Twitch category IDs where this response is applicable"`
	OnlineOnly        bool      `json:"online_only" description:"Whether this response is only applicable when the stream is online"`
	OfflineOnly       bool      `json:"offline_only" description:"Whether this response is only applicable when the stream is offline"`
}

type commandRoleCooldownResponseDto struct {
	RoleID   uuid.UUID `json:"role_id" description:"The ID of the role"`
	Cooldown int       `json:"cooldown" description:"The cooldown time in seconds for this role"`
}

type Expire struct {
	ExpiresAt   time.Time `json:"expires_at" format:"date-time"`
	ExpiresType string    `json:"expires_type" enum:"DISABLE,DELETE"`
}
