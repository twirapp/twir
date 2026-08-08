package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	entity "github.com/twirapp/twir/libs/entities/youtube_bot"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	youtubebotsrepo "github.com/twirapp/twir/libs/repositories/youtube_bots"
)

func TestCompleteYouTubeBotSetupReassignsBindingsAndResubscribesAffectedChannels(t *testing.T) {
	// Given
	admin := usersmodel.User{ID: uuid.New(), IsBotAdmin: true}
	sessions := &vkVideoBotSetupSessions{userID: admin.ID}
	users := &vkVideoBotUsersRepositoryFake{users: map[uuid.UUID]usersmodel.User{admin.ID: admin}}
	bindings := &youtubeBotBindingRepositoryFake{channelIDs: []uuid.UUID{uuid.New(), uuid.New()}}
	publisher := &oauthEventSubPublisher{}
	auth := &Auth{
		config:               cfg.Config{TokensCipherKey: "pnyfwfiulmnqlhkvixaeligpprcnlyke"},
		sessions:             sessions,
		usersRepo:            users,
		youtubeBotProvider:   &youtubeBotSetupProviderFake{},
		youtubeBotsRepo:      &youtubeBotRepositoryFake{},
		channelPlatformsRepo: bindings,
		transactionRunner:    &vkVideoBotTransactionFake{},
		kv:                   &vkVideoBotKVFake{values: make(map[string][]byte)},
		eventSubPublisher:    publisher,
	}
	link, err := auth.StartYouTubeBotSetup(context.Background())
	if err != nil {
		t.Fatalf("start YouTube bot setup: %v", err)
	}
	parsedLink, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse setup URL: %v", err)
	}

	// When
	err = auth.CompleteYouTubeBotSetup(context.Background(), "code", parsedLink.Query().Get("state"))

	// Then
	if err != nil {
		t.Fatalf("complete YouTube bot setup: %v", err)
	}
	if bindings.assignCalls != 1 || bindings.botUserID == uuid.Nil {
		t.Fatalf("assignment calls = %d, bot user ID = %s", bindings.assignCalls, bindings.botUserID)
	}
	if len(publisher.requests) != len(bindings.channelIDs) {
		t.Fatalf("published %d subscriptions, want %d", len(publisher.requests), len(bindings.channelIDs))
	}
	for _, request := range publisher.requests {
		if request.Platform != platform.PlatformYouTube {
			t.Fatalf("published platform = %s, want %s", request.Platform, platform.PlatformYouTube)
		}
	}
}

func TestYouTubeBotCallbackRedirectsToAdminPanelAfterSetup(t *testing.T) {
	// Given
	admin := usersmodel.User{ID: uuid.New(), IsBotAdmin: true}
	sessions := &vkVideoBotSetupSessions{userID: admin.ID}
	_, api := humatest.New(t)
	auth := &Auth{
		config:               cfg.Config{TokensCipherKey: "pnyfwfiulmnqlhkvixaeligpprcnlyke"},
		sessions:             sessions,
		usersRepo:            &vkVideoBotUsersRepositoryFake{users: map[uuid.UUID]usersmodel.User{admin.ID: admin}},
		youtubeBotProvider:   &youtubeBotSetupProviderFake{},
		youtubeBotsRepo:      &youtubeBotRepositoryFake{},
		channelPlatformsRepo: &youtubeBotBindingRepositoryFake{},
		transactionRunner:    &vkVideoBotTransactionFake{},
		kv:                   &vkVideoBotKVFake{values: make(map[string][]byte)},
	}
	huma.Get(api, "/auth/youtube/bot-callback", func(ctx context.Context, input *struct {
		Code  string `query:"code"`
		State string `query:"state"`
	}) (*youtubeBotCallbackOutput, error) {
		return auth.completeYouTubeBotCallback(ctx, input.Code, input.State)
	})
	link, err := auth.StartYouTubeBotSetup(context.Background())
	if err != nil {
		t.Fatalf("start YouTube bot setup: %v", err)
	}
	parsedLink, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse setup URL: %v", err)
	}

	// When
	response := api.Get("/auth/youtube/bot-callback?code=code&state=" + url.QueryEscape(parsedLink.Query().Get("state")))

	// Then
	if response.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusFound)
	}
	if location := response.Header().Get("Location"); location != "/dashboard/admin" {
		t.Fatalf("Location = %q, want %q", location, "/dashboard/admin")
	}
}

type youtubeBotSetupProviderFake struct{}

func (youtubeBotSetupProviderFake) GetBotSetupAuthURL(state, _ string) string {
	return "https://youtube.example.test/authorize?state=" + url.QueryEscape(state)
}

func (youtubeBotSetupProviderFake) ExchangeBotSetupCode(context.Context, string, string) (*appplatform.PlatformTokens, error) {
	return &appplatform.PlatformTokens{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 3600, Scopes: []string{"youtube.read"}}, nil
}

func (youtubeBotSetupProviderFake) GetUser(context.Context, string) (*appplatform.PlatformUser, error) {
	return &appplatform.PlatformUser{ID: "youtube-bot", Login: "youtube-bot", DisplayName: "YouTube Bot"}, nil
}

type youtubeBotRepositoryFake struct {
	bot entity.YouTubeBot
}

func (r *youtubeBotRepositoryFake) Get(context.Context) (entity.YouTubeBot, error) {
	if r.bot.ID == uuid.Nil {
		return entity.Nil, youtubebotsrepo.ErrNotFound
	}
	return r.bot, nil
}

func (*youtubeBotRepositoryFake) Lock(context.Context) error { return nil }

func (r *youtubeBotRepositoryFake) Upsert(_ context.Context, input youtubebotsrepo.UpsertInput) (entity.YouTubeBot, error) {
	if r.bot.ID == uuid.Nil {
		r.bot.ID = uuid.New()
	}
	r.bot.YouTubeUserID = input.YouTubeUserID
	return r.bot, nil
}

func (*youtubeBotRepositoryFake) Update(context.Context, youtubebotsrepo.UpdateInput) (entity.YouTubeBot, error) {
	return entity.Nil, errors.New("unexpected update")
}

type youtubeBotBindingRepositoryFake struct {
	channelplatformsrepo.Repository
	botUserID   uuid.UUID
	channelIDs  []uuid.UUID
	assignCalls int
}

func (r *youtubeBotBindingRepositoryFake) AssignYouTubeBot(_ context.Context, botUserID uuid.UUID) ([]uuid.UUID, error) {
	r.assignCalls++
	r.botUserID = botUserID
	return r.channelIDs, nil
}

var _ youtubebotsrepo.Repository = (*youtubeBotRepositoryFake)(nil)
