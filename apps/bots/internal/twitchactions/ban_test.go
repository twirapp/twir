package twitchactions

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/stretchr/testify/require"
	kvinmemory "github.com/twirapp/kv/stores/inmemory"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	genericcacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestBanUsesSelectedTwitchBindingSafety(t *testing.T) {
	const (
		broadcasterID = "twitch-broadcaster"
		botID         = "twitch-bot"
		targetUserID  = "target-user"
	)

	twitchUserID := uuid.New()
	validTwitchBinding := channelplatformsmodel.ChannelPlatform{
		Platform:          platform.PlatformTwitch,
		PlatformChannelID: broadcasterID,
		UserID:            twitchUserID,
		Enabled:           true,
		BotConfig: json.RawMessage(
			`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":false}`,
		),
	}

	tests := []struct {
		name        string
		binding     channelplatformsmodel.ChannelPlatform
		wantErr     bool
		wantAPICall bool
	}{
		{
			name:        "uses selected Twitch binding when Kick comes first",
			binding:     validTwitchBinding,
			wantAPICall: true,
		},
		{
			name: "rejects disabled binding",
			binding: channelplatformsmodel.ChannelPlatform{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: broadcasterID,
				UserID:            twitchUserID,
				Enabled:           false,
				BotConfig:         validTwitchBinding.BotConfig,
			},
			wantErr: true,
		},
		{
			name: "rejects missing bot config",
			binding: channelplatformsmodel.ChannelPlatform{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: broadcasterID,
				UserID:            twitchUserID,
				Enabled:           true,
			},
		},
		{
			name: "rejects malformed bot config",
			binding: channelplatformsmodel.ChannelPlatform{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: broadcasterID,
				UserID:            twitchUserID,
				Enabled:           true,
				BotConfig:         json.RawMessage(`{`),
			},
			wantErr: true,
		},
		{
			name: "rejects non moderator bot config",
			binding: channelplatformsmodel.ChannelPlatform{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: broadcasterID,
				UserID:            twitchUserID,
				Enabled:           true,
				BotConfig:         json.RawMessage(`{"bot_id":"twitch-bot","is_bot_mod":false}`),
			},
		},
		{
			name: "rejects banned bot config",
			binding: channelplatformsmodel.ChannelPlatform{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: broadcasterID,
				UserID:            twitchUserID,
				Enabled:           true,
				BotConfig:         json.RawMessage(`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &banCaptureTransport{}
			client := newBanTestHelixClient(t, transport)
			var userClientIDs []uuid.UUID
			var botClientIDs []string

			actions := &TwitchActions{
				channelsByTwitchIDCache: newBanTestChannelCache(channelsmodel.Channel{
					Bindings: []channelplatformsmodel.ChannelPlatform{
						{
							Platform:          platform.PlatformKick,
							PlatformChannelID: "kick-channel",
							UserID:            uuid.New(),
							Enabled:           true,
						},
						tt.binding,
					},
				}),
				newUserClient: func(_ context.Context, userID uuid.UUID) (*helix.Client, error) {
					userClientIDs = append(userClientIDs, userID)
					return client, nil
				},
				newBotClient: func(_ context.Context, gotBotID string) (*helix.Client, error) {
					botClientIDs = append(botClientIDs, gotBotID)
					return client, nil
				},
			}

			err := actions.Ban(
				context.Background(),
				BanOpts{
					BroadcasterID: broadcasterID,
					ModeratorID:   "caller-supplied-bot",
					UserID:        targetUserID,
					Reason:        "test reason",
					Duration:      60,
				},
			)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if !tt.wantAPICall {
				require.Empty(t, userClientIDs)
				require.Empty(t, botClientIDs)
				require.Zero(t, transport.calls)
				return
			}

			require.Equal(t, []uuid.UUID{twitchUserID}, userClientIDs)
			require.Equal(t, []string{botID}, botClientIDs)
			require.Equal(t, 1, transport.calls)
			require.Equal(t, http.MethodPost, transport.method)
			require.Equal(t, "/helix/moderation/bans", transport.path)
			require.Equal(t, broadcasterID, transport.query.Get("broadcaster_id"))
			require.Equal(t, botID, transport.query.Get("moderator_id"))

			var body struct {
				Data struct {
					Duration int    `json:"duration"`
					Reason   string `json:"reason"`
					UserID   string `json:"user_id"`
				} `json:"data"`
			}
			require.NoError(t, json.Unmarshal([]byte(transport.body), &body))
			require.Equal(t, 60, body.Data.Duration)
			require.Equal(t, "test reason", body.Data.Reason)
			require.Equal(t, targetUserID, body.Data.UserID)
		})
	}
}

func newBanTestChannelCache(channel channelsmodel.Channel) *channelcache.TwitchUserIDCacher {
	return &channelcache.TwitchUserIDCacher{
		GenericCacher: genericcacher.New(
			genericcacher.Opts[channelsmodel.Channel]{
				KV:        kvinmemory.New(),
				KeyPrefix: "test:twitch-ban:",
				LoadFn: func(context.Context, string) (channelsmodel.Channel, error) {
					return channel, nil
				},
				Ttl: time.Minute,
			},
		),
	}
}

func newBanTestHelixClient(t *testing.T, transport http.RoundTripper) *helix.Client {
	t.Helper()

	client, err := helix.NewClient(
		&helix.Options{
			ClientID: "test-client",
			HTTPClient: &http.Client{
				Transport: transport,
			},
		},
	)
	require.NoError(t, err)

	return client
}

type banCaptureTransport struct {
	calls  int
	method string
	path   string
	query  url.Values
	body   string
}

func (t *banCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.calls++
	t.method = req.Method
	t.path = req.URL.Path
	t.query = req.URL.Query()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	t.body = string(body)

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
		Request:    req,
	}, nil
}
