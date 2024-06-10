package resolvers

const (
	notificationsSubscriptionKey = "api.newNotifications"
	chatOverlaySubscriptionKey   = "api.chatOverlaySettings"
)

func chatOverlaySubscriptionKeyCreate(id, userId string) string {
	return chatOverlaySubscriptionKey + "." + userId + "." + id
}
