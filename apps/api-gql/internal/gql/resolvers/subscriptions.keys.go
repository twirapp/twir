package resolvers

const (
	notificationsSubscriptionKey     = "api.newNotifications"
	chatOverlaySubscriptionKey       = "api.chatOverlaySettings"
	nowPlayingOverlaySubscriptionKey = "api.nowPlayingOverlaySettings"
)

func chatOverlaySubscriptionKeyCreate(id, userId string) string {
	return chatOverlaySubscriptionKey + "." + userId + "." + id
}

func nowPlayingOverlaySubscriptionKeyCreate(id, userId string) string {
	return nowPlayingOverlaySubscriptionKey + "." + userId + "." + id
}
