package resolvers

import (
	"reflect"
	"testing"

	channelsservice "github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	twir_events "github.com/twirapp/twir/apps/api-gql/internal/services/twir-events"
	chatmessagesrepo "github.com/twirapp/twir/libs/repositories/chat_messages"
)

func TestBuildTwirEventSubscriptionKeysSkipsKickPlatformID(t *testing.T) {
	got := buildTwirEventSubscriptionKeys(channelsservice.ApiKeyChannelIdentity{
		InternalChannelID: "internal-channel-id",
		ChatTargets: []chatmessagesrepo.PlatformChannelIdentity{
			{Platform: "twitch", PlatformChannelID: "123"},
			{Platform: "kick", PlatformChannelID: "123"},
		},
	})

	want := []string{
		twir_events.CreateSubscribeKey("internal-channel-id"),
		twir_events.CreateSubscribeKey("123"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildTwirEventSubscriptionKeys() = %#v, want %#v", got, want)
	}
}
