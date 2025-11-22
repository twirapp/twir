package api

const (
	TriggerTtsSaySubject  = "api.tts.say.trigger"
	TriggerTtsSkipSubject = "api.tts.skip.trigger"
)

type TriggerTtsSay struct {
	ChannelId string
	Text      string
	Voice     string
	Rate      string
	Pitch     string
	Volume    string
}

type TriggerTtsSkip struct {
	ChannelId string
}

