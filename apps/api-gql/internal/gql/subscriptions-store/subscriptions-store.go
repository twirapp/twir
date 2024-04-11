package subscriptions_store

import (
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

func New() *SubscriptionsStore {
	return &SubscriptionsStore{
		NewNotificationsChannels: make(map[string]chan *gqlmodel.UserNotification),
	}
}

type SubscriptionsStore struct {
	// key is authenticated user id
	NewNotificationsChannels map[string]chan *gqlmodel.UserNotification
}
