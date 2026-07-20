package twir_users

import (
	"context"
	"testing"

	"github.com/twirapp/twir/libs/repositories/users_with_channel"
	"github.com/twirapp/twir/libs/repositories/users_with_channel/model"
)

type usersWithChannelsRepositoryStub struct {
	getManyCalls int
	countCalls   int
}

func (s *usersWithChannelsRepositoryStub) GetByID(context.Context, string) (model.UserWithChannel, error) {
	return model.UserWithChannel{}, nil
}

func (s *usersWithChannelsRepositoryStub) GetManyByIDS(
	context.Context,
	users_with_channel.GetManyInput,
) ([]model.UserWithChannel, error) {
	s.getManyCalls++
	return nil, nil
}

func (s *usersWithChannelsRepositoryStub) GetManyCount(context.Context, users_with_channel.GetManyInput) (int, error) {
	s.countCalls++
	return 0, nil
}

func TestGetManyRejectsInvalidBadgeID(t *testing.T) {
	repository := &usersWithChannelsRepositoryStub{}
	service := &Service{usersWithChannelsRepository: repository}

	_, err := service.GetMany(context.Background(), GetManyInput{
		HasBadges: []string{"not-a-uuid"},
	})
	if err == nil {
		t.Fatal("expected invalid badge ID to be rejected")
	}
	if repository.getManyCalls != 0 || repository.countCalls != 0 {
		t.Fatalf(
			"expected no repository calls, got getMany=%d count=%d",
			repository.getManyCalls,
			repository.countCalls,
		)
	}
}
