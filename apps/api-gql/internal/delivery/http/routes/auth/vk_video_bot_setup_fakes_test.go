package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"sync"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	entity "github.com/twirapp/twir/libs/entities/vk_video_bot"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	vkvideobotsrepo "github.com/twirapp/twir/libs/repositories/vk_video_bots"
)

type vkVideoBotSetupFixture struct {
	auth        *Auth
	sessions    *vkVideoBotSetupSessions
	provider    *vkVideoBotSetupProviderFake
	bots        *vkVideoBotRepositoryFake
	bindings    *vkVideoBotBindingRepositoryFake
	users       *vkVideoBotUsersRepositoryFake
	transaction *vkVideoBotTransactionFake
	kv          *vkVideoBotKVFake
	publisher   *oauthEventSubPublisher
}

func newVKVideoBotSetupFixture(admin usersmodel.User) *vkVideoBotSetupFixture {
	sessions := &vkVideoBotSetupSessions{userID: admin.ID}
	provider := &vkVideoBotSetupProviderFake{
		authorizationURL: "https://vk.example.test/authorize",
		tokens:           &appplatform.PlatformTokens{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 3600, Scopes: []string{"chat"}},
		user:             &appplatform.PlatformUser{ID: "vk-user", Login: "vk-login", DisplayName: "VK User"},
	}
	users := &vkVideoBotUsersRepositoryFake{users: map[uuid.UUID]usersmodel.User{admin.ID: admin}}
	bots := &vkVideoBotRepositoryFake{}
	bindings := &vkVideoBotBindingRepositoryFake{}
	transaction := &vkVideoBotTransactionFake{}
	kvStore := &vkVideoBotKVFake{values: make(map[string][]byte)}
	publisher := &oauthEventSubPublisher{}

	return &vkVideoBotSetupFixture{
		auth: &Auth{
			config:               testVKVideoBotConfig(),
			sessions:             sessions,
			usersRepo:            users,
			vkVideoBotProvider:   provider,
			vkVideoBotsRepo:      bots,
			channelPlatformsRepo: bindings,
			transactionRunner:    transaction,
			kv:                   kvStore,
			eventSubPublisher:    publisher,
		},
		sessions: sessions, provider: provider, bots: bots, bindings: bindings, users: users, transaction: transaction, kv: kvStore, publisher: publisher,
	}
}

type vkVideoBotSetupSessions struct {
	userID            uuid.UUID
	setIdentityCalls  int
	setPlatformCalls  int
	setDashboardCalls int
}

func (s *vkVideoBotSetupSessions) GetInternalUserID(context.Context) (uuid.UUID, error) {
	return s.userID, nil
}
func (s *vkVideoBotSetupSessions) SetSessionInternalUserID(context.Context, uuid.UUID) error {
	s.setIdentityCalls++
	return nil
}
func (s *vkVideoBotSetupSessions) SetSessionCurrentPlatform(context.Context, string) error {
	s.setPlatformCalls++
	return nil
}
func (s *vkVideoBotSetupSessions) SetSessionSelectedDashboard(context.Context, string) error {
	s.setDashboardCalls++
	return nil
}
func (*vkVideoBotSetupSessions) SetSessionTwitchUser(context.Context, helix.User) error { return nil }
func (*vkVideoBotSetupSessions) SetSessionKickUser(context.Context, authsessions.KickSessionUser) error {
	return nil
}
func (*vkVideoBotSetupSessions) SetOAuthAttempt(context.Context, string, authsessions.OAuthAttempt) error {
	return nil
}
func (*vkVideoBotSetupSessions) GetOAuthAttempt(context.Context, string) (authsessions.OAuthAttempt, error) {
	return authsessions.OAuthAttempt{}, errors.New("unexpected")
}
func (*vkVideoBotSetupSessions) DeleteOAuthAttempt(context.Context, string) error { return nil }

type vkVideoBotSetupProviderFake struct {
	authorizationURL string
	tokens           *appplatform.PlatformTokens
	user             *appplatform.PlatformUser
	exchangeCalls    int
}

func (p *vkVideoBotSetupProviderFake) GetBotSetupAuthURL(state string) (string, error) {
	parsed, err := url.Parse(p.authorizationURL)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("state", state)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
func (p *vkVideoBotSetupProviderFake) ExchangeBotSetupCode(context.Context, string) (*appplatform.PlatformTokens, error) {
	p.exchangeCalls++
	return p.tokens, nil
}
func (p *vkVideoBotSetupProviderFake) GetUser(context.Context, string) (*appplatform.PlatformUser, error) {
	return p.user, nil
}

type vkVideoBotRepositoryFake struct {
	bot         entity.VKVideoBot
	lockCalls   int
	upsertCalls int
}

func (r *vkVideoBotRepositoryFake) Get(context.Context) (entity.VKVideoBot, error) {
	if r.bot.IsNil() || r.bot.ID == uuid.Nil {
		return entity.Nil, vkvideobotsrepo.ErrNotFound
	}
	return r.bot, nil
}
func (r *vkVideoBotRepositoryFake) Lock(context.Context) error { r.lockCalls++; return nil }
func (r *vkVideoBotRepositoryFake) Upsert(_ context.Context, input vkvideobotsrepo.UpsertInput) (entity.VKVideoBot, error) {
	r.upsertCalls++
	if r.bot.ID == uuid.Nil {
		r.bot.ID = uuid.New()
	}
	r.bot.EncryptedAccessToken = input.EncryptedAccessToken
	r.bot.EncryptedRefreshToken = input.EncryptedRefreshToken
	r.bot.Scopes = input.Scopes
	r.bot.ExpiresIn = input.ExpiresIn
	r.bot.ObtainmentTimestamp = input.ObtainmentTimestamp
	r.bot.VKUserID = input.VKUserID
	return r.bot, nil
}
func (*vkVideoBotRepositoryFake) Update(context.Context, vkvideobotsrepo.UpdateInput) (entity.VKVideoBot, error) {
	return entity.Nil, errors.New("unexpected update")
}

type vkVideoBotBindingRepositoryFake struct {
	channelplatformsrepo.Repository
	assignedUserID   uuid.UUID
	assignCalls      int
	assignChannelIDs []uuid.UUID
}

func (r *vkVideoBotBindingRepositoryFake) AssignVKVideoLiveBot(_ context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	r.assignCalls++
	r.assignedUserID = userID
	if r.assignChannelIDs != nil {
		return r.assignChannelIDs, nil
	}
	return []uuid.UUID{uuid.New()}, nil
}

type vkVideoBotUsersRepositoryFake struct {
	users map[uuid.UUID]usersmodel.User
}

func (r *vkVideoBotUsersRepositoryFake) GetByID(_ context.Context, id uuid.UUID) (usersmodel.User, error) {
	user, ok := r.users[id]
	if !ok {
		return usersmodel.Nil, usersmodel.ErrNotFound
	}
	return user, nil
}
func (r *vkVideoBotUsersRepositoryFake) GetByPlatformID(_ context.Context, platform platformentity.Platform, platformID string) (usersmodel.User, error) {
	for _, user := range r.users {
		if user.Platform == platform && user.PlatformID == platformID {
			return user, nil
		}
	}
	return usersmodel.Nil, usersmodel.ErrNotFound
}
func (r *vkVideoBotUsersRepositoryFake) Create(_ context.Context, input usersrepo.CreateInput) (usersmodel.User, error) {
	user := usersmodel.User{ID: uuid.New(), Platform: input.Platform, PlatformID: input.PlatformID, Login: input.Login, DisplayName: input.DisplayName}
	r.users[user.ID] = user
	return user, nil
}
func (*vkVideoBotUsersRepositoryFake) Update(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error) {
	return usersmodel.User{}, nil
}
func (*vkVideoBotUsersRepositoryFake) GetManyByIDS(context.Context, usersrepo.GetManyInput) ([]usersmodel.User, error) {
	return nil, errors.New("unexpected")
}
func (*vkVideoBotUsersRepositoryFake) GetRandomOnlineUser(context.Context, usersrepo.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, errors.New("unexpected")
}
func (*vkVideoBotUsersRepositoryFake) GetOnlineUsersWithFilters(context.Context, usersrepo.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	return nil, errors.New("unexpected")
}
func (*vkVideoBotUsersRepositoryFake) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	return usersmodel.Nil, errors.New("unexpected")
}

type vkVideoBotTransactionFake struct{ calls int }

func (t *vkVideoBotTransactionFake) Do(ctx context.Context, fn func(context.Context) error) error {
	t.calls++
	return fn(ctx)
}

type vkVideoBotKVFake struct {
	mu     sync.Mutex
	values map[string][]byte
}

func (s *vkVideoBotKVFake) Get(_ context.Context, key string) kv.Valuer {
	s.mu.Lock()
	defer s.mu.Unlock()
	return vkVideoBotValue{bytes: append([]byte(nil), s.values[key]...), err: mapLookupError(s.values, key)}
}
func (s *vkVideoBotKVFake) Set(_ context.Context, key string, value any, _ ...kvoptions.Option) error {
	bytes, ok := value.([]byte)
	if !ok {
		var err error
		bytes, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = append([]byte(nil), bytes...)
	return nil
}
func (s *vkVideoBotKVFake) SetMany(context.Context, []kv.SetMany) error {
	return errors.New("unexpected")
}
func (s *vkVideoBotKVFake) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.values, key)
	return nil
}
func (s *vkVideoBotKVFake) DeleteMany(context.Context, []string) error {
	return errors.New("unexpected")
}
func (s *vkVideoBotKVFake) Exists(context.Context, string) (bool, error) {
	return false, errors.New("unexpected")
}
func (s *vkVideoBotKVFake) ExistsMany(context.Context, []string) ([]bool, error) {
	return nil, errors.New("unexpected")
}
func (s *vkVideoBotKVFake) GetKeysByPattern(context.Context, string) ([]string, error) {
	return nil, errors.New("unexpected")
}

type vkVideoBotValue struct {
	bytes []byte
	err   error
}

func (v vkVideoBotValue) Int() (int64, error)     { return 0, v.err }
func (v vkVideoBotValue) String() (string, error) { return string(v.bytes), v.err }
func (v vkVideoBotValue) Bytes() ([]byte, error)  { return v.bytes, v.err }
func (v vkVideoBotValue) Bool() (bool, error)     { return false, v.err }
func (v vkVideoBotValue) Float() (float64, error) { return 0, v.err }
func (v vkVideoBotValue) Scan(any) error          { return v.err }
func (v vkVideoBotValue) Err() error              { return v.err }
func mapLookupError(values map[string][]byte, key string) error {
	if _, ok := values[key]; !ok {
		return kv.ErrKeyNil
	}
	return nil
}
