package kick

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"

	buscore "github.com/twirapp/twir/libs/bus-core"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/stretchr/testify/require"
)

func TestSplitMessage_UsesByteLimit(t *testing.T) {
	parts := splitMessage(strings.Repeat("ы", 251))

	require.Len(t, parts, 2)
	require.Len(t, []byte(parts[0]), 500)
	require.Len(t, []byte(parts[1]), 2)
}

func TestSplitMessage_PreservesUTF8Boundaries(t *testing.T) {
	parts := splitMessage(strings.Repeat("界", 167))

	require.Len(t, parts, 2)
	for _, part := range parts {
		require.True(t, utf8.ValidString(part))
		require.LessOrEqual(t, len([]byte(part)), 500)
	}
}

func TestSendMessage_RequestsKickTokenFromBus(t *testing.T) {
	t.Parallel()

	requester := &fakeBotTokenRequester{
		resp: &buscore.QueueResponse[buscoretokens.TokenResponse]{
			Data: buscoretokens.TokenResponse{AccessToken: "kick-access-token"},
		},
	}
	transport := &captureTransport{}

	client := &ChatClient{
		config: cfg.Config{},
		httpClient: &http.Client{
			Transport: transport,
		},
		requestBotToken: requester,
	}

	err := client.SendMessage(context.Background(), "42", "hello")
	require.NoError(t, err)
	require.Equal(t, 1, requester.calls)
	require.Equal(t, platformentity.PlatformKick, requester.req.Platform)
	require.Equal(t, "Bearer kick-access-token", transport.authorization)
}

type fakeBotTokenRequester struct {
	calls int
	req   buscoretokens.GetBotTokenRequest
	resp  *buscore.QueueResponse[buscoretokens.TokenResponse]
	err   error
}

func (f *fakeBotTokenRequester) Publish(ctx context.Context, data buscoretokens.GetBotTokenRequest) error {
	panic("unexpected call")
}

func (f *fakeBotTokenRequester) Request(ctx context.Context, data buscoretokens.GetBotTokenRequest) (*buscore.QueueResponse[buscoretokens.TokenResponse], error) {
	f.calls++
	f.req = data
	return f.resp, f.err
}

func (f *fakeBotTokenRequester) SubscribeGroup(queueGroup string, data buscore.QueueSubscribeCallback[buscoretokens.GetBotTokenRequest, buscoretokens.TokenResponse]) error {
	panic("unexpected call")
}

func (f *fakeBotTokenRequester) Subscribe(data buscore.QueueSubscribeCallback[buscoretokens.GetBotTokenRequest, buscoretokens.TokenResponse]) error {
	panic("unexpected call")
}

func (f *fakeBotTokenRequester) Unsubscribe() {}

type captureTransport struct {
	authorization string
}

func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.authorization = req.Header.Get("Authorization")

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"message":"OK","data":{"is_sent":true,"message_id":"msg-1"}}`)),
		Request:    req,
	}, nil
}
