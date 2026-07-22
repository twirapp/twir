package timers

import (
	"context"
	"strings"
	"testing"

	"github.com/nicklaw5/helix/v2"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
)

func TestBuildOnlineTwitchChannelRowsQueryUsesTwitchBindingWithoutFilteringDisabledBindings(t *testing.T) {
	t.Parallel()

	db := newDryRunPostgresDB(t)
	statement := buildOnlineTwitchChannelRowsQuery(
		db,
		context.Background(),
		[]string{"twitch-channel-id"},
	).Scan(&[]twitchChannelRow{}).Statement
	sql := statement.SQL.String()

	assertQueryContains(t, sql,
		`c."isEnabled" AS is_enabled`,
		"cp.platform_channel_id AS platform_id",
		"cp.user_id AS user_id",
		"FROM channel_platforms cp",
		"JOIN channels c ON c.id = cp.channel_id",
		"JOIN users u ON u.id = cp.user_id AND u.platform = 'twitch'",
		"cp.platform = 'twitch'",
		"cp.platform_channel_id IN ($1)",
	)
	assertQueryExcludes(t, sql, "c.twitch_user_id", "users.platform_id AS platform_id", "cp.enabled")

	if got := strings.Count(sql, "channel_platforms"); got != 1 {
		t.Fatalf("channel_platforms references = %d, want 1 selected Twitch binding source: %s", got, sql)
	}
}

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

	stream := onlineStream{
		stream:  streamsmodel.Stream{UserId: row.PlatformID},
		channel: row.toChannel(),
	}

	if stream.channel.User == nil {
		t.Fatal("channel owner must be attached to the Twitch stream")
	}
	if stream.channel.User.ID != row.UserID {
		t.Fatalf("owner ID = %q, want %q", stream.channel.User.ID, row.UserID)
	}
	if (&onlineUsers{}).shouldSkipStream(stream) {
		t.Fatal("an enabled Twitch channel with an unbanned owner must be processed")
	}
}
