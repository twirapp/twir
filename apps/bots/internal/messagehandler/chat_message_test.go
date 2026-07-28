package messagehandler

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
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

type chatMessageEmoteQueue[Req, Res any] struct {
	mu         sync.Mutex
	requests   []Req
	response   *buscore.QueueResponse[Res]
	requestErr error
}

func (q *chatMessageEmoteQueue[Req, Res]) Publish(context.Context, Req) error { return nil }

func (q *chatMessageEmoteQueue[Req, Res]) Request(
	_ context.Context,
	request Req,
) (*buscore.QueueResponse[Res], error) {
	q.mu.Lock()
	q.requests = append(q.requests, request)
	q.mu.Unlock()

	return q.response, q.requestErr
}

func (q *chatMessageEmoteQueue[Req, Res]) SubscribeGroup(string, buscore.QueueSubscribeCallback[Req, Res]) error {
	return nil
}

func (q *chatMessageEmoteQueue[Req, Res]) Subscribe(buscore.QueueSubscribeCallback[Req, Res]) error {
	return nil
}

func (q *chatMessageEmoteQueue[Req, Res]) Unsubscribe() {}

func (q *chatMessageEmoteQueue[Req, Res]) requestSnapshot() []Req {
	q.mu.Lock()
	defer q.mu.Unlock()

	return append([]Req(nil), q.requests...)
}

func TestChatMessageCountEmotes(t *testing.T) {
	tests := []struct {
		name                string
		message             generic.ChatMessage
		globalEmotes        []emotes_cacher.Emote
		channelEmotes       []emotes_cacher.Emote
		want                map[string]int
		wantGlobalRequests  []emotes_cacher.GetGlobalEmotesRequest
		wantChannelRequests []emotes_cacher.GetChannelEmotesRequest
	}{
		{
			name:    "nil message returns nil without cache requests",
			message: generic.ChatMessage{},
			want:    nil,
		},
		{
			name: "twitch counts native and third-party emotes",
			message: generic.ChatMessage{
				Platform:          string(platform.PlatformTwitch),
				PlatformChannelID: "twitch-provider-id",
				Message: &generic.ChatMessageMessage{
					Text: "Kappa Pog Party",
					Fragments: []generic.ChatMessageMessageFragment{
						{
							Type:  generic.FragmentType_EMOTE,
							Text:  "Kappa",
							Emote: &generic.ChatMessageMessageFragmentEmote{ID: "native-kappa"},
						},
					},
				},
			},
			globalEmotes:       []emotes_cacher.Emote{{Name: "Pog"}},
			channelEmotes:      []emotes_cacher.Emote{{Name: "Party"}},
			want:               map[string]int{"Kappa": 1, "Pog": 1, "Party": 1},
			wantGlobalRequests: []emotes_cacher.GetGlobalEmotesRequest{{}},
			wantChannelRequests: []emotes_cacher.GetChannelEmotesRequest{{
				Platform:  platform.PlatformTwitch,
				ChannelID: "twitch-provider-id",
			}},
		},
		{
			name: "kick counts native and third-party emotes by provider channel",
			message: generic.ChatMessage{
				Platform:          string(platform.PlatformKick),
				ChannelID:         "internal-channel-uuid",
				PlatformChannelID: "kick-provider-id",
				Message: &generic.ChatMessageMessage{
					Text: "vahui LUL KEKW",
					Fragments: []generic.ChatMessageMessageFragment{
						{
							Type:  generic.FragmentType_EMOTE,
							Text:  "vahui",
							Emote: &generic.ChatMessageMessageFragmentEmote{ID: "native-kick"},
						},
					},
				},
			},
			globalEmotes:       []emotes_cacher.Emote{{Name: "LUL"}},
			channelEmotes:      []emotes_cacher.Emote{{Name: "KEKW"}},
			want:               map[string]int{"vahui": 1, "LUL": 1, "KEKW": 1},
			wantGlobalRequests: []emotes_cacher.GetGlobalEmotesRequest{{}},
			wantChannelRequests: []emotes_cacher.GetChannelEmotesRequest{{
				Platform:  platform.PlatformKick,
				ChannelID: "kick-provider-id",
			}},
		},
		{
			name: "native fragments are not counted again by token matching",
			message: generic.ChatMessage{
				Platform:          string(platform.PlatformTwitch),
				PlatformChannelID: "twitch-provider-id",
				Message: &generic.ChatMessageMessage{
					Text: "Kappa Kappa Pog",
					Fragments: []generic.ChatMessageMessageFragment{
						{
							Type:  generic.FragmentType_EMOTE,
							Text:  "Kappa",
							Emote: &generic.ChatMessageMessageFragmentEmote{ID: "native-kappa-1"},
						},
						{
							Type:  generic.FragmentType_EMOTE,
							Text:  "Kappa",
							Emote: &generic.ChatMessageMessageFragmentEmote{ID: "native-kappa-2"},
						},
					},
				},
			},
			globalEmotes:       []emotes_cacher.Emote{{Name: "Pog"}},
			want:               map[string]int{"Kappa": 2, "Pog": 1},
			wantGlobalRequests: []emotes_cacher.GetGlobalEmotesRequest{{}},
			wantChannelRequests: []emotes_cacher.GetChannelEmotesRequest{{
				Platform:  platform.PlatformTwitch,
				ChannelID: "twitch-provider-id",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globalQueue := &chatMessageEmoteQueue[
				emotes_cacher.GetGlobalEmotesRequest,
				emotes_cacher.Response,
			]{
				response: &buscore.QueueResponse[emotes_cacher.Response]{
					Data: emotes_cacher.Response{Emotes: tt.globalEmotes},
				},
			}
			channelQueue := &chatMessageEmoteQueue[
				emotes_cacher.GetChannelEmotesRequest,
				emotes_cacher.Response,
			]{
				response: &buscore.QueueResponse[emotes_cacher.Response]{
					Data: emotes_cacher.Response{Emotes: tt.channelEmotes},
				},
			}

			handler := &MessageHandler{}
			if tt.message.Message != nil {
				bus := buscore.NewNatsBus(nil)
				bus.EmotesCacher.GetGlobalEmotes = globalQueue
				bus.EmotesCacher.GetChannelEmotes = channelQueue
				handler.twirBus = bus
			}

			got, err := handler.chatMessageCountEmotes(context.Background(), tt.message)
			if err != nil {
				t.Fatalf("chatMessageCountEmotes() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("chatMessageCountEmotes() = %#v, want %#v", got, tt.want)
			}
			if got := globalQueue.requestSnapshot(); !reflect.DeepEqual(got, tt.wantGlobalRequests) {
				t.Fatalf("global emote requests = %#v, want %#v", got, tt.wantGlobalRequests)
			}
			if got := channelQueue.requestSnapshot(); !reflect.DeepEqual(got, tt.wantChannelRequests) {
				t.Fatalf("channel emote requests = %#v, want %#v", got, tt.wantChannelRequests)
			}
		})
	}
}
