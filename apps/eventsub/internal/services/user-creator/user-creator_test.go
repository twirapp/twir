package user_creator

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/users"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestEnsureUserExistsEnrichesExistingUserChatIdentity(t *testing.T) {
	t.Parallel()

	repository := &usersRepositoryStub{
		user: usermodel.User{
			ID:         uuid.New(),
			Platform:   platform.PlatformTwitch,
			PlatformID: "123",
		},
	}
	service := &UserCreatorService{usersRepo: repository}
	input := CreateUserInput{
		UserID:      "123",
		PlatformID:  "123",
		Platform:    platform.PlatformTwitch,
		Login:       "alice",
		DisplayName: "Alice",
	}

	user, err := service.ensureUserExists(context.Background(), input)
	if err != nil {
		t.Fatalf("ensure user exists: %v", err)
	}

	if repository.updateCalls != 1 {
		t.Fatalf("expected one identity update, got %d", repository.updateCalls)
	}
	if user.Login != "alice" {
		t.Fatalf("expected login %q, got %q", "alice", user.Login)
	}
	if user.DisplayName != "Alice" {
		t.Fatalf("expected display name %q, got %q", "Alice", user.DisplayName)
	}
}

type usersRepositoryStub struct {
	user        usermodel.User
	updateCalls int
}

func (s *usersRepositoryStub) GetByID(context.Context, uuid.UUID) (usermodel.User, error) {
	return usermodel.Nil, usermodel.ErrNotFound
}

func (s *usersRepositoryStub) GetByPlatformID(
	context.Context,
	platform.Platform,
	string,
) (usermodel.User, error) {
	return s.user, nil
}

func (s *usersRepositoryStub) GetManyByIDS(context.Context, users.GetManyInput) ([]usermodel.User, error) {
	return nil, errors.New("unexpected GetManyByIDS call")
}

func (s *usersRepositoryStub) Update(
	_ context.Context,
	_ uuid.UUID,
	input users.UpdateInput,
) (usermodel.User, error) {
	s.updateCalls++
	if input.Login != nil {
		s.user.Login = *input.Login
	}
	if input.DisplayName != nil {
		s.user.DisplayName = *input.DisplayName
	}
	return s.user, nil
}

func (s *usersRepositoryStub) GetRandomOnlineUser(
	context.Context,
	users.GetRandomOnlineUserInput,
) (usermodel.OnlineUser, error) {
	return usermodel.NilOnlineUser, errors.New("unexpected GetRandomOnlineUser call")
}

func (s *usersRepositoryStub) GetOnlineUsersWithFilters(
	context.Context,
	users.GetOnlineUsersWithFiltersInput,
) ([]usermodel.OnlineUser, error) {
	return nil, errors.New("unexpected GetOnlineUsersWithFilters call")
}

func (s *usersRepositoryStub) GetByApiKey(context.Context, string) (usermodel.User, error) {
	return usermodel.Nil, errors.New("unexpected GetByApiKey call")
}

func (s *usersRepositoryStub) Create(context.Context, users.CreateInput) (usermodel.User, error) {
	return usermodel.Nil, errors.New("unexpected Create call")
}
