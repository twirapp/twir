package api

const (
	TriggerBrbStartSubject = "api.brb.start.trigger"
	TriggerBrbStopSubject  = "api.brb.stop.trigger"
)

type TriggerBrbStart struct {
	ChannelId string
	Minutes   int32
	Text      *string
}

type TriggerBrbStop struct {
	ChannelId string
}
