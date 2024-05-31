package eventsub

const (
	EventsubSubscribeSubject = "eventsub.subscribe"
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
