package resolvers

const (
	notificationsSubscriptionKey          = "api.newNotifications"
	chatOverlaySubscriptionKey            = "api.chatOverlaySettings"
	nowPlayingOverlaySubscriptionKey      = "api.nowPlayingOverlaySettings"
	dashboardWidgetsLayoutSubscriptionKey = "api.dashboardWidgetsLayout"
)

func chatOverlaySubscriptionKeyCreate(id, userId string) string {
	return chatOverlaySubscriptionKey + "." + userId + "." + id
}

func nowPlayingOverlaySubscriptionKeyCreate(id, userId string) string {
	return nowPlayingOverlaySubscriptionKey + "." + userId + "." + id
}

func dashboardWidgetsLayoutSubscriptionKeyCreate(channelID string) string {
	return dashboardWidgetsLayoutSubscriptionKey + "." + channelID
}
