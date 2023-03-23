package eventsub_bindings

import (
	"encoding/json"
	"testing"
)

var jsonDataBytes = []byte(`
{
	"id": "f1c2a387-161a-49f9-a165-0f21d7a4e1c4",
	"status": "enabled",
	"type": "channel.follow",
	"version": "1",
	"cost": 1,
	"condition": {
		"broadcaster_user_id": "12826"
	},
	"transport": {
		"method": "webhook",
		"callback": "https://example.com/webhooks/callback"
	},
	"created_at": "2019-11-16T10:11:12.123Z"
}
`)

func TestSubscription_ConditionChannelFollow(t *testing.T) {
	var sub Subscription
	if err := json.Unmarshal(jsonDataBytes, &sub); err != nil {
		t.Fatalf("unmarshal data: %v", err)
		return
	}

	condition, err := sub.ConditionChannelFollow()
	if err != nil {
		t.Fatalf("condition: %v", err)
		return
	}

	if condition.BroadcasterUserID != "12826" {
		t.Fatalf("expected condition.BroadcasterUserID = \"12826\" but got %q", condition.BroadcasterUserID)
	}
}
