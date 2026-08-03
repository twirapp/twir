package postgres

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	"github.com/twirapp/twir/libs/entities/platform"
	channelspgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	mcpOAuth "github.com/twirapp/twir/libs/repositories/mcp_oauth"
	"github.com/twirapp/twir/libs/repositories/users"
	userspgx "github.com/twirapp/twir/libs/repositories/users/pgx"
)

const integrationDatabaseURLEnv = "TWIR_MCP_OAUTH_TEST_DATABASE_URL"

type integrationFixture struct {
	pool      *pgxpool.Pool
	clientID  string
	channelID uuid.UUID
	userID    uuid.UUID
}

func newIntegrationFixture(t *testing.T, ctx context.Context) (*Pgx, integrationFixture) {
	t.Helper()

	databaseURL := os.Getenv(integrationDatabaseURLEnv)
	if databaseURL == "" {
		t.Skipf("set %s to run PostgreSQL integration tests", integrationDatabaseURLEnv)
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping PostgreSQL: %v", err)
	}

	user, err := userspgx.New(userspgx.Opts{PgxPool: pool}).Create(ctx, users.CreateInput{
		Platform:          platform.PlatformTwitch,
		PlatformID:        uuid.NewString(),
		IsBotAdmin:        false,
		IsBanned:          false,
		HideOnLandingPage: true,
		Login:             "mcp-oauth-test-" + uuid.NewString(),
		DisplayName:       "MCP OAuth Test",
		Avatar:            "",
	})
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, user.ID); err != nil {
			t.Errorf("delete test user: %v", err)
		}
	})

	channel, err := channelspgx.New(channelspgx.Opts{PgxPool: pool}).Create(ctx)
	if err != nil {
		t.Fatalf("create test channel: %v", err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM channels WHERE id = $1`, channel.ID); err != nil {
			t.Errorf("delete test channel: %v", err)
		}
	})

	clientID := testClientID()
	repository := New(Opts{PgxPool: pool})
	_, err = repository.CreateClient(ctx, mcpOAuth.CreateClientInput{
		ClientID:                clientID,
		Metadata:                json.RawMessage(`{"client_name":"MCP OAuth Test"}`),
		RedirectURIs:            []string{"https://example.com/callback"},
		GrantTypes:              []string{"authorization_code", "refresh_token"},
		ResponseTypes:           []string{"code"},
		TokenEndpointAuthMethod: "none",
		Scopes:                  []entity.Scope{entity.ScopeRead, entity.ScopeWrite},
	})
	if err != nil {
		t.Fatalf("create test client: %v", err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM mcp_oauth_clients WHERE client_id = $1`, clientID); err != nil {
			t.Errorf("delete test client: %v", err)
		}
	})

	return repository, integrationFixture{pool: pool, clientID: clientID, channelID: channel.ID, userID: user.ID}
}

func (fixture integrationFixture) createTokenInput(accessHash, refreshHash entity.CredentialHash) mcpOAuth.CreateTokenInput {
	return mcpOAuth.CreateTokenInput{
		ClientID:         fixture.clientID,
		ChannelID:        fixture.channelID,
		UserID:           fixture.userID,
		AccessTokenHash:  accessHash,
		RefreshTokenHash: refreshHash,
		Scopes:           []entity.Scope{entity.ScopeRead, entity.ScopeWrite},
		Resource:         "https://mcp.example.com",
		AccessExpiresAt:  time.Now().Add(time.Hour),
		RefreshExpiresAt: time.Now().Add(24 * time.Hour),
	}
}
