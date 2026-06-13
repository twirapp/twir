package manager

import (
	"testing"

	"github.com/kvizyx/twitchy/eventsub"
)

func TestGetConditionForTopicUsesBotIDForChatTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelChatMessage, "broadcaster-123", "bot-456")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	chatCondition, ok := condition.(eventsub.ChannelChatMessageCondition)
	if !ok {
		t.Fatalf("expected ChannelChatMessageCondition, got %T", condition)
	}

	if chatCondition.BroadcasterUserId != "broadcaster-123" {
		t.Fatalf("expected broadcaster_user_id=broadcaster-123, got %q", chatCondition.BroadcasterUserId)
	}

	if chatCondition.UserId != "bot-456" {
		t.Fatalf("expected user_id=bot-456, got %q", chatCondition.UserId)
	}
}

func TestGetConditionForTopicUsesBotIDForModeratorTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelModerate, "broadcaster-123", "bot-456")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	moderateCondition, ok := condition.(eventsub.ChannelModerateV2Condition)
	if !ok {
		t.Fatalf("expected ChannelModerateV2Condition, got %T", condition)
	}

	if moderateCondition.BroadcasterUserId != "broadcaster-123" {
		t.Fatalf("expected broadcaster_user_id=broadcaster-123, got %q", moderateCondition.BroadcasterUserId)
	}

	if moderateCondition.ModeratorUserId != "bot-456" {
		t.Fatalf("expected moderator_user_id=bot-456, got %q", moderateCondition.ModeratorUserId)
	}
}

func TestGetConditionForTopicUsesBotIDForFollowTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelFollow, "broadcaster-123", "bot-456")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	followCondition, ok := condition.(eventsub.ChannelFollowCondition)
	if !ok {
		t.Fatalf("expected ChannelFollowCondition, got %T", condition)
	}

	if followCondition.BroadcasterUserId != "broadcaster-123" {
		t.Fatalf("expected broadcaster_user_id=broadcaster-123, got %q", followCondition.BroadcasterUserId)
	}

	if followCondition.ModeratorUserId != "bot-456" {
		t.Fatalf("expected moderator_user_id=bot-456, got %q", followCondition.ModeratorUserId)
	}
}
