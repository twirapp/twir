package discord

// ChannelType represents the type of a Discord channel
type ChannelType int

const (
	ChannelTypeVoice ChannelType = iota
	ChannelTypeText
)

// Subject names for Discord bus
const (
	GetGuildChannelsSubject = "discord.get_guild_channels"
	GetGuildInfoSubject     = "discord.get_guild_info"
	LeaveGuildSubject       = "discord.leave_guild"
	GetGuildRolesSubject    = "discord.get_guild_roles"
)

// GetGuildChannelsRequest is the request for getting guild channels
type GetGuildChannelsRequest struct {
	GuildID string
}

// GuildChannel represents a Discord guild channel
type GuildChannel struct {
	ID              string
	Name            string
	Type            ChannelType
	CanSendMessages bool
}

// GetGuildChannelsResponse is the response for getting guild channels
type GetGuildChannelsResponse struct {
	Channels []GuildChannel
}

// GetGuildInfoRequest is the request for getting guild info
type GetGuildInfoRequest struct {
	GuildID string
}

// Role represents a Discord role
type Role struct {
	ID    string
	Name  string
	Color string
}

// GetGuildInfoResponse is the response for getting guild info
type GetGuildInfoResponse struct {
	ID       string
	Name     string
	Icon     string
	Channels []GuildChannel
	Roles    []Role
}

// LeaveGuildRequest is the request for leaving a guild
type LeaveGuildRequest struct {
	GuildID string
}

// GetGuildRolesRequest is the request for getting guild roles
type GetGuildRolesRequest struct {
	GuildID string
}

// GetGuildRolesResponse is the response for getting guild roles
type GetGuildRolesResponse struct {
	Roles []Role
}
