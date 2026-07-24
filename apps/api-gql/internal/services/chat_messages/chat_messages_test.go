package chat_messages

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/bus-core/events"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/wsrouter"
)

func TestHandleChannelBanEventUsesSelectedEventPlatformBinding(t *testing.T) {
	channelID := uuid.New()
	platformChannelID := "kick-channel"
	key := chatOverlayModerationSubscriptionKeyCreate(
		platformentity.PlatformKick.String(),
		platformChannelID,
	)
	router := &chatMessagesTestWsRouter{}
	lookup := &chatMessagesTestChannelLookup{channels: map[uuid.UUID]channelentity.Channel{
		channelID: {
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{
				{
					Platform:          platformentity.PlatformVKVideoLive,
					PlatformChannelID: "vk-channel",
					BotConfig:         json.RawMessage(`{"unexpected":"vk"}`),
				},
				{
					Platform:          platformentity.PlatformTwitch,
					PlatformChannelID: "twitch-channel",
					BotConfig:         json.RawMessage(`{"unexpected":"twitch"}`),
				},
				{
					Platform:          platformentity.PlatformKick,
					PlatformChannelID: platformChannelID,
					BotConfig:         json.RawMessage(`{"unexpected":"kick"}`),
				},
			},
		},
	}}
	service := &Service{
		channelService: lookup,
		wsRouter:       router,
		chanSubs:       map[string]struct{}{key: {}},
	}

	_, err := service.handleChannelBanEvent(context.Background(), events.ChannelBanMessage{
		BaseInfo: events.BaseInfo{
			Platform:          platformentity.PlatformKick,
			ChannelPlatformID: channelID.String(),
		},
		UserLogin: "banned-user",
	})
	if err != nil {
		t.Fatalf("handle channel ban event: %v", err)
	}

	if len(router.published) != 1 {
		t.Fatalf("published events = %d, want 1", len(router.published))
	}
	if router.published[0].key != key {
		t.Fatalf("published key = %q, want %q", router.published[0].key, key)
	}

	event, ok := router.published[0].data.(entity.ChatOverlayModerationEvent)
	if !ok {
		t.Fatalf("published data = %T, want ChatOverlayModerationEvent", router.published[0].data)
	}
	if event.Platform != platformentity.PlatformKick.String() || event.UserLogin != "banned-user" {
		t.Fatalf("published event = %#v, want Kick ban for banned-user", event)
	}
}

func TestHandleChannelBanEventSkipsMissingEventPlatformBinding(t *testing.T) {
	channelID := uuid.New()
	wrongProviderID := uuid.New()
	wrongProviderChannelID := wrongProviderID.String()
	wrongProviderRoute := chatOverlayModerationSubscriptionKeyCreate(
		platformentity.PlatformKick.String(),
		wrongProviderChannelID,
	)
	wrongProviderLookupRoute := chatOverlayModerationSubscriptionKeyCreate(
		platformentity.PlatformKick.String(),
		"wrong-kick-channel",
	)
	router := &chatMessagesTestWsRouter{}
	lookup := &chatMessagesTestChannelLookup{channels: map[uuid.UUID]channelentity.Channel{
		channelID: {
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platformentity.PlatformVKVideoLive, PlatformChannelID: wrongProviderChannelID},
			},
		},
		wrongProviderID: {
			ID: wrongProviderID,
			Bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platformentity.PlatformKick, PlatformChannelID: "wrong-kick-channel"},
			},
		},
	}}
	service := &Service{
		channelService: lookup,
		wsRouter:       router,
		chanSubs: map[string]struct{}{
			wrongProviderRoute:       {},
			wrongProviderLookupRoute: {},
		},
	}

	_, err := service.handleChannelBanEvent(context.Background(), events.ChannelBanMessage{
		BaseInfo: events.BaseInfo{
			Platform:          platformentity.PlatformKick,
			ChannelPlatformID: channelID.String(),
		},
	})
	if err != nil {
		t.Fatalf("handle channel ban event: %v", err)
	}
	if len(lookup.channelIDs) != 1 || lookup.channelIDs[0] != channelID {
		t.Fatalf("channel lookups = %#v, want only parsed internal channel ID %s", lookup.channelIDs, channelID)
	}
	if len(router.published) != 0 {
		t.Fatalf("published events = %d, want 0", len(router.published))
	}
}

type chatMessagesTestChannelLookup struct {
	channels   map[uuid.UUID]channelentity.Channel
	channelIDs []uuid.UUID
}

func (r *chatMessagesTestChannelLookup) GetChannelByID(
	_ context.Context,
	channelID uuid.UUID,
) (channelentity.Channel, error) {
	r.channelIDs = append(r.channelIDs, channelID)
	channel, ok := r.channels[channelID]
	if !ok {
		return channelentity.Nil, nil
	}

	return channel, nil
}

type chatMessagesPublishedEvent struct {
	key  string
	data any
}

type chatMessagesTestWsRouter struct {
	published []chatMessagesPublishedEvent
}

func (*chatMessagesTestWsRouter) Subscribe([]string) (wsrouter.WsRouterSubscription, error) {
	return nil, nil
}

func (r *chatMessagesTestWsRouter) Publish(key string, data any) error {
	r.published = append(r.published, chatMessagesPublishedEvent{key: key, data: data})
	return nil
}
