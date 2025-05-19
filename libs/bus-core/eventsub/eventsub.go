package eventsub

const (
	EventsubSubscribeAllSubject = "eventsub.subscribeAll"
	EventsubSubscribeSubject    = "eventsub.subscribe"
	EventsubInitChannelsSubject = "eventsub.initChannels"
	EventsubUnsubscribeSubject  = "eventsub.unsubscribe"
)

type EventsubSubscribeToAllEventsRequest struct {
	ChannelID string
}

type EventsubSubscribeRequest struct {
	ChannelID     string
	Topic         string
	ConditionType string
	Version       string
}
