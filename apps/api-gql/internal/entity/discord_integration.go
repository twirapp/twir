package entity

type DiscordGuild struct {
	ID                               string
	Name                             string
	Icon                             *string
	LiveNotificationEnabled          bool
	LiveNotificationChannelsIds      []string
	LiveNotificationShowTitle        bool
	LiveNotificationShowCategory     bool
	LiveNotificationShowViewers      bool
	LiveNotificationMessage          string
	LiveNotificationShowPreview      bool
	LiveNotificationShowProfileImage bool
	OfflineNotificationMessage       string
	ShouldDeleteMessageOnOffline     bool
	AdditionalUsersIdsForLiveCheck   []string
}

type DiscordIntegrationData struct {
	Guilds []DiscordGuild
}

type DiscordGuildChannel struct {
	ID              string
	Name            string
	Type            DiscordChannelType
	CanSendMessages bool
}

type DiscordChannelType int

const (
DiscordChannelTypeText DiscordChannelType = iota
DiscordChannelTypeVoice
DiscordChannelTypeDM
DiscordChannelTypeGroupDM
DiscordChannelTypeCategory
DiscordChannelTypeAnnouncement
DiscordChannelTypeAnnouncementThread
DiscordChannelTypePublicThread
DiscordChannelTypePrivateThread
DiscordChannelTypeStageVoice
DiscordChannelTypeDirectory
DiscordChannelTypeForum
DiscordChannelTypeMedia
)

type DiscordGuildRole struct {
	ID    string
	Name  string
	Color string
}

type DiscordGuildInfo struct {
	ID       string
	Name     string
	Icon     *string
	Channels []DiscordGuildChannel
	Roles    []DiscordGuildRole
}
