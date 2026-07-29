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
	"github.com/kvizyx/twitchy/helix"
	"github.com/stretchr/testify/require"
	kvinmemory "github.com/twirapp/kv/stores/inmemory"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	genericcacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestBanUsesSelectedTwitchBindingSafety(t *testing.T) {
	const (
		broadcasterID = "twitch-broadcaster"
		botID         = "twitch-bot"
		targetUserID  = "target-user"
	)

	twitchUserID := uuid.New()
	validTwitchBinding := channelplatformentity.ChannelPlatform{
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
		binding     channelplatformentity.ChannelPlatform
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
			binding: channelplatformentity.ChannelPlatform{
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
			binding: channelplatformentity.ChannelPlatform{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: broadcasterID,
				UserID:            twitchUserID,
				Enabled:           true,
			},
		},
		{
			name: "rejects malformed bot config",
			binding: channelplatformentity.ChannelPlatform{
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
			binding: channelplatformentity.ChannelPlatform{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: broadcasterID,
				UserID:            twitchUserID,
				Enabled:           true,
				BotConfig:         json.RawMessage(`{"bot_id":"twitch-bot","is_bot_mod":false}`),
			},
		},
		{
			name: "rejects banned bot config",
			binding: channelplatformentity.ChannelPlatform{
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
			client := newBanTestHelixClient(t, transport, botID)
			var userClientIDs []uuid.UUID
			var channelBotRoutes [][2]string

			actions := &TwitchActions{
				channelsByTwitchIDCache: newBanTestChannelCache(channelentity.Channel{
					Bindings: []channelplatformentity.ChannelPlatform{
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
				newChannelBotClient: func(_ context.Context, gotBotID string, gotChannelID string) (*helix.Client, error) {
					channelBotRoutes = append(channelBotRoutes, [2]string{gotBotID, gotChannelID})
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
				require.Empty(t, channelBotRoutes)
				require.Zero(t, transport.calls)
				return
			}

			require.Equal(t, []uuid.UUID{twitchUserID}, userClientIDs)
			require.Equal(t, [][2]string{{botID, broadcasterID}}, channelBotRoutes)
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

func newBanTestChannelCache(channel channelentity.Channel) *channelcache.TwitchUserIDCacher {
	return &channelcache.TwitchUserIDCacher{
		GenericCacher: genericcacher.New(
			genericcacher.Opts[channelentity.Channel]{
				KV:        kvinmemory.New(),
				KeyPrefix: "test:twitch-ban:",
				LoadFn: func(context.Context, string) (channelentity.Channel, error) {
					return channel, nil
				},
				Ttl: time.Minute,
			},
		),
	}
}

func newBanTestHelixClient(t *testing.T, transport http.RoundTripper, userID string) *helix.Client {
	t.Helper()

	client, err := helix.New(
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{
			AccessToken: "test-token",
			ClientID:    "test-client",
			TokenClass:  helix.TokenClassUser,
			UserID:      userID,
			Scopes:      []helix.AuthorizationScope{helix.ScopeModeratorManageBannedUsers},
		}),
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
