package twitch

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	"github.com/redis/go-redis/v9"
)

type twitchyCaptureTransport struct {
	calls int
	paths []string
}

func (t *twitchyCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.calls++
	t.paths = append(t.paths, req.URL.Path)

	body := `{"data":[]}`
	switch req.URL.Path {
	case "/helix/users":
		body = `{"data":[{"id":"user-id","login":"viewer","display_name":"Viewer","profile_image_url":"avatar"}]}`
	case "/helix/games":
		body = `{"data":[{"id":"game-id","name":"Game","box_art_url":"art"}]}`
	case "/helix/channel_points/custom_rewards":
		body = `{"data":[{"id":"reward-id","title":"Reward","cost":100,"image":null,"default_image":{"url_1x":"default-1x","url_2x":"default-2x","url_4x":"default-4x"}}]}`
	case "/helix/channels/followers":
		body = `{"data":[{"user_id":"one"},{"user_id":"two"}]}`
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}, nil
}

func newTwitchyCacheTestClient(t *testing.T, transport http.RoundTripper) *helix.Client {
	t.Helper()

	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{
			AccessToken: "test-token",
			ClientID:    "test-client",
			TokenClass:  helix.TokenClassUser,
			UserID:      "broadcaster",
			Scopes: []helix.AuthorizationScope{
				helix.ScopeChannelReadRedemptions,
				helix.ScopeChannelManageRedemptions,
				helix.ScopeModeratorReadFollowers,
			},
		}),
	)
	if err != nil {
		t.Fatalf("new Twitchy client: %v", err)
	}

	return client
}

func newTwitchyCacheTestService(t *testing.T, transport *twitchyCaptureTransport) *CachedTwitchClient {
	t.Helper()

	redisServer := miniredis.RunT(t)
	client := newTwitchyCacheTestClient(t, transport)
	return &CachedTwitchClient{
		client: client,
		redis:  redis.NewClient(&redis.Options{Addr: redisServer.Addr()}),
		newUserClient: func(context.Context, uuid.UUID) (*helix.Client, error) {
			return client, nil
		},
	}
}

func TestCachedTwitchClientMapsAppReadsAndCachesUsers(t *testing.T) {
	// Given
	transport := &twitchyCaptureTransport{}
	cache := newTwitchyCacheTestService(t, transport)

	// When
	user, err := cache.GetUserById(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	game, err := cache.GetGame(context.Background(), "game-id")
	if err != nil {
		t.Fatalf("get game: %v", err)
	}
	_, err = cache.GetUserById(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("get cached user: %v", err)
	}

	// Then
	if user.Login != "viewer" || user.ProfileImageURL != "avatar" {
		t.Fatalf("user = %#v, want mapped Twitch user", user)
	}
	if game.Name != "Game" || game.BoxArtURL != "art" {
		t.Fatalf("game = %#v, want mapped Twitch game", game)
	}
	if transport.calls != 2 {
		t.Fatalf("HTTP calls = %d, want 2 after cached user read", transport.calls)
	}
}

func TestCachedTwitchClientMapsRewardsAndCountsFollowersWithBroadcasterClient(t *testing.T) {
	// Given
	transport := &twitchyCaptureTransport{}
	cache := newTwitchyCacheTestService(t, transport)
	broadcasterID := uuid.New()

	// When
	rewards, err := cache.GetChannelRewards(context.Background(), broadcasterID, "broadcaster")
	if err != nil {
		t.Fatalf("get rewards: %v", err)
	}
	followers, err := cache.GetChannelFollowersCountByChannelId(context.Background(), broadcasterID, "broadcaster")
	if err != nil {
		t.Fatalf("get followers: %v", err)
	}

	// Then
	if len(rewards) != 1 || rewards[0].Image == nil || rewards[0].Image.URL1x != "default-1x" {
		t.Fatalf("rewards = %#v, want default-image-mapped reward", rewards)
	}
	if followers != 2 {
		t.Fatalf("followers = %d, want 2", followers)
	}
	if transport.calls != 2 || transport.paths[0] != "/helix/channel_points/custom_rewards" || transport.paths[1] != "/helix/channels/followers" {
		t.Fatalf("requests = %#v, want rewards then followers", transport.paths)
	}
}
