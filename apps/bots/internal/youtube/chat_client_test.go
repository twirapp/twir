package youtube

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	buscore "github.com/twirapp/twir/libs/bus-core"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func TestSplitMessage_SplitsAt200Runes(t *testing.T) {
	// Given
	message := strings.Repeat("界", 201)

	// When
	parts := splitMessage(message)

	// Then
	require.Equal(t, []string{strings.Repeat("界", 200), "界"}, parts)
}

func TestChatClient_SendMessage_postsOfficialLiveChatMessage(t *testing.T) {
	// Given
	ownerID := uuid.New()
	cache := newFakeLiveChatCache()
	ownerRequester := &fakeUserTokenRequester{
		response: &buscore.QueueResponse[buscoretokens.TokenResponse]{
			Data: buscoretokens.TokenResponse{AccessToken: "streamer-access-token"},
		},
	}
	botRequester := &fakeBotTokenRequester{
		response: &buscore.QueueResponse[buscoretokens.TokenResponse]{
			Data: buscoretokens.TokenResponse{AccessToken: "bot-access-token"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/liveBroadcasts":
			require.Equal(t, http.MethodGet, request.Method)
			require.Equal(t, "Bearer streamer-access-token", request.Header.Get("Authorization"))
			require.Equal(t, "snippet", request.URL.Query().Get("part"))
			require.Equal(t, "true", request.URL.Query().Get("mine"))
			require.Equal(t, "active", request.URL.Query().Get("broadcastStatus"))
			writer.Header().Set("Content-Type", "application/json")
			_, err := writer.Write([]byte(`{"items":[{"snippet":{"liveChatId":"live-chat-42"}}]}`))
			require.NoError(t, err)
		case "/liveChat/messages":
			require.Equal(t, http.MethodPost, request.Method)
			require.Equal(t, "Bearer bot-access-token", request.Header.Get("Authorization"))
			require.Equal(t, "application/json", request.Header.Get("Content-Type"))
			require.Equal(t, "snippet", request.URL.Query().Get("part"))

			body, err := io.ReadAll(request.Body)
			require.NoError(t, err)
			require.JSONEq(t, `{"snippet":{"liveChatId":"live-chat-42","type":"textMessageEvent","textMessageDetails":{"messageText":"hello from bot"}}}`, string(body))
			writer.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected request path %q", request.URL.Path)
		}
	}))
	defer server.Close()

	client := &ChatClient{
		apiBaseURL:       server.URL,
		httpClient:       server.Client(),
		liveChatIDCache:  cache,
		requestUserToken: ownerRequester,
		requestBotToken:  botRequester,
	}
	binding := channelplatformentity.ChannelPlatform{
		ChannelID:         uuid.New(),
		Platform:          platformentity.PlatformYouTube,
		PlatformChannelID: "youtube-channel-42",
		UserID:            ownerID,
	}

	// When
	err := client.SendMessage(context.Background(), binding, "hello from bot")

	// Then
	require.NoError(t, err)
	require.Equal(t, buscoretokens.GetUserTokenRequest{UserId: ownerID}, ownerRequester.request)
	require.Equal(t, buscoretokens.GetBotTokenRequest{Platform: platformentity.PlatformYouTube}, botRequester.request)
	require.Len(t, cache.sets, 1)
	require.Equal(t, "youtube:livechat_id:"+binding.ChannelID.String(), cache.sets[0].key)
	require.Equal(t, "live-chat-42", cache.sets[0].value)
}

type fakeUserTokenRequester struct {
	request  buscoretokens.GetUserTokenRequest
	response *buscore.QueueResponse[buscoretokens.TokenResponse]
	err      error
}

func (f *fakeUserTokenRequester) Request(
	_ context.Context,
	request buscoretokens.GetUserTokenRequest,
) (*buscore.QueueResponse[buscoretokens.TokenResponse], error) {
	f.request = request
	return f.response, f.err
}

type fakeBotTokenRequester struct {
	request  buscoretokens.GetBotTokenRequest
	response *buscore.QueueResponse[buscoretokens.TokenResponse]
	err      error
}

func (f *fakeBotTokenRequester) Request(
	_ context.Context,
	request buscoretokens.GetBotTokenRequest,
) (*buscore.QueueResponse[buscoretokens.TokenResponse], error) {
	f.request = request
	return f.response, f.err
}

type fakeLiveChatCache struct {
	values map[string]string
	sets   []fakeLiveChatCacheSet
}

type fakeLiveChatCacheSet struct {
	key   string
	value string
}

func newFakeLiveChatCache() *fakeLiveChatCache {
	return &fakeLiveChatCache{values: make(map[string]string)}
}

func (f *fakeLiveChatCache) Get(_ context.Context, key string) kv.Valuer {
	value, ok := f.values[key]
	if !ok {
		return fakeCacheValue{err: kv.ErrKeyNil}
	}

	return fakeCacheValue{value: value}
}

func (f *fakeLiveChatCache) Set(
	_ context.Context,
	key string,
	value any,
	options ...kvoptions.Option,
) error {
	stringValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string cache value, got %T", value)
	}

	f.values[key] = stringValue
	f.sets = append(f.sets, fakeLiveChatCacheSet{key: key, value: stringValue})
	return nil
}

type fakeCacheValue struct {
	value string
	err   error
}

func (v fakeCacheValue) Int() (int64, error) { return 0, v.err }

func (v fakeCacheValue) String() (string, error) { return v.value, v.err }

func (v fakeCacheValue) Bytes() ([]byte, error) { return []byte(v.value), v.err }

func (v fakeCacheValue) Bool() (bool, error) { return false, v.err }

func (v fakeCacheValue) Float() (float64, error) { return 0, v.err }

func (v fakeCacheValue) Scan(any) error { return v.err }

func (v fakeCacheValue) Err() error { return v.err }
