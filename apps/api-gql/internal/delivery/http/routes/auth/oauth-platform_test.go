package auth

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	tokensrepo "github.com/twirapp/twir/libs/repositories/tokens"
	tokensmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestUpsertPlatformUserToken_BindsCreatedTokenToUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	createdTokenID := uuid.New()

	tokensRepository := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return nil, tokensrepo.ErrNotFound
		},
		createUserTokenFunc: func(_ context.Context, input tokensrepo.CreateInput) (*tokensmodel.Token, error) {
			if input.UserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, input.UserID)
			}

			return &tokensmodel.Token{ID: createdTokenID}, nil
		},
	}

	usersRepository := &fakeUsersRepository{
		updateFunc: func(_ context.Context, id uuid.UUID, input usersrepo.UpdateInput) (usersmodel.User, error) {
			if id != userID {
				t.Fatalf("expected update userID %s, got %s", userID, id)
			}

			if input.TokenID == nil {
				t.Fatal("expected token binding update")
			}

			if got, want := *input.TokenID, createdTokenID.String(); got != want {
				t.Fatalf("expected tokenID %s, got %s", want, got)
			}

			return usersmodel.User{}, nil
		},
	}

	authHandler := newTestAuth(tokensRepository, usersRepository)

	err := authHandler.upsertPlatformUserToken(context.Background(), userID, testPlatformTokens())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokensRepository.createUserTokenCalls != 1 {
		t.Fatalf("expected CreateUserToken to be called once, got %d", tokensRepository.createUserTokenCalls)
	}

	if usersRepository.updateCalls != 1 {
		t.Fatalf("expected users update to be called once, got %d", usersRepository.updateCalls)
	}

	if tokensRepository.updateTokenCalls != 0 {
		t.Fatalf("expected UpdateTokenByID not to be called, got %d", tokensRepository.updateTokenCalls)
	}
}

func TestUpsertPlatformUserToken_UpdatesExistingTokenWithoutRebinding(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	existingTokenID := uuid.New()

	tokensRepository := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return &tokensmodel.Token{ID: existingTokenID}, nil
		},
		updateTokenByIDFunc: func(_ context.Context, id uuid.UUID, input tokensrepo.UpdateTokenInput) (*tokensmodel.Token, error) {
			if id != existingTokenID {
				t.Fatalf("expected tokenID %s, got %s", existingTokenID, id)
			}

			if input.AccessToken == nil || input.RefreshToken == nil {
				t.Fatal("expected encrypted access and refresh tokens")
			}

			return &tokensmodel.Token{ID: existingTokenID}, nil
		},
	}

	usersRepository := &fakeUsersRepository{
		updateFunc: func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error) {
			t.Fatal("users update should not be called for existing token")
			return usersmodel.User{}, nil
		},
	}

	authHandler := newTestAuth(tokensRepository, usersRepository)

	err := authHandler.upsertPlatformUserToken(context.Background(), userID, testPlatformTokens())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokensRepository.updateTokenCalls != 1 {
		t.Fatalf("expected UpdateTokenByID to be called once, got %d", tokensRepository.updateTokenCalls)
	}

	if tokensRepository.createUserTokenCalls != 0 {
		t.Fatalf("expected CreateUserToken not to be called, got %d", tokensRepository.createUserTokenCalls)
	}

	if usersRepository.updateCalls != 0 {
		t.Fatalf("expected users update not to be called, got %d", usersRepository.updateCalls)
	}
}

func TestUpsertPlatformUserToken_ReturnsErrorWhenBindingFails(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	createdTokenID := uuid.New()
	bindingErr := errors.New("users update failed")

	tokensRepository := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return nil, tokensrepo.ErrNotFound
		},
		createUserTokenFunc: func(context.Context, tokensrepo.CreateInput) (*tokensmodel.Token, error) {
			return &tokensmodel.Token{ID: createdTokenID}, nil
		},
	}

	usersRepository := &fakeUsersRepository{
		updateFunc: func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error) {
			return usersmodel.User{}, bindingErr
		},
	}

	authHandler := newTestAuth(tokensRepository, usersRepository)

	err := authHandler.upsertPlatformUserToken(context.Background(), userID, testPlatformTokens())
	if !errors.Is(err, bindingErr) {
		t.Fatalf("expected binding error %v, got %v", bindingErr, err)
	}

	if usersRepository.updateCalls != 1 {
		t.Fatalf("expected users update to be called once, got %d", usersRepository.updateCalls)
	}
}

func newTestAuth(tokensRepository tokensrepo.Repository, usersRepository usersrepo.Repository) *Auth {
	return &Auth{
		config:           cfg.Config{TokensCipherKey: "pnyfwfiulmnqlhkvixaeligpprcnlyke"},
		logger:           slog.New(slog.NewTextHandler(testWriter{t: nil}, nil)),
		tokensRepository: tokensRepository,
		usersRepo:        usersRepository,
	}
}

func testPlatformTokens() *appplatform.PlatformTokens {
	return &appplatform.PlatformTokens{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    7200,
		Scopes:       []string{"chat:read", "chat:write"},
	}
}

type fakeTokensRepository struct {
	getByUserIDFunc      func(context.Context, uuid.UUID) (*tokensmodel.Token, error)
	createUserTokenFunc  func(context.Context, tokensrepo.CreateInput) (*tokensmodel.Token, error)
	updateTokenByIDFunc  func(context.Context, uuid.UUID, tokensrepo.UpdateTokenInput) (*tokensmodel.Token, error)
	createUserTokenCalls int
	updateTokenCalls     int
}

func (f *fakeTokensRepository) GetByID(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
	panic("unexpected GetByID call")
}

func (f *fakeTokensRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*tokensmodel.Token, error) {
	if f.getByUserIDFunc == nil {
		panic("unexpected GetByUserID call")
	}

	return f.getByUserIDFunc(ctx, userID)
}

func (f *fakeTokensRepository) GetByBotID(context.Context, string) (*tokensmodel.Token, error) {
	panic("unexpected GetByBotID call")
}

func (f *fakeTokensRepository) CreateUserToken(ctx context.Context, input tokensrepo.CreateInput) (*tokensmodel.Token, error) {
	f.createUserTokenCalls++
	if f.createUserTokenFunc == nil {
		panic("unexpected CreateUserToken call")
	}

	return f.createUserTokenFunc(ctx, input)
}

func (f *fakeTokensRepository) UpdateTokenByID(
	ctx context.Context,
	id uuid.UUID,
	input tokensrepo.UpdateTokenInput,
) (*tokensmodel.Token, error) {
	f.updateTokenCalls++
	if f.updateTokenByIDFunc == nil {
		panic("unexpected UpdateTokenByID call")
	}

	return f.updateTokenByIDFunc(ctx, id, input)
}

type fakeUsersRepository struct {
	updateFunc  func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error)
	updateCalls int
}

func (f *fakeUsersRepository) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	panic("unexpected GetByID call")
}

func (f *fakeUsersRepository) GetByPlatformID(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
	panic("unexpected GetByPlatformID call")
}

func (f *fakeUsersRepository) GetManyByIDS(context.Context, usersrepo.GetManyInput) ([]usersmodel.User, error) {
	panic("unexpected GetManyByIDS call")
}

func (f *fakeUsersRepository) Update(ctx context.Context, id uuid.UUID, input usersrepo.UpdateInput) (usersmodel.User, error) {
	f.updateCalls++
	if f.updateFunc == nil {
		panic("unexpected Update call")
	}

	return f.updateFunc(ctx, id, input)
}

func (f *fakeUsersRepository) GetRandomOnlineUser(context.Context, usersrepo.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	panic("unexpected GetRandomOnlineUser call")
}

func (f *fakeUsersRepository) GetOnlineUsersWithFilters(context.Context, usersrepo.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	panic("unexpected GetOnlineUsersWithFilters call")
}

func (f *fakeUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	panic("unexpected GetByApiKey call")
}

func (f *fakeUsersRepository) Create(context.Context, usersrepo.CreateInput) (usersmodel.User, error) {
	panic("unexpected Create call")
}

type testWriter struct {
	t *testing.T
}

func (w testWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var _ = time.Now
