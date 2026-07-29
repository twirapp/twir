package helpers

import (
	"errors"

	"github.com/olahol/melody"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/zap"
)

var ErrUserNotFound = errors.New("no user found")

// CheckChannelByApiKey resolves the internal channel ID by a channel or a user
// API key and stores it in the session under the "userId" key, which is then
// compared against channel IDs across the overlay namespaces.
func CheckChannelByApiKey(
	session *melody.Session,
	channelsRepo channels.Repository,
	usersRepo users.Repository,
) error {
	apiKey := session.Request.URL.Query().Get("apiKey")
	if apiKey == "" {
		session.Close()
		return errors.New("no api key")
	}

	ctx := session.Request.Context()

	channel, err := channelsRepo.GetByApiKey(ctx, apiKey)
	if err != nil && !errors.Is(err, channels.ErrNotFound) {
		session.Close()
		return err
	}

	if !channel.IsNil() {
		session.Set("userId", channel.ID.String())
		return nil
	}

	user, err := usersRepo.GetByApiKey(ctx, apiKey)
	if err != nil {
		zap.S().Errorf(apiKey, err)
		session.Close()
		return ErrUserNotFound
	}

	channel, err = channelsRepo.GetByBindingUserID(ctx, user.Platform, user.ID)
	if err != nil {
		zap.S().Errorf(apiKey, err)
		session.Close()
		return ErrUserNotFound
	}

	session.Set("userId", channel.ID.String())

	return nil
}
