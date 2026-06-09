package eventsub

import platformentity "github.com/twirapp/twir/libs/entities/platform"

const (
	EventsubSubscribeAllSubject = "eventsub.subscribeAll"
	EventsubSubscribeSubject    = "eventsub.subscribe"
	EventsubInitChannelsSubject = "eventsub.initChannels"
	EventsubUnsubscribeSubject  = "eventsub.unsubscribe"
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

type EventsubUnsubscribeRequest struct {
	ChannelID string
	Platform  platformentity.Platform
}
