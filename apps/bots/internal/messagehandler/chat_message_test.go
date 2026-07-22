package messagehandler

import (
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestFindChatMessageBindingUsesCanonicalBindingID(t *testing.T) {
	platformBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformTwitch,
	}
	canonicalBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformTwitch,
	}
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{platformBinding, canonicalBinding},
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
	platformBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformTwitch,
	}
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{platformBinding},
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
