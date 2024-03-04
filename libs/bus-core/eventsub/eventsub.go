package eventsub

const (
	EventsubSubscribeSubject = "eventsub.subscribe"
)

type EventsubSubscribeRequest struct {
	ChannelID string
}
