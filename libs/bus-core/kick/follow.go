package kick

const KickChannelFollowSubject = "kick.channel-follow"

type KickChannelFollow struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
	FollowerUserID    string `json:"follower_user_id"`
	FollowerUserLogin string `json:"follower_user_login"`
}
