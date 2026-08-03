package postgres

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	mcpOAuth "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

func TestPgx_ConsumeAuthorizationCode_allows_one_concurrent_consumer(t *testing.T) {
	// Given
	ctx := context.Background()
	repository, fixture := newIntegrationFixture(t, ctx)
	codeHash := testHash(1)
	_, err := repository.CreateAuthorizationCode(ctx, mcpOAuth.CreateAuthorizationCodeInput{
		CodeHash:      codeHash,
		ClientID:      fixture.clientID,
		ChannelID:     fixture.channelID,
		UserID:        fixture.userID,
		RedirectURI:   "https://example.com/callback",
		PKCEChallenge: "challenge",
		Scopes:        []entity.Scope{entity.ScopeRead},
		Resource:      "https://mcp.example.com",
		ExpiresAt:     time.Now().Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("create authorization code: %v", err)
	}

	// When
	errorsByConsumer := consumeAuthorizationCodeConcurrently(ctx, repository, codeHash)

	// Then
	var successfulConsumers int
	var missingCodeConsumers int
	for _, consumeErr := range errorsByConsumer {
		if consumeErr == nil {
			successfulConsumers++
			continue
		}
		if errors.Is(consumeErr, mcpOAuth.ErrAuthorizationCodeNotFound) {
			missingCodeConsumers++
			continue
		}
		t.Fatalf("consume authorization code: %v", consumeErr)
	}
	if successfulConsumers != 1 || missingCodeConsumers != 1 {
		t.Fatalf("successful consumers = %d, missing consumers = %d, want 1 each", successfulConsumers, missingCodeConsumers)
	}
}

func TestPgx_RotateRefreshToken_revokes_the_family_on_reuse(t *testing.T) {
	// Given
	ctx := context.Background()
	repository, fixture := newIntegrationFixture(t, ctx)
	presentedRefreshHash := testHash(2)
	initial, err := repository.CreateToken(ctx, fixture.createTokenInput(testHash(3), presentedRefreshHash))
	if err != nil {
		t.Fatalf("create initial token: %v", err)
	}

	// When
	errorsByConsumer := rotateRefreshTokenConcurrently(ctx, repository, presentedRefreshHash)

	// Then
	var successfulRotations int
	var reuseErrors int
	for _, rotateErr := range errorsByConsumer {
		if rotateErr == nil {
			successfulRotations++
			continue
		}
		var reuseError *mcpOAuth.RefreshTokenReuseError
		if errors.As(rotateErr, &reuseError) && reuseError.FamilyID == initial.FamilyID {
			reuseErrors++
			continue
		}
		t.Fatalf("rotate refresh token: %v", rotateErr)
	}
	if successfulRotations != 1 || reuseErrors != 1 {
		t.Fatalf("successful rotations = %d, reuse errors = %d, want 1 each", successfulRotations, reuseErrors)
	}

	var activeTokens int
	if err := fixture.pool.QueryRow(ctx, `SELECT COUNT(*) FROM mcp_oauth_tokens WHERE family_id = $1 AND revoked_at IS NULL`, initial.FamilyID).Scan(&activeTokens); err != nil {
		t.Fatalf("count active family tokens: %v", err)
	}
	if activeTokens != 0 {
		t.Fatalf("active family tokens = %d, want 0 after refresh reuse", activeTokens)
	}
}

func TestPgx_CreateToken_round_trips_granular_scopes(t *testing.T) {
	// Given
	ctx := context.Background()
	repository, fixture := newIntegrationFixture(t, ctx)
	accessHash := testHash(30)
	granularScopes := []entity.Scope{"commands:read", "timers:edit", "dashboard:read"}

	// When
	token, err := repository.CreateToken(ctx, mcpOAuth.CreateTokenInput{
		ClientID:         fixture.clientID,
		ChannelID:        fixture.channelID,
		UserID:           fixture.userID,
		AccessTokenHash:  accessHash,
		RefreshTokenHash: testHash(31),
		Scopes:           granularScopes,
		Resource:         "https://mcp.example.com",
		AccessExpiresAt:  time.Now().Add(time.Hour),
		RefreshExpiresAt: time.Now().Add(24 * time.Hour),
	})

	// Then
	if err != nil {
		t.Fatalf("create token: %v", err)
	}
	if !slices.Equal(token.Scopes, granularScopes) {
		t.Fatalf("created scopes = %v, want %v", token.Scopes, granularScopes)
	}
	fetched, err := repository.GetTokenByAccessTokenHash(ctx, accessHash)
	if err != nil {
		t.Fatalf("get token: %v", err)
	}
	if !slices.Equal(fetched.Scopes, granularScopes) {
		t.Fatalf("fetched scopes = %v, want %v", fetched.Scopes, granularScopes)
	}
}

func TestPgx_RevokeToken_revokes_family_and_ignores_unknown_hashes(t *testing.T) {
	// Given
	ctx := context.Background()
	repository, fixture := newIntegrationFixture(t, ctx)
	accessHash := testHash(8)
	token, err := repository.CreateToken(ctx, fixture.createTokenInput(accessHash, testHash(9)))
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	// When
	err = repository.RevokeToken(ctx, fixture.clientID, accessHash)
	if err != nil {
		t.Fatalf("revoke access token: %v", err)
	}
	if err := repository.RevokeToken(ctx, fixture.clientID, testHash(10)); err != nil {
		t.Fatalf("revoke unknown token: %v", err)
	}

	// Then
	assertFamilyRevoked(t, ctx, fixture, token.FamilyID)

	refreshToken, err := repository.CreateToken(ctx, fixture.createTokenInput(testHash(11), testHash(12)))
	if err != nil {
		t.Fatalf("create refresh-revocation token: %v", err)
	}
	if err := repository.RevokeToken(ctx, fixture.clientID, refreshToken.RefreshTokenHash); err != nil {
		t.Fatalf("revoke refresh token: %v", err)
	}
	assertFamilyRevoked(t, ctx, fixture, refreshToken.FamilyID)
}

func consumeAuthorizationCodeConcurrently(ctx context.Context, repository *Pgx, codeHash entity.CredentialHash) []error {
	return concurrently(2, func() error {
		_, err := repository.ConsumeAuthorizationCode(ctx, codeHash)
		return err
	})
}

func rotateRefreshTokenConcurrently(ctx context.Context, repository *Pgx, refreshHash entity.CredentialHash) []error {
	var hashCounter byte = 20
	var hashLock sync.Mutex
	return concurrently(2, func() error {
		hashLock.Lock()
		accessHash := testHash(hashCounter)
		hashCounter++
		nextRefreshHash := testHash(hashCounter)
		hashCounter++
		hashLock.Unlock()
		_, err := repository.RotateRefreshToken(ctx, mcpOAuth.RotateRefreshTokenInput{
			PresentedRefreshTokenHash: refreshHash,
			NextAccessTokenHash:       accessHash,
			NextRefreshTokenHash:      nextRefreshHash,
			AccessExpiresAt:           time.Now().Add(time.Hour),
			RefreshExpiresAt:          time.Now().Add(24 * time.Hour),
		})
		return err
	})
}

func concurrently(count int, action func() error) []error {
	start := make(chan struct{})
	results := make(chan error, count)
	var group sync.WaitGroup
	group.Add(count)
	for range count {
		go func() {
			defer group.Done()
			<-start
			results <- action()
		}()
	}
	close(start)
	group.Wait()
	close(results)

	errorsByAction := make([]error, 0, count)
	for err := range results {
		errorsByAction = append(errorsByAction, err)
	}
	return errorsByAction
}

func testHash(value byte) entity.CredentialHash {
	return entity.CredentialHash{0: value, 31: value}
}

func assertFamilyRevoked(t *testing.T, ctx context.Context, fixture integrationFixture, familyID uuid.UUID) {
	t.Helper()

	var activeTokens int
	if err := fixture.pool.QueryRow(ctx, `SELECT COUNT(*) FROM mcp_oauth_tokens WHERE family_id = $1 AND revoked_at IS NULL`, familyID).Scan(&activeTokens); err != nil {
		t.Fatalf("count revoked family tokens: %v", err)
	}
	if activeTokens != 0 {
		t.Fatalf("active family tokens = %d, want 0", activeTokens)
	}
}

func testClientID() string {
	return "mcp-oauth-test-" + uuid.NewString()
}
