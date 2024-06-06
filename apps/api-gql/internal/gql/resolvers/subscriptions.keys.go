package resolvers

const (
	notificationsSubscriptionKey = "api.newNotifications"
	chatOverlaySubscriptionKey   = "api.chatOverlaySettings"
)

func chatOverlaySubscriptionKeyCreate(userId string) string {
	return chatOverlaySubscriptionKey + "." + userId
}
