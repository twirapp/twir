package platforms

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

type testChatAdapter struct {
	platform     platform.Platform
	capabilities platform.Capabilities
}

func (a testChatAdapter) Platform() platform.Platform {
	return a.platform
}

func (a testChatAdapter) Capabilities() platform.Capabilities {
	return a.capabilities
}

func (testChatAdapter) SendMessage(
	context.Context,
	channelplatformsmodel.ChannelPlatform,
	string,
	string,
	ChatOptions,
) error {
	return nil
}

func TestNewRegistryRegistersChatAdapters(t *testing.T) {
	t.Parallel()

	registry := newRegistry(
		testChatAdapter{
			platform:     platform.PlatformTwitch,
			capabilities: platform.Capabilities{platform.CapabilityChatWrite},
		},
		testChatAdapter{
			platform:     platform.PlatformKick,
			capabilities: platform.Capabilities{platform.CapabilityChatWrite},
		},
	)

	for _, current := range []platform.Platform{platform.PlatformTwitch, platform.PlatformKick} {
		if _, err := registry.Require(current, platform.CapabilityChatWrite); err != nil {
			t.Errorf("expected %s chat adapter: %v", current, err)
		}
	}
}

type fakeTwitchSender struct {
	binding channelplatformsmodel.ChannelPlatform
	opts    twitchactions.SendMessageOpts
}

func (s *fakeTwitchSender) SendMessage(
	_ context.Context,
	binding channelplatformsmodel.ChannelPlatform,
	opts twitchactions.SendMessageOpts,
) error {
	s.binding = binding
	s.opts = opts
	return nil
}

func TestTwitchChatAdapterForwardsMessageOptions(t *testing.T) {
	t.Parallel()

	sender := &fakeTwitchSender{}
	adapter := NewTwitchChatAdapter(sender)
	binding := channelplatformsmodel.ChannelPlatform{Platform: platform.PlatformTwitch}

	err := adapter.SendMessage(
		context.Background(),
		binding,
		"hello",
		"reply-id",
		ChatOptions{
			IsAnnounce:        true,
			SkipToxicityCheck: true,
			SkipRateLimits:    true,
			AnnounceColor:     bots.AnnounceColorBlue,
		},
	)
	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	if !reflect.DeepEqual(sender.binding, binding) {
		t.Errorf("expected binding %#v, got %#v", binding, sender.binding)
	}
	if sender.opts.Message != "hello" || sender.opts.ReplyParentMessageID != "reply-id" {
		t.Errorf("unexpected message options: %#v", sender.opts)
	}
	if !sender.opts.IsAnnounce || !sender.opts.SkipToxicityCheck || !sender.opts.SkipRateLimits ||
		sender.opts.AnnounceColor != bots.AnnounceColorBlue {
		t.Errorf("chat options were not forwarded: %#v", sender.opts)
	}
}

type fakeKickSender struct {
	binding channelplatformsmodel.ChannelPlatform
	message string
	replyID string
}

func (s *fakeKickSender) SendMessage(
	_ context.Context,
	binding channelplatformsmodel.ChannelPlatform,
	message string,
	replyID string,
) error {
	s.binding = binding
	s.message = message
	s.replyID = replyID
	return nil
}

func TestKickChatAdapterForwardsMessage(t *testing.T) {
	t.Parallel()

	sender := &fakeKickSender{}
	adapter := NewKickChatAdapter(sender)
	binding := channelplatformsmodel.ChannelPlatform{Platform: platform.PlatformKick}

	err := adapter.SendMessage(context.Background(), binding, "hello", "reply-id", ChatOptions{})
	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	if !reflect.DeepEqual(sender.binding, binding) {
		t.Errorf("expected binding %#v, got %#v", binding, sender.binding)
	}
	if sender.message != "hello" || sender.replyID != "reply-id" {
		t.Errorf("unexpected forwarded Kick message: %#v", sender)
	}
}

func TestNewChatRegistryRegistersTwitchAndKick(t *testing.T) {
	t.Parallel()

	registry := NewChatRegistry(nil, nil)

	for _, current := range []platform.Platform{platform.PlatformTwitch, platform.PlatformKick} {
		if _, err := registry.Require(current, platform.CapabilityChatWrite); err != nil {
			t.Errorf("expected %s chat.write adapter: %v", current, err)
		}
		if _, err := registry.Require(current, platform.CapabilityChatReply); err != nil {
			t.Errorf("expected %s chat.reply adapter: %v", current, err)
		}
	}
}

type recordingChatAdapter struct {
	platform     platform.Platform
	capabilities platform.Capabilities
	bindings     []channelplatformsmodel.ChannelPlatform
	messages     []string
	replies      []string
	err          error
}

func (a *recordingChatAdapter) Platform() platform.Platform {
	return a.platform
}

func (a *recordingChatAdapter) Capabilities() platform.Capabilities {
	return a.capabilities
}

func (a *recordingChatAdapter) SendMessage(
	_ context.Context,
	binding channelplatformsmodel.ChannelPlatform,
	message string,
	replyID string,
	_ ChatOptions,
) error {
	a.bindings = append(a.bindings, binding)
	a.messages = append(a.messages, message)
	a.replies = append(a.replies, replyID)
	return a.err
}

func TestDispatchDispatchesEnabledBindingsWithoutPlatformFilter(t *testing.T) {
	t.Parallel()

	twitch := &recordingChatAdapter{
		platform:     platform.PlatformTwitch,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
	}
	kick := &recordingChatAdapter{
		platform:     platform.PlatformKick,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
	}

	err := Dispatch(
		context.Background(),
		newRegistry(twitch, kick),
		[]channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformKick, PlatformChannelID: "kick-channel", Enabled: true},
			{Platform: platform.PlatformTwitch, PlatformChannelID: "twitch-channel", Enabled: true},
			{Platform: platform.PlatformTwitch, PlatformChannelID: "disabled-channel", Enabled: false},
		},
		nil,
		"hello",
		"reply-id",
		ChatOptions{},
	)
	if err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}
	if len(twitch.bindings) != 1 || twitch.bindings[0].PlatformChannelID != "twitch-channel" {
		t.Errorf("unexpected Twitch dispatches: %#v", twitch.bindings)
	}
	if len(kick.bindings) != 1 || kick.bindings[0].PlatformChannelID != "kick-channel" {
		t.Errorf("unexpected Kick dispatches: %#v", kick.bindings)
	}
	if twitch.messages[0] != "hello" || twitch.replies[0] != "reply-id" {
		t.Errorf("unexpected Twitch message: %#v %#v", twitch.messages, twitch.replies)
	}
}

func TestDispatchSkipsBindingsWithoutChatWriteWhenNoPlatformRequested(t *testing.T) {
	t.Parallel()

	twitch := &recordingChatAdapter{
		platform:     platform.PlatformTwitch,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
	}
	unsupported := &recordingChatAdapter{
		platform:     platform.PlatformVKVideoLive,
		capabilities: platform.Capabilities{platform.CapabilityChatRead},
	}

	err := Dispatch(
		context.Background(),
		newRegistry(twitch, unsupported),
		[]channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformVKVideoLive, Enabled: true},
			{Platform: platform.PlatformTwitch, Enabled: true},
		},
		nil,
		"hello",
		"",
		ChatOptions{},
	)
	if err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}
	if len(unsupported.bindings) != 0 {
		t.Errorf("unsupported adapter received bindings: %#v", unsupported.bindings)
	}
	if len(twitch.bindings) != 1 {
		t.Errorf("expected one Twitch dispatch, got %#v", twitch.bindings)
	}
}

func TestDispatchReturnsUnsupportedCapabilityForExplicitSelection(t *testing.T) {
	t.Parallel()

	twitch := &recordingChatAdapter{
		platform:     platform.PlatformTwitch,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
	}
	unsupported := &recordingChatAdapter{
		platform:     platform.PlatformVKVideoLive,
		capabilities: platform.Capabilities{platform.CapabilityChatRead},
	}

	err := Dispatch(
		context.Background(),
		newRegistry(twitch, unsupported),
		[]channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformTwitch, Enabled: true},
			{Platform: platform.PlatformVKVideoLive, Enabled: true},
		},
		[]platform.Platform{platform.PlatformVKVideoLive},
		"hello",
		"",
		ChatOptions{},
	)

	var unsupportedErr platform.ErrUnsupportedCapability
	if !errors.As(err, &unsupportedErr) {
		t.Fatalf("expected ErrUnsupportedCapability, got %v", err)
	}
	if unsupportedErr.Platform != platform.PlatformVKVideoLive ||
		unsupportedErr.Capability != platform.CapabilityChatWrite {
		t.Errorf("unexpected unsupported capability error: %#v", unsupportedErr)
	}
	if len(twitch.bindings) != 0 {
		t.Errorf("explicit VK selection dispatched Twitch bindings: %#v", twitch.bindings)
	}
	if len(unsupported.bindings) != 0 {
		t.Errorf("unsupported adapter received bindings: %#v", unsupported.bindings)
	}
}

func TestDispatchReturnsUnsupportedCapabilityForExplicitPlatformWithoutBinding(t *testing.T) {
	t.Parallel()

	twitch := &recordingChatAdapter{
		platform:     platform.PlatformTwitch,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
	}

	err := Dispatch(
		context.Background(),
		newRegistry(twitch),
		[]channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformTwitch, Enabled: true},
		},
		[]platform.Platform{platform.PlatformVKVideoLive},
		"hello",
		"",
		ChatOptions{},
	)

	var unsupportedErr platform.ErrUnsupportedCapability
	if !errors.As(err, &unsupportedErr) {
		t.Fatalf("expected ErrUnsupportedCapability, got %v", err)
	}
	if unsupportedErr.Platform != platform.PlatformVKVideoLive ||
		unsupportedErr.Capability != platform.CapabilityChatWrite {
		t.Errorf("unexpected unsupported capability error: %#v", unsupportedErr)
	}
	if len(twitch.bindings) != 0 {
		t.Errorf("explicit VK selection dispatched Twitch bindings: %#v", twitch.bindings)
	}
}

func TestDispatchContinuesAfterBindingFailure(t *testing.T) {
	t.Parallel()

	twitchErr := errors.New("Twitch send failed")
	twitch := &recordingChatAdapter{
		platform:     platform.PlatformTwitch,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
		err:          twitchErr,
	}
	kick := &recordingChatAdapter{
		platform:     platform.PlatformKick,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
	}

	err := Dispatch(
		context.Background(),
		newRegistry(twitch, kick),
		[]channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformTwitch, Enabled: true},
			{Platform: platform.PlatformKick, Enabled: true},
		},
		nil,
		"hello",
		"",
		ChatOptions{},
	)
	if !errors.Is(err, twitchErr) {
		t.Fatalf("expected Twitch error, got %v", err)
	}
	if len(kick.bindings) != 1 {
		t.Errorf("expected Kick dispatch after Twitch failure, got %#v", kick.bindings)
	}
}

func TestDispatchAggregatesErrorsFromEachBinding(t *testing.T) {
	t.Parallel()

	twitchErr := errors.New("Twitch send failed")
	kickErr := errors.New("Kick send failed")
	twitch := &recordingChatAdapter{
		platform:     platform.PlatformTwitch,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
		err:          twitchErr,
	}
	kick := &recordingChatAdapter{
		platform:     platform.PlatformKick,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
		err:          kickErr,
	}

	err := Dispatch(
		context.Background(),
		newRegistry(twitch, kick),
		[]channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformTwitch, Enabled: true},
			{Platform: platform.PlatformKick, Enabled: true},
		},
		nil,
		"hello",
		"",
		ChatOptions{},
	)
	if !errors.Is(err, twitchErr) || !errors.Is(err, kickErr) {
		t.Fatalf("expected both binding errors, got %v", err)
	}
	if len(twitch.bindings) != 1 || len(kick.bindings) != 1 {
		t.Errorf("expected both adapters to be called, got Twitch=%#v Kick=%#v", twitch.bindings, kick.bindings)
	}
}
