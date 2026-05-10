package manager

import (
	"testing"

	"github.com/kvizyx/twitchy/eventsub"
)

func TestGetConditionForTopicUsesBroadcasterIDForChatTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelChatMessage, "123")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	chatCondition, ok := condition.(eventsub.ChannelChatMessageCondition)
	if !ok {
		t.Fatalf("expected ChannelChatMessageCondition, got %T", condition)
	}

	if chatCondition.BroadcasterUserId != "123" {
		t.Fatalf("expected broadcaster_user_id=123, got %q", chatCondition.BroadcasterUserId)
	}

	if chatCondition.UserId != "123" {
		t.Fatalf("expected user_id=123, got %q", chatCondition.UserId)
	}
}

func TestGetConditionForTopicUsesBroadcasterIDForModeratorTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelModerate, "456")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	moderateCondition, ok := condition.(eventsub.ChannelModerateV2Condition)
	if !ok {
		t.Fatalf("expected ChannelModerateV2Condition, got %T", condition)
	}

	if moderateCondition.BroadcasterUserId != "456" {
		t.Fatalf("expected broadcaster_user_id=456, got %q", moderateCondition.BroadcasterUserId)
	}

	if moderateCondition.ModeratorUserId != "456" {
		t.Fatalf("expected moderator_user_id=456, got %q", moderateCondition.ModeratorUserId)
	}
}

func TestGetConditionForTopicUsesBroadcasterIDForFollowTopics(t *testing.T) {
	m := &Manager{}

	condition, err := m.getConditionForTopic(eventsub.EventTypeChannelFollow, "789")
	if err != nil {
		t.Fatalf("getConditionForTopic returned error: %v", err)
	}

	followCondition, ok := condition.(eventsub.ChannelFollowCondition)
	if !ok {
		t.Fatalf("expected ChannelFollowCondition, got %T", condition)
	}

	if followCondition.BroadcasterUserId != "789" {
		t.Fatalf("expected broadcaster_user_id=789, got %q", followCondition.BroadcasterUserId)
	}

	if followCondition.ModeratorUserId != "789" {
		t.Fatalf("expected moderator_user_id=789, got %q", followCondition.ModeratorUserId)
	}
}
