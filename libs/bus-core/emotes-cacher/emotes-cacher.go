package emotes_cacher

const (
	EmotesCacherGetChannelEmotesSubject = "emotes-cacher.get-channel-emotes"
	EmotesCacherGetGlobalEmotesSubject  = "emotes-cacher.get-global-emotes"
)

type ServiceName string

const (
	ServiceNameSevenTV ServiceName = "7tv"
	ServiceNameBTTV    ServiceName = "bttv"
	ServiceNameFFZ     ServiceName = "ffz"
)

type Emote struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Service ServiceName `json:"service"`
}

type GetChannelEmotesRequest struct {
	ChannelID string        `json:"channel_id"`
	ServiceIn []ServiceName `json:"service_in"`
}

type GetGlobalEmotesRequest struct {
	ServiceIn []ServiceName `json:"service_in"`
}

type Response struct {
	Emotes []Emote `json:"emotes"`
}
