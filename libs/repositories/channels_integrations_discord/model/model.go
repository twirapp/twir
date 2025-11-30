package model

type ChannelIntegrationDiscord struct {
	ID                               int
	ChannelID                        string
	GuildID                          string
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

	isNil bool
}

func (c ChannelIntegrationDiscord) IsNil() bool {
	return c.isNil
}

var Nil = ChannelIntegrationDiscord{
	isNil: true,
}
