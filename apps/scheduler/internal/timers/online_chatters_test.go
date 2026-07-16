package timers

import (
	"testing"

	model "github.com/twirapp/twir/libs/gomodels"
)

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
