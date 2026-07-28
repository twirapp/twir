package generic

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestChatMessageJSONIsCanonicalTransport(t *testing.T) {
	var message ChatMessage
	if err := json.Unmarshal(
		[]byte(`{
			"id":"event-123",
			"message_id":"message-123",
			"platform":"kick",
			"channel_id":"channel-123",
			"channel_binding_id":"binding-123",
			"platform_channel_id":"provider-channel-123",
			"user_id":"user-123",
			"message":{"text":"hello Kappa","fragments":[{"type":2,"text":"Kappa","emote":{"id":"25"}}]},
			"badges":[{"id":"moderator","set_id":"moderator","text":"Moderator"}],
			"is_broadcaster":true,
			"is_moderator":true,
			"is_vip":true,
			"is_subscriber":true
		}`),
		&message,
	); err != nil {
		t.Fatalf("unmarshal canonical chat message: %v", err)
	}

	encoded, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("marshal canonical chat message: %v", err)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatalf("unmarshal marshaled chat message: %v", err)
	}

	for _, field := range []string{
		"enriched_data",
		"channel",
		"stream",
		"user",
		"user_stats",
	} {
		if _, ok := payload[field]; ok {
			t.Fatalf("canonical payload must not contain repository enrichment %q: %s", field, encoded)
		}
	}

	assertChatMessageJSONField(t, payload, "id", `"event-123"`)
	assertChatMessageJSONField(t, payload, "message_id", `"message-123"`)
	assertChatMessageJSONField(t, payload, "platform", `"kick"`)
	assertChatMessageJSONField(t, payload, "channel_id", `"channel-123"`)
	assertChatMessageJSONField(t, payload, "channel_binding_id", `"binding-123"`)
	assertChatMessageJSONField(t, payload, "platform_channel_id", `"provider-channel-123"`)
	assertChatMessageJSONField(t, payload, "user_id", `"user-123"`)
	assertChatMessageJSONField(t, payload, "is_broadcaster", "true")
	assertChatMessageJSONField(t, payload, "is_moderator", "true")
	assertChatMessageJSONField(t, payload, "is_vip", "true")
	assertChatMessageJSONField(t, payload, "is_subscriber", "true")

	var content struct {
		Text      string `json:"text"`
		Fragments []struct {
			Text  string `json:"text"`
			Emote *struct {
				ID string `json:"id"`
			} `json:"emote"`
		} `json:"fragments"`
	}
	if err := json.Unmarshal(payload["message"], &content); err != nil {
		t.Fatalf("unmarshal message content: %v", err)
	}
	if content.Text != "hello Kappa" || len(content.Fragments) != 1 ||
		content.Fragments[0].Text != "Kappa" || content.Fragments[0].Emote == nil ||
		content.Fragments[0].Emote.ID != "25" {
		t.Fatalf("message content/fragments were not preserved: %s", payload["message"])
	}

	var badges []ChatMessageBadge
	if err := json.Unmarshal(payload["badges"], &badges); err != nil {
		t.Fatalf("unmarshal badges: %v", err)
	}
	if len(badges) != 1 || badges[0].SetID != "moderator" || badges[0].Text != "Moderator" {
		t.Fatalf("badges were not preserved: %s", payload["badges"])
	}
}

func assertChatMessageJSONField(
	t *testing.T,
	payload map[string]json.RawMessage,
	field string,
	want string,
) {
	t.Helper()

	if got := string(payload[field]); got != want {
		t.Fatalf("payload field %q = %s, want %s", field, got, want)
	}
}

func TestChatMessageRoleHelpers(t *testing.T) {
	msg := ChatMessage{
		BroadcasterUserId: "100",
		ChatterUserId:     "200",
		Badges: []ChatMessageBadge{
			{SetID: "moderator", Text: "Moderator"},
			{SetID: "subscriber", Text: "Subscriber"},
		},
	}

	if !msg.IsChatterModerator() {
		t.Fatalf("expected moderator badge to be detected")
	}

	if !msg.IsChatterSubscriber() {
		t.Fatalf("expected subscriber badge to be detected")
	}

	if msg.IsChatterBroadcaster() {
		t.Fatalf("did not expect broadcaster without matching ids or flag")
	}

	kickMsg := ChatMessage{
		BroadcasterUserId: "300",
		ChatterUserId:     "300",
		Badges:            []ChatMessageBadge{{ID: "moderator", Text: "Moderator"}},
	}

	if !kickMsg.IsChatterBroadcaster() {
		t.Fatalf("expected broadcaster when chatter and broadcaster ids match")
	}

	if !kickMsg.IsChatterModerator() {
		t.Fatalf("expected moderator from normalized badge id")
	}

	flagMsg := ChatMessage{
		IsVip:        true,
		IsModerator:  true,
		IsSubscriber: true,
	}

	if !flagMsg.IsChatterVip() || !flagMsg.IsChatterModerator() || !flagMsg.IsChatterSubscriber() {
		t.Fatalf("expected enriched role flags to be honored")
	}
}
