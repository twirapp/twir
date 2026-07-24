package auth

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestVKVideoBotSetupRejectsUnauthorizedAndNonAdminUsers(t *testing.T) {
	// Given
	for _, admin := range []usersmodel.User{
		{ID: uuid.New(), IsBotAdmin: false},
		{ID: uuid.New(), IsBotAdmin: true, IsBanned: true},
	} {
		fixture := newVKVideoBotSetupFixture(admin)

		// When
		_, err := fixture.auth.StartVKVideoBotSetup(context.Background())

		// Then
		if err == nil {
			t.Fatalf("setup link for %+v succeeded", admin)
		}
	}
}

func TestVKVideoBotSetupBindsOpaqueStateToTheLiveAdmin(t *testing.T) {
	// Given
	admin := usersmodel.User{ID: uuid.New(), IsBotAdmin: true}
	fixture := newVKVideoBotSetupFixture(admin)
	link, err := fixture.auth.StartVKVideoBotSetup(context.Background())
	if err != nil {
		t.Fatalf("start setup: %v", err)
	}
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse setup URL: %v", err)
	}
	state := parsed.Query().Get("state")
	if len(state) < 32 {
		t.Fatalf("state length = %d, want at least 32", len(state))
	}
	otherAdmin := usersmodel.User{ID: uuid.New(), IsBotAdmin: true}
	fixture.users.users[otherAdmin.ID] = otherAdmin
	fixture.sessions.userID = otherAdmin.ID

	// When
	err = fixture.auth.CompleteVKVideoBotSetup(context.Background(), "code", state)

	// Then
	if !errors.Is(err, ErrVKVideoBotSetupStateInvalid) {
		t.Fatalf("completion error = %v, want invalid state", err)
	}
	if fixture.provider.exchangeCalls != 0 || fixture.bots.upsertCalls != 0 {
		t.Fatalf("mismatched session exchanged or persisted a bot")
	}
}

func TestVKVideoBotSetupConsumesStateBeforeExchangeAndDoesNotMutateSession(t *testing.T) {
	// Given
	admin := usersmodel.User{ID: uuid.New(), IsBotAdmin: true}
	fixture := newVKVideoBotSetupFixture(admin)
	link, err := fixture.auth.StartVKVideoBotSetup(context.Background())
	if err != nil {
		t.Fatalf("start setup: %v", err)
	}
	state := mustVKVideoBotSetupState(t, link)

	// When
	err = fixture.auth.CompleteVKVideoBotSetup(context.Background(), "code", state)
	if err != nil {
		t.Fatalf("complete setup: %v", err)
	}
	replayErr := fixture.auth.CompleteVKVideoBotSetup(context.Background(), "code", state)

	// Then
	if !errors.Is(replayErr, ErrVKVideoBotSetupStateInvalid) {
		t.Fatalf("replay error = %v, want invalid state", replayErr)
	}
	if fixture.provider.exchangeCalls != 1 {
		t.Fatalf("exchange calls = %d, want 1", fixture.provider.exchangeCalls)
	}
	if fixture.sessions.setIdentityCalls != 0 || fixture.sessions.setPlatformCalls != 0 || fixture.sessions.setDashboardCalls != 0 {
		t.Fatalf("setup mutated browser session: %+v", fixture.sessions)
	}
}

func TestVKVideoBotSetupCreatesOrReplacesSingletonAndBackfillsBindings(t *testing.T) {
	// Given
	admin := usersmodel.User{ID: uuid.New(), IsBotAdmin: true}
	fixture := newVKVideoBotSetupFixture(admin)
	firstState := startVKVideoBotSetup(t, fixture)

	// When
	if err := fixture.auth.CompleteVKVideoBotSetup(context.Background(), "first-code", firstState); err != nil {
		t.Fatalf("complete first setup: %v", err)
	}
	firstBotID := fixture.bots.bot.ID
	fixture.provider.tokens.AccessToken = "replacement-access"
	secondState := startVKVideoBotSetup(t, fixture)
	if err := fixture.auth.CompleteVKVideoBotSetup(context.Background(), "replacement-code", secondState); err != nil {
		t.Fatalf("complete replacement: %v", err)
	}

	// Then
	if fixture.bots.bot.ID != firstBotID || fixture.bots.upsertCalls != 2 {
		t.Fatalf("singleton replacement = %#v, upserts = %d", fixture.bots.bot, fixture.bots.upsertCalls)
	}
	if fixture.bots.lockCalls != 2 || fixture.bindings.assignCalls != 2 || fixture.bindings.assignedUserID != fixture.bots.bot.VKUserID {
		t.Fatalf("replacement did not lock and backfill: locks=%d assignments=%d user=%s", fixture.bots.lockCalls, fixture.bindings.assignCalls, fixture.bindings.assignedUserID)
	}
	if strings.Contains(fixture.bots.bot.EncryptedAccessToken, "replacement-access") {
		t.Fatal("access token was persisted without encryption")
	}
}

func TestVKVideoBotBindingConfigRequiresAndUsesTheSingleton(t *testing.T) {
	// Given
	admin := usersmodel.User{ID: uuid.New(), IsBotAdmin: true}
	fixture := newVKVideoBotSetupFixture(admin)

	// When
	_, missingErr := fixture.auth.vkVideoBotBindingConfig(context.Background())
	botUserID := uuid.New()
	fixture.bots.bot.ID = uuid.New()
	fixture.bots.bot.VKUserID = botUserID
	config, configuredErr := fixture.auth.vkVideoBotBindingConfig(context.Background())

	// Then
	if !errors.Is(missingErr, ErrVKVideoBotNotConfigured) {
		t.Fatalf("missing singleton error = %v", missingErr)
	}
	if configuredErr != nil || config.BotUserID == nil || *config.BotUserID != botUserID {
		t.Fatalf("future binding config = %#v, error = %v", config, configuredErr)
	}
	if fixture.bots.lockCalls != 2 {
		t.Fatalf("singleton lock calls = %d, want 2", fixture.bots.lockCalls)
	}
}

func startVKVideoBotSetup(t *testing.T, fixture *vkVideoBotSetupFixture) string {
	t.Helper()
	link, err := fixture.auth.StartVKVideoBotSetup(context.Background())
	if err != nil {
		t.Fatalf("start setup: %v", err)
	}
	return mustVKVideoBotSetupState(t, link)
}

func mustVKVideoBotSetupState(t *testing.T, link string) string {
	t.Helper()
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse setup URL: %v", err)
	}
	state := parsed.Query().Get("state")
	if state == "" {
		t.Fatal("setup URL has no state")
	}
	return state
}

func testVKVideoBotConfig() cfg.Config {
	return cfg.Config{TokensCipherKey: "pnyfwfiulmnqlhkvixaeligpprcnlyke", VKVideoClientID: "client", VKVideoClientSecret: "secret"}
}

var _ = platformentity.PlatformVKVideoLive
