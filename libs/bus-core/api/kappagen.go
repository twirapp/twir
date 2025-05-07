package api

const (
	TriggerKappagenSubject = "api.kappagen.trigger"
)

type TriggerKappagenMessage struct {
	ChannelId string
	Text      string
	Emotes    []TriggerKappagenEmote
}

type TriggerKappagenEmote struct {
	Id        string
	Positions []string
}
