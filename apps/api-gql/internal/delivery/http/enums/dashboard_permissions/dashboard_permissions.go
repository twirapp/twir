package dashboard_permissions

type ChannelRolePermissionEnum string

const (
	ChannelRolePermissionEnumCanAccessDashboard    ChannelRolePermissionEnum = "CAN_ACCESS_DASHBOARD"
	ChannelRolePermissionEnumUpdateChannelTitle    ChannelRolePermissionEnum = "UPDATE_CHANNEL_TITLE"
	ChannelRolePermissionEnumUpdateChannelCategory ChannelRolePermissionEnum = "UPDATE_CHANNEL_CATEGORY"
	ChannelRolePermissionEnumViewCommands          ChannelRolePermissionEnum = "VIEW_COMMANDS"
	ChannelRolePermissionEnumManageCommands        ChannelRolePermissionEnum = "MANAGE_COMMANDS"
	ChannelRolePermissionEnumViewKeywords          ChannelRolePermissionEnum = "VIEW_KEYWORDS"
	ChannelRolePermissionEnumManageKeywords        ChannelRolePermissionEnum = "MANAGE_KEYWORDS"
	ChannelRolePermissionEnumViewTimers            ChannelRolePermissionEnum = "VIEW_TIMERS"
	ChannelRolePermissionEnumManageTimers          ChannelRolePermissionEnum = "MANAGE_TIMERS"
	ChannelRolePermissionEnumViewIntegrations      ChannelRolePermissionEnum = "VIEW_INTEGRATIONS"
	ChannelRolePermissionEnumManageIntegrations    ChannelRolePermissionEnum = "MANAGE_INTEGRATIONS"
	ChannelRolePermissionEnumViewSongRequests      ChannelRolePermissionEnum = "VIEW_SONG_REQUESTS"
	ChannelRolePermissionEnumManageSongRequests    ChannelRolePermissionEnum = "MANAGE_SONG_REQUESTS"
	ChannelRolePermissionEnumViewModeration        ChannelRolePermissionEnum = "VIEW_MODERATION"
	ChannelRolePermissionEnumManageModeration      ChannelRolePermissionEnum = "MANAGE_MODERATION"
	ChannelRolePermissionEnumViewVariables         ChannelRolePermissionEnum = "VIEW_VARIABLES"
	ChannelRolePermissionEnumManageVariables       ChannelRolePermissionEnum = "MANAGE_VARIABLES"
	ChannelRolePermissionEnumViewGreetings         ChannelRolePermissionEnum = "VIEW_GREETINGS"
	ChannelRolePermissionEnumManageGreetings       ChannelRolePermissionEnum = "MANAGE_GREETINGS"
	ChannelRolePermissionEnumViewOverlays          ChannelRolePermissionEnum = "VIEW_OVERLAYS"
	ChannelRolePermissionEnumManageOverlays        ChannelRolePermissionEnum = "MANAGE_OVERLAYS"
	ChannelRolePermissionEnumViewRoles             ChannelRolePermissionEnum = "VIEW_ROLES"
	ChannelRolePermissionEnumManageRoles           ChannelRolePermissionEnum = "MANAGE_ROLES"
	ChannelRolePermissionEnumViewEvents            ChannelRolePermissionEnum = "VIEW_EVENTS"
	ChannelRolePermissionEnumManageEvents          ChannelRolePermissionEnum = "MANAGE_EVENTS"
	ChannelRolePermissionEnumViewAlerts            ChannelRolePermissionEnum = "VIEW_ALERTS"
	ChannelRolePermissionEnumManageAlerts          ChannelRolePermissionEnum = "MANAGE_ALERTS"
	ChannelRolePermissionEnumViewGames             ChannelRolePermissionEnum = "VIEW_GAMES"
	ChannelRolePermissionEnumManageGames           ChannelRolePermissionEnum = "MANAGE_GAMES"
	ChannelRolePermissionEnumViewBotSettings       ChannelRolePermissionEnum = "VIEW_BOT_SETTINGS"
	ChannelRolePermissionEnumManageBotSettings     ChannelRolePermissionEnum = "MANAGE_BOT_SETTINGS"
)

var AllChannelRolePermissionEnum = []ChannelRolePermissionEnum{
	ChannelRolePermissionEnumCanAccessDashboard,
	ChannelRolePermissionEnumUpdateChannelTitle,
	ChannelRolePermissionEnumUpdateChannelCategory,
	ChannelRolePermissionEnumViewCommands,
	ChannelRolePermissionEnumManageCommands,
	ChannelRolePermissionEnumViewKeywords,
	ChannelRolePermissionEnumManageKeywords,
	ChannelRolePermissionEnumViewTimers,
	ChannelRolePermissionEnumManageTimers,
	ChannelRolePermissionEnumViewIntegrations,
	ChannelRolePermissionEnumManageIntegrations,
	ChannelRolePermissionEnumViewSongRequests,
	ChannelRolePermissionEnumManageSongRequests,
	ChannelRolePermissionEnumViewModeration,
	ChannelRolePermissionEnumManageModeration,
	ChannelRolePermissionEnumViewVariables,
	ChannelRolePermissionEnumManageVariables,
	ChannelRolePermissionEnumViewGreetings,
	ChannelRolePermissionEnumManageGreetings,
	ChannelRolePermissionEnumViewOverlays,
	ChannelRolePermissionEnumManageOverlays,
	ChannelRolePermissionEnumViewRoles,
	ChannelRolePermissionEnumManageRoles,
	ChannelRolePermissionEnumViewEvents,
	ChannelRolePermissionEnumManageEvents,
	ChannelRolePermissionEnumViewAlerts,
	ChannelRolePermissionEnumManageAlerts,
	ChannelRolePermissionEnumViewGames,
	ChannelRolePermissionEnumManageGames,
	ChannelRolePermissionEnumViewBotSettings,
	ChannelRolePermissionEnumManageBotSettings,
}

func (e ChannelRolePermissionEnum) IsValid() bool {
	switch e {
	case ChannelRolePermissionEnumCanAccessDashboard, ChannelRolePermissionEnumUpdateChannelTitle, ChannelRolePermissionEnumUpdateChannelCategory, ChannelRolePermissionEnumViewCommands, ChannelRolePermissionEnumManageCommands, ChannelRolePermissionEnumViewKeywords, ChannelRolePermissionEnumManageKeywords, ChannelRolePermissionEnumViewTimers, ChannelRolePermissionEnumManageTimers, ChannelRolePermissionEnumViewIntegrations, ChannelRolePermissionEnumManageIntegrations, ChannelRolePermissionEnumViewSongRequests, ChannelRolePermissionEnumManageSongRequests, ChannelRolePermissionEnumViewModeration, ChannelRolePermissionEnumManageModeration, ChannelRolePermissionEnumViewVariables, ChannelRolePermissionEnumManageVariables, ChannelRolePermissionEnumViewGreetings, ChannelRolePermissionEnumManageGreetings, ChannelRolePermissionEnumViewOverlays, ChannelRolePermissionEnumManageOverlays, ChannelRolePermissionEnumViewRoles, ChannelRolePermissionEnumManageRoles, ChannelRolePermissionEnumViewEvents, ChannelRolePermissionEnumManageEvents, ChannelRolePermissionEnumViewAlerts, ChannelRolePermissionEnumManageAlerts, ChannelRolePermissionEnumViewGames, ChannelRolePermissionEnumManageGames, ChannelRolePermissionEnumViewBotSettings, ChannelRolePermissionEnumManageBotSettings:
		return true
	}
	return false
}

func (e ChannelRolePermissionEnum) String() string {
	return string(e)
}
