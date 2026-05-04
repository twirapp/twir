package kick

const KickStreamOnlineSubject = "kick.stream-online"

type KickStreamOnline struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
}
