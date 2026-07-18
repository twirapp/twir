package timers

import (
	"testing"

	"github.com/nicklaw5/helix/v2"
	model "github.com/twirapp/twir/libs/gomodels"
)

func TestAppendUniqueChattersRetainsEveryPage(t *testing.T) {
	t.Parallel()

	chatters := make([]helix.ChatChatter, 0)
	indices := make(map[string]int)
	chatters = appendUniqueChatters(chatters, indices, []helix.ChatChatter{
		{UserID: "one"},
		{UserID: "two"},
	})
	chatters = appendUniqueChatters(chatters, indices, []helix.ChatChatter{
		{UserID: "three"},
		{UserID: "four"},
	})

	if len(chatters) != 4 {
		t.Fatalf("chatters count = %d, want 4", len(chatters))
	}
	for index, userID := range []string{"one", "two", "three", "four"} {
		if chatters[index].UserID != userID {
			t.Fatalf("chatter at index %d = %q, want %q", index, chatters[index].UserID, userID)
		}
	}
}

func TestOrderChattersForUserInsertKeepsSnapshotIntact(t *testing.T) {
	t.Parallel()

	snapshot := []helix.ChatChatter{
		{UserID: "3"},
		{UserID: "1"},
		{UserID: "2"},
	}

	ordered := orderChattersForUserInsert(snapshot)

	for index, userID := range []string{"1", "2", "3"} {
		if ordered[index].UserID != userID {
			t.Fatalf("ordered chatter at index %d = %q, want %q", index, ordered[index].UserID, userID)
		}
	}
	if snapshot[0].UserID != "3" {
		t.Fatalf("snapshot was changed to start with %q", snapshot[0].UserID)
	}
}

func TestTwitchChannelRowToChannelKeepsOwner(t *testing.T) {
	t.Parallel()

	row := twitchChannelRow{
		ChannelID:  "channel-id",
		PlatformID: "twitch-platform-id",
		UserID:     "owner-id",
		IsEnabled:  true,
		IsBanned:   false,
	}

	stream := &model.ChannelsStreams{UserId: row.PlatformID}
	stream.Channel = row.toChannel()

	if stream.Channel.User == nil {
		t.Fatal("channel owner must be attached to the Twitch stream")
	}
	if stream.Channel.User.ID != row.UserID {
		t.Fatalf("owner ID = %q, want %q", stream.Channel.User.ID, row.UserID)
	}
	if (&onlineUsers{}).shouldSkipStream(stream) {
		t.Fatal("an enabled Twitch channel with an unbanned owner must be processed")
	}
}
