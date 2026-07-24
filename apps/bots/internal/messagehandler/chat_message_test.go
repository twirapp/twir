package messagehandler

import (
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/generic"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestFindChatMessageBindingUsesCanonicalBindingID(t *testing.T) {
	platformBinding := channelplatformentity.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformTwitch,
	}
	canonicalBinding := channelplatformentity.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformTwitch,
	}
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{platformBinding, canonicalBinding},
	}

	binding, found, err := findChatMessageBinding(
		channel,
		generic.ChatMessage{ChannelBindingID: canonicalBinding.ID.String()},
		platform.PlatformTwitch,
	)
	if err != nil {
		t.Fatalf("find chat message binding: %v", err)
	}
	if !found {
		t.Fatal("expected canonical binding")
	}
	if binding.ID != canonicalBinding.ID {
		t.Fatalf("binding ID = %s, want %s", binding.ID, canonicalBinding.ID)
	}
}

func TestFindChatMessageBindingFallsBackOnlyWhenBindingIDIsAbsent(t *testing.T) {
	platformBinding := channelplatformentity.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformTwitch,
	}
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{platformBinding},
	}

	binding, found, err := findChatMessageBinding(
		channel,
		generic.ChatMessage{},
		platform.PlatformTwitch,
	)
	if err != nil {
		t.Fatalf("find chat message binding without ID: %v", err)
	}
	if !found || binding.ID != platformBinding.ID {
		t.Fatalf("binding = %#v, found = %t, want platform binding %s", binding, found, platformBinding.ID)
	}

	_, found, err = findChatMessageBinding(
		channel,
		generic.ChatMessage{ChannelBindingID: uuid.New().String()},
		platform.PlatformTwitch,
	)
	if err != nil {
		t.Fatalf("find missing canonical binding: %v", err)
	}
	if found {
		t.Fatal("expected no binding when the supplied binding ID is unknown")
	}

	_, found, err = findChatMessageBinding(
		channel,
		generic.ChatMessage{ChannelBindingID: "not-a-uuid"},
		platform.PlatformTwitch,
	)
	if err == nil {
		t.Fatal("expected an invalid supplied binding ID to fail")
	}
	if found {
		t.Fatal("expected no binding for an invalid supplied binding ID")
	}
}
