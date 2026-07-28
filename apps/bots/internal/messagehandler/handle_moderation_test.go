package messagehandler

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/generic"
	channelsmoderationsettingsmodel "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
)

func TestModerationOneManSpam_UsesProviderMessageIDForRedisState(t *testing.T) {
	// Given
	redisClient, recorder := newMessageIDRedisClient(t)
	handler := &MessageHandler{redis: redisClient}
	message := enrichedChatMessage{ChatMessage: generic.ChatMessage{
		ID:                uuid.NewString(),
		MessageID:         "123456789",
		BroadcasterUserId: "channel-provider-id",
		ChatterUserId:     "chatter-provider-id",
		Message:           &generic.ChatMessageMessage{Text: "repeated message"},
	}}
	settings := channelsmoderationsettingsmodel.ChannelModerationSettings{
		TriggerLength:                   1,
		OneManSpamMessageMemorySeconds:  60,
		OneManSpamMinimumStoredMessages: 1,
	}

	// When
	handler.moderationOneManSpam(context.Background(), settings, message)

	// Then
	recorder.mu.Lock()
	_, internalIDStored := recorder.members["channels:channel-provider-id:moderation:one_man_spam:chatter-provider-id\x00"+message.ID]
	_, providerIDStored := recorder.members["channels:channel-provider-id:moderation:one_man_spam:chatter-provider-id\x00"+message.MessageID]
	_, internalIDExpires := recorder.expirations["channels:channel-provider-id:moderation:one_man_spam:chatter-provider-id\x00"+message.ID]
	_, providerIDExpires := recorder.expirations["channels:channel-provider-id:moderation:one_man_spam:chatter-provider-id\x00"+message.MessageID]
	recorder.mu.Unlock()
	if internalIDStored || !providerIDStored || internalIDExpires || !providerIDExpires {
		t.Fatalf(
			"one-man spam Redis state stored internal=%t provider=%t expires internal=%t provider=%t, want false true false true",
			internalIDStored,
			providerIDStored,
			internalIDExpires,
			providerIDExpires,
		)
	}
}
