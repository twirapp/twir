package kick

const KickStreamOfflineSubject = "kick.stream-offline"

type KickStreamOffline struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
}
