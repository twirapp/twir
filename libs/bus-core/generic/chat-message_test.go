package generic

import "testing"

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
		EnrichedData: ChatMessageEnrichedData{
			IsChatterVip:        true,
			IsChatterModerator:  true,
			IsChatterSubscriber: true,
		},
	}

	if !flagMsg.IsChatterVip() || !flagMsg.IsChatterModerator() || !flagMsg.IsChatterSubscriber() {
		t.Fatalf("expected enriched role flags to be honored")
	}
}
