package badges_users

import (
	"context"
	"testing"

	"github.com/google/uuid"
	repository "github.com/twirapp/twir/libs/repositories/badges_users"
	"github.com/twirapp/twir/libs/repositories/badges_users/model"
)

type repositoryStub struct {
	createCalls int
	deleteCalls int
}

func (s *repositoryStub) GetMany(context.Context, repository.GetManyInput) ([]model.BadgeUser, error) {
	return nil, nil
}

func (s *repositoryStub) Create(context.Context, repository.CreateInput) (model.BadgeUser, error) {
	s.createCalls++
	return model.BadgeUser{}, nil
}

func (s *repositoryStub) Delete(context.Context, repository.DeleteInput) error {
	s.deleteCalls++
	return nil
}

func TestCreateRejectsInvalidUserID(t *testing.T) {
	repository := &repositoryStub{}
	service := &Service{badgesUsersRepository: repository}

	_, err := service.Create(context.Background(), CreateInput{
		BadgeID: uuid.New(),
		UserID:  "not-a-uuid",
	})
	if err == nil {
		t.Fatal("expected invalid user ID to be rejected")
	}
	if repository.createCalls != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.createCalls)
	}
}

func TestDeleteRejectsInvalidUserID(t *testing.T) {
	repository := &repositoryStub{}
	service := &Service{badgesUsersRepository: repository}

	err := service.Delete(context.Background(), DeleteInput{
		BadgeID: uuid.New(),
		UserID:  "not-a-uuid",
	})
	if err == nil {
		t.Fatal("expected invalid user ID to be rejected")
	}
	if repository.deleteCalls != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.deleteCalls)
	}
}
