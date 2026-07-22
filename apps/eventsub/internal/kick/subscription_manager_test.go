package kick

import (
	"testing"

	"github.com/google/uuid"
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
)

func TestRedisKeyIncludesBindingEventAndTransportKind(t *testing.T) {
	bindingID := uuid.MustParse("9847545a-c603-4a30-88ca-7c2e621d2069")

	got := redisKey(bindingID, "channel.followed", buscoreeventsub.TransportWebhook)
	want := "kick:sub:9847545a-c603-4a30-88ca-7c2e621d2069:channel.followed:webhook"
	if got != want {
		t.Errorf("redisKey() = %q, want %q", got, want)
	}
}
