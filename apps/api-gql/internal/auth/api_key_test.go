package auth

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

type apiKeyContextKey struct{}

func TestGetAuthenticatedUserByApiKeyResolvesChannelOwnerFromBindings(t *testing.T) {
	tests := []struct {
		name     string
		platform platform.Platform
	}{
		{name: "Twitch", platform: platform.PlatformTwitch},
		{name: "Kick", platform: platform.PlatformKick},
		{name: "VK Video Live", platform: platform.PlatformVKVideoLive},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ownerID := uuid.New()
			owner := usersmodel.User{
				ID:         ownerID,
				Platform:   tt.platform,
				PlatformID: tt.platform.String() + "-owner",
				ApiKey:     "owner-api-key",
			}
			channels := &apiKeyChannelsRepository{
				channel: channelsmodel.Channel{
					Bindings: []channelplatformsmodel.ChannelPlatform{{
						Platform: tt.platform,
						UserID:   ownerID,
					}},
				},
			}
			users := &apiKeyUsersRepository{usersByID: map[uuid.UUID]usersmodel.User{ownerID: owner}}
			auth := newAPIKeyAuthForTest(channels, users)
			ctx := authenticatedAPIKeyContext("channel-api-key")

			got, err := auth.GetAuthenticatedUserByApiKey(ctx)
			if err != nil {
				t.Fatalf("GetAuthenticatedUserByApiKey() error = %v", err)
			}
			if got.ID != ownerID.String() {
				t.Fatalf("owner ID = %q, want %q", got.ID, ownerID)
			}
			if got.PlatformID != owner.PlatformID {
				t.Fatalf("owner platform ID = %q, want %q", got.PlatformID, owner.PlatformID)
			}
			if got.ApiKey != owner.ApiKey {
				t.Fatalf("owner API key = %q, want %q", got.ApiKey, owner.ApiKey)
			}
			if len(channels.apiKeyLookups) != 1 || channels.apiKeyLookups[0] != "channel-api-key" {
				t.Fatalf("channel API key lookups = %#v, want [channel-api-key]", channels.apiKeyLookups)
			}
			if len(users.idLookups) != 1 || users.idLookups[0] != ownerID {
				t.Fatalf("owner ID lookups = %#v, want [%s]", users.idLookups, ownerID)
			}
			if channels.contexts[0].Value(apiKeyContextKey{}) != "request-context" {
				t.Fatal("channel lookup did not receive the request context")
			}
			if users.contexts[0].Value(apiKeyContextKey{}) != "request-context" {
				t.Fatal("owner lookup did not receive the request context")
			}
		})
	}
}

func TestGetAuthenticatedUserByApiKeyPreservesBindingOwnershipPriority(t *testing.T) {
	twitchOwnerID := uuid.New()
	kickOwnerID := uuid.New()
	vkOwnerID := uuid.New()
	channels := &apiKeyChannelsRepository{
		channel: channelsmodel.Channel{
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{Platform: platform.PlatformKick, UserID: kickOwnerID},
				{Platform: platform.PlatformVKVideoLive, UserID: vkOwnerID},
				{Platform: platform.PlatformTwitch, UserID: twitchOwnerID},
			},
		},
	}
	users := &apiKeyUsersRepository{usersByID: map[uuid.UUID]usersmodel.User{
		twitchOwnerID: {ID: twitchOwnerID, Platform: platform.PlatformTwitch},
		kickOwnerID:   {ID: kickOwnerID, Platform: platform.PlatformKick},
		vkOwnerID:     {ID: vkOwnerID, Platform: platform.PlatformVKVideoLive},
	}}

	got, err := newAPIKeyAuthForTest(channels, users).GetAuthenticatedUserByApiKey(
		authenticatedAPIKeyContext("channel-api-key"),
	)
	if err != nil {
		t.Fatalf("GetAuthenticatedUserByApiKey() error = %v", err)
	}
	if got.ID != twitchOwnerID.String() {
		t.Fatalf("owner ID = %q, want Twitch owner %q", got.ID, twitchOwnerID)
	}
	if len(users.idLookups) != 1 || users.idLookups[0] != twitchOwnerID {
		t.Fatalf("owner ID lookups = %#v, want Twitch owner %s", users.idLookups, twitchOwnerID)
	}
}

func TestGetAuthenticatedUserByApiKeyFallsBackToUserKeyWhenChannelOwnershipCannotResolve(t *testing.T) {
	tests := []struct {
		name     string
		channels *apiKeyChannelsRepository
		users    *apiKeyUsersRepository
	}{
		{
			name:     "channel has no bindings",
			channels: &apiKeyChannelsRepository{channel: channelsmodel.Channel{}},
			users: &apiKeyUsersRepository{userByAPIKey: usersmodel.User{
				ID: uuid.New(),
			}},
		},
		{
			name:     "channel lookup fails",
			channels: &apiKeyChannelsRepository{getByAPIKeyErr: errors.New("channel lookup failed")},
			users: &apiKeyUsersRepository{userByAPIKey: usersmodel.User{
				ID: uuid.New(),
			}},
		},
		{
			name: "binding owner lookup fails",
			channels: &apiKeyChannelsRepository{channel: channelsmodel.Channel{Bindings: []channelplatformsmodel.ChannelPlatform{{
				Platform: platform.PlatformTwitch,
				UserID:   uuid.New(),
			}}}},
			users: &apiKeyUsersRepository{
				getByIDErr: errors.New("owner lookup failed"),
				userByAPIKey: usersmodel.User{
					ID: uuid.New(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newAPIKeyAuthForTest(tt.channels, tt.users).GetAuthenticatedUserByApiKey(
				authenticatedAPIKeyContext("user-api-key"),
			)
			if err != nil {
				t.Fatalf("GetAuthenticatedUserByApiKey() error = %v", err)
			}
			if got.ID != tt.users.userByAPIKey.ID.String() {
				t.Fatalf("fallback user ID = %q, want %q", got.ID, tt.users.userByAPIKey.ID)
			}
			if len(tt.users.apiKeyLookups) != 1 || tt.users.apiKeyLookups[0] != "user-api-key" {
				t.Fatalf("user API key lookups = %#v, want [user-api-key]", tt.users.apiKeyLookups)
			}
		})
	}
}

func TestGetAuthenticatedUserByApiKeyFallsBackWhenBindingOwnerIsNilWithoutError(t *testing.T) {
	bindingOwnerID := uuid.New()
	fallbackUserID := uuid.New()
	channels := &apiKeyChannelsRepository{
		channel: channelsmodel.Channel{Bindings: []channelplatformsmodel.ChannelPlatform{{
			Platform: platform.PlatformTwitch,
			UserID:   bindingOwnerID,
		}}},
	}
	users := &apiKeyUsersRepository{
		usersByID: map[uuid.UUID]usersmodel.User{
			bindingOwnerID: usersmodel.Nil,
		},
		userByAPIKey: usersmodel.User{ID: fallbackUserID},
	}

	got, err := newAPIKeyAuthForTest(channels, users).GetAuthenticatedUserByApiKey(
		authenticatedAPIKeyContext("channel-api-key"),
	)
	if err != nil {
		t.Fatalf("GetAuthenticatedUserByApiKey() error = %v", err)
	}
	if got.ID == uuid.Nil.String() {
		t.Fatal("GetAuthenticatedUserByApiKey() authenticated a zero-ID binding owner")
	}
	if got.ID != fallbackUserID.String() {
		t.Fatalf("authenticated user ID = %q, want fallback user %q", got.ID, fallbackUserID)
	}
	if len(users.idLookups) != 1 || users.idLookups[0] != bindingOwnerID {
		t.Fatalf("binding owner lookups = %#v, want [%s]", users.idLookups, bindingOwnerID)
	}
	if len(users.apiKeyLookups) != 1 || users.apiKeyLookups[0] != "channel-api-key" {
		t.Fatalf("user API key lookups = %#v, want [channel-api-key]", users.apiKeyLookups)
	}
}

func TestGetAuthenticatedUserByApiKeyWrapsUserLookupErrors(t *testing.T) {
	userLookupErr := errors.New("user lookup failed")
	auth := newAPIKeyAuthForTest(
		&apiKeyChannelsRepository{},
		&apiKeyUsersRepository{getByAPIKeyErr: userLookupErr},
	)

	_, err := auth.GetAuthenticatedUserByApiKey(authenticatedAPIKeyContext("user-api-key"))
	if !errors.Is(err, userLookupErr) {
		t.Fatalf("GetAuthenticatedUserByApiKey() error = %v, want wrapped %v", err, userLookupErr)
	}
	if !strings.Contains(err.Error(), "cannot get user from db") {
		t.Fatalf("GetAuthenticatedUserByApiKey() error = %q, want user database context", err)
	}
}

func authenticatedAPIKeyContext(apiKey string) context.Context {
	ctx := context.WithValue(context.Background(), apiKeyContextKey{}, "request-context")
	return context.WithValue(ctx, WsApiKeyContextKey{}, apiKey)
}

func newAPIKeyAuthForTest(
	channels channelsrepo.Repository,
	users usersrepo.Repository,
) *Auth {
	return &Auth{
		channelService: channelservice.NewChannelService(
			channels,
			&buscore.Bus{},
			config.Config{},
			nil,
			nil,
		),
		usersRepo: users,
	}
}

type apiKeyChannelsRepository struct {
	channel        channelsmodel.Channel
	getByAPIKeyErr error
	apiKeyLookups  []string
	contexts       []context.Context
}

func (r *apiKeyChannelsRepository) GetByApiKey(
	ctx context.Context,
	apiKey string,
) (channelsmodel.Channel, error) {
	r.contexts = append(r.contexts, ctx)
	r.apiKeyLookups = append(r.apiKeyLookups, apiKey)
	if r.getByAPIKeyErr != nil {
		return channelsmodel.Nil, r.getByAPIKeyErr
	}

	return r.channel, nil
}

func (*apiKeyChannelsRepository) GetMany(context.Context, channelsrepo.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (*apiKeyChannelsRepository) GetAllByBindingPlatform(context.Context, platform.Platform) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (*apiKeyChannelsRepository) GetByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (*apiKeyChannelsRepository) GetByBindingUserID(
	context.Context,
	platform.Platform,
	uuid.UUID,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (*apiKeyChannelsRepository) GetByPlatformChannelID(
	context.Context,
	platform.Platform,
	string,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (*apiKeyChannelsRepository) GetBySlug(context.Context, channelsrepo.GetBySlugInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (*apiKeyChannelsRepository) GetCount(context.Context, channelsrepo.GetCountInput) (int, error) {
	return 0, nil
}

func (*apiKeyChannelsRepository) Update(
	context.Context,
	uuid.UUID,
	channelsrepo.UpdateInput,
) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (*apiKeyChannelsRepository) Create(context.Context, channelsrepo.CreateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

type apiKeyUsersRepository struct {
	usersByID      map[uuid.UUID]usersmodel.User
	userByAPIKey   usersmodel.User
	getByIDErr     error
	getByAPIKeyErr error
	idLookups      []uuid.UUID
	apiKeyLookups  []string
	contexts       []context.Context
}

func (r *apiKeyUsersRepository) GetByID(ctx context.Context, id uuid.UUID) (usersmodel.User, error) {
	r.contexts = append(r.contexts, ctx)
	r.idLookups = append(r.idLookups, id)
	if r.getByIDErr != nil {
		return usersmodel.Nil, r.getByIDErr
	}

	user, ok := r.usersByID[id]
	if !ok {
		return usersmodel.Nil, errors.New("user not found")
	}

	return user, nil
}

func (r *apiKeyUsersRepository) GetByApiKey(ctx context.Context, apiKey string) (usersmodel.User, error) {
	r.contexts = append(r.contexts, ctx)
	r.apiKeyLookups = append(r.apiKeyLookups, apiKey)
	if r.getByAPIKeyErr != nil {
		return usersmodel.Nil, r.getByAPIKeyErr
	}

	return r.userByAPIKey, nil
}

func (*apiKeyUsersRepository) GetByPlatformID(
	context.Context,
	platform.Platform,
	string,
) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func (*apiKeyUsersRepository) GetManyByIDS(
	context.Context,
	usersrepo.GetManyInput,
) ([]usersmodel.User, error) {
	return nil, nil
}

func (*apiKeyUsersRepository) Update(
	context.Context,
	uuid.UUID,
	usersrepo.UpdateInput,
) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func (*apiKeyUsersRepository) GetRandomOnlineUser(
	context.Context,
	usersrepo.GetRandomOnlineUserInput,
) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, nil
}

func (*apiKeyUsersRepository) GetOnlineUsersWithFilters(
	context.Context,
	usersrepo.GetOnlineUsersWithFiltersInput,
) ([]usersmodel.OnlineUser, error) {
	return nil, nil
}

func (*apiKeyUsersRepository) Create(context.Context, usersrepo.CreateInput) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}
