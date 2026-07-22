package tts

import (
	"context"
	"testing"

	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestResolveChannelIDByAPIKeyUsesUserPlatformBinding(t *testing.T) {
	for _, platform := range []platformentity.Platform{
		platformentity.PlatformTwitch,
		platformentity.PlatformKick,
		platformentity.PlatformVKVideoLive,
	} {
		t.Run(platform.String(), func(t *testing.T) {
			userID := uuid.New()
			channelID := uuid.New()
			lookup := &ttsTestChannelLookup{channel: channelsmodel.Channel{ID: channelID}}
			service := &Service{
				usersRepository: ttsTestUsersRepository{user: usersmodel.User{
					ID:       userID,
					Platform: platform,
				}},
				channelService: lookup,
			}

			got, err := service.ResolveChannelIDByAPIKey(context.Background(), "api-key")
			if err != nil {
				t.Fatalf("resolve channel ID by API key: %v", err)
			}
			if got != channelID.String() {
				t.Fatalf("channel ID = %q, want %q", got, channelID)
			}
			if lookup.platform != platform || lookup.userID != userID {
				t.Fatalf(
					"binding lookup = (%s, %s), want (%s, %s)",
					lookup.platform,
					lookup.userID,
					platform,
					userID,
				)
			}
		})
	}
}

type ttsTestUsersRepository struct {
	user usersmodel.User
}

func (r ttsTestUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	return r.user, nil
}

type ttsTestChannelLookup struct {
	channel  channelsmodel.Channel
	platform platformentity.Platform
	userID   uuid.UUID
}

func (r *ttsTestChannelLookup) GetChannelByBindingUserID(
	_ context.Context,
	p platformentity.Platform,
	userID uuid.UUID,
) (channelsmodel.Channel, error) {
	r.platform = p
	r.userID = userID
	return r.channel, nil
}
