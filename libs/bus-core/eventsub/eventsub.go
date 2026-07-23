package eventsub

import platformentity "github.com/twirapp/twir/libs/entities/platform"

const (
	EventsubSubscribeAllSubject = "eventsub.subscribeAll"
	EventsubSubscribeSubject    = "eventsub.subscribe"
	EventsubInitChannelsSubject = "eventsub.initChannels"
	EventsubUnsubscribeSubject  = "eventsub.unsubscribe"
)

type TransportKind string

const (
	TransportWebhook   TransportKind = "webhook"
	TransportWebSocket TransportKind = "websocket"
)

type EventsubSubscribeToAllEventsRequest struct {
	ChannelID string
	Platform  platformentity.Platform
}

type EventsubSubscribeRequest struct {
	ChannelID string
	Topic     string
	Version   string
}

type EventsubBindingSnapshot struct {
	ID                string
	UserID            string
	PlatformChannelID string
}

type EventsubUnsubscribeRequest struct {
	ChannelID string
	Platform  platformentity.Platform
	Binding   *EventsubBindingSnapshot
}
