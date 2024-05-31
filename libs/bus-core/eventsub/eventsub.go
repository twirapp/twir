package eventsub

const (
	EventsubSubscribeAllSubject = "eventsub.subscribeAll"
	EventsubSubscribeSubject    = "eventsub.subscribe"
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
