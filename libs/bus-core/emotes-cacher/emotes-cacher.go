package emotes_cacher

const (
	EMOTES_CACHER_GLOBAL_EMOTES_SUBJECT  = "emotes-cacher.global-emotes"
	EMOTES_CACHER_CHANNEL_EMOTES_SUBJECT = "emotes-cacher.channel-emotes"
)

type EmotesCacheRequest struct {
	ChannelID string
}
