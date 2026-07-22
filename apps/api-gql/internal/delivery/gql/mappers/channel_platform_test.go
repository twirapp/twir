package mappers

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestChannelPlatformBindingToGraphQLMapsProfileEnabledStateAndCapabilities(t *testing.T) {
	t.Parallel()

	bindingID := uuid.New()
	userID := uuid.New()
	got, err := ChannelPlatformBindingToGraphQL(
		channelplatformsmodel.ChannelPlatform{
			ID:                bindingID,
			Platform:          platformentity.PlatformTwitch,
			UserID:            userID,
			PlatformChannelID: "twitch-channel",
			Enabled:           true,
		},
		usersmodel.User{
			ID:          userID,
			PlatformID:  "twitch-user",
			Login:       "twitch-login",
			DisplayName: "Twitch Name",
			Avatar:      "https://example.com/twitch.png",
		},
		platformentity.PlatformTwitch.Capabilities(),
	)
	if err != nil {
		t.Fatalf("ChannelPlatformBindingToGraphQL() error = %v", err)
	}

	if got.ID != bindingID || got.Platform != gqlmodel.PlatformTwitch || got.UserID != userID || got.PlatformChannelID != "twitch-channel" || !got.Enabled {
		t.Fatalf("binding fields = %#v", got)
	}
	if got.PlatformUserID != "twitch-user" || got.PlatformLogin != "twitch-login" || got.PlatformDisplayName != "Twitch Name" || got.PlatformAvatar == nil || *got.PlatformAvatar != "https://example.com/twitch.png" {
		t.Fatalf("profile fields = %#v", got)
	}

	gotCapabilities := make([]string, 0, len(got.Capabilities))
	for _, capability := range got.Capabilities {
		gotCapabilities = append(gotCapabilities, capability.Name)
	}
	wantCapabilities := []string{
		"chat.read",
		"chat.write",
		"chat.reply",
		"moderation.delete",
		"streams.read",
		"events.follow",
		"events.raid",
		"events.reward",
	}
	if !reflect.DeepEqual(gotCapabilities, wantCapabilities) {
		t.Fatalf("capability strings = %#v, want %#v", gotCapabilities, wantCapabilities)
	}
}

func TestPlatformMappersSupportVKVideoLive(t *testing.T) {
	t.Parallel()

	entityPlatform, err := GraphQLPlatformToEntity(gqlmodel.PlatformVkVideoLive)
	if err != nil {
		t.Fatalf("GraphQLPlatformToEntity() error = %v", err)
	}
	if entityPlatform != platformentity.PlatformVKVideoLive {
		t.Fatalf("GraphQLPlatformToEntity() = %q, want %q", entityPlatform, platformentity.PlatformVKVideoLive)
	}

	graphqlPlatform, err := EntityPlatformToGraphQL(platformentity.PlatformVKVideoLive)
	if err != nil {
		t.Fatalf("EntityPlatformToGraphQL() error = %v", err)
	}
	if graphqlPlatform != gqlmodel.PlatformVkVideoLive {
		t.Fatalf("EntityPlatformToGraphQL() = %q, want %q", graphqlPlatform, gqlmodel.PlatformVkVideoLive)
	}
}

func TestMapChannelModelToGqlPublicUserMapsProfilesFromBindings(t *testing.T) {
	t.Parallel()

	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	channel := channelsmodel.Channel{
		ID: uuid.New(),
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{Platform: platformentity.PlatformTwitch, UserID: twitchUserID},
			{Platform: platformentity.PlatformKick, UserID: kickUserID},
		},
	}

	got := MapChannelModelToGqlPublicUser(channel, map[uuid.UUID]usersmodel.User{
		twitchUserID: {
			ID: twitchUserID, PlatformID: "twitch-user", Login: "twitch-login", DisplayName: "Twitch Name", Avatar: "https://example.com/twitch.png",
		},
		kickUserID: {
			ID: kickUserID, PlatformID: "kick-user", Login: "kick-login", DisplayName: "Kick Name", Avatar: "https://example.com/kick.png",
		},
	})

	if got.ID != channel.ID || got.TwitchProfile == nil || got.KickProfile == nil {
		t.Fatalf("public user = %#v, want Twitch and Kick profiles", got)
	}
	if got.TwitchProfile.ID != "twitch-user" || got.TwitchProfile.DisplayName != "Twitch Name" {
		t.Fatalf("Twitch profile = %#v", got.TwitchProfile)
	}
	if got.KickProfile.ID != "kick-user" || got.KickProfile.DisplayName != "Kick Name" || got.KickProfile.ProfilePicture == nil || *got.KickProfile.ProfilePicture != "https://example.com/kick.png" {
		t.Fatalf("Kick profile = %#v", got.KickProfile)
	}
}
