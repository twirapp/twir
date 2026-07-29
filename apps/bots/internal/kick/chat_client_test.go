package kick

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
	cfg "github.com/twirapp/twir/libs/config"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/oauth"
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

	requester := &fakeBotTokenSource{token: oauth.Credential{AccessToken: "kick-access-token"}}
	transport := &captureTransport{}

	client := &ChatClient{
		config: cfg.Config{},
		httpClient: &http.Client{
			Transport: transport,
		},
		botTokens: requester,
	}

	err := client.SendMessage(context.Background(), kickBinding("42"), "hello", "")
	require.NoError(t, err)
	require.Equal(t, 1, requester.calls)
	require.Equal(t, "Bearer kick-access-token", transport.authorization)
}

type fakeBotTokenSource struct {
	calls int
	token oauth.Credential
	err   error
}

func (f *fakeBotTokenSource) Token(context.Context) (oauth.Credential, error) {
	f.calls++
	return f.token, f.err
}

type captureTransport struct {
	authorization string
	body          string
}

func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.authorization = req.Header.Get("Authorization")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	t.body = string(body)

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"message":"OK","data":{"is_sent":true,"message_id":"msg-1"}}`)),
		Request:    req,
	}, nil
}

func TestSendMessage_OmitsReplyToWhenEmpty(t *testing.T) {
	t.Parallel()

	requester := &fakeBotTokenSource{token: oauth.Credential{AccessToken: "kick-access-token"}}
	transport := &captureTransport{}

	client := &ChatClient{
		config: cfg.Config{},
		httpClient: &http.Client{
			Transport: transport,
		},
		botTokens: requester,
	}

	err := client.SendMessage(context.Background(), kickBinding("42"), "hello", "")
	require.NoError(t, err)
	require.Contains(t, transport.body, `"content":"hello"`)
	require.NotContains(t, transport.body, `"reply_to_message_id"`)
}

func TestSendMessage_IncludesReplyToWhenProvided(t *testing.T) {
	t.Parallel()

	requester := &fakeBotTokenSource{token: oauth.Credential{AccessToken: "kick-access-token"}}
	transport := &captureTransport{}

	client := &ChatClient{
		config: cfg.Config{},
		httpClient: &http.Client{
			Transport: transport,
		},
		botTokens: requester,
	}

	replyToMessageID := "opaque-reply-id-123"
	err := client.SendMessage(context.Background(), kickBinding("42"), "hello", replyToMessageID)
	require.NoError(t, err)
	require.Contains(t, transport.body, `"reply_to_message_id":"opaque-reply-id-123"`)
}

func kickBinding(platformChannelID string) channelplatformentity.ChannelPlatform {
	return channelplatformentity.ChannelPlatform{
		Platform:          platformentity.PlatformKick,
		PlatformChannelID: platformChannelID,
		Enabled:           true,
	}
}
