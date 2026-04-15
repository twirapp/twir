package kick

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type stubUsersRepo struct {
	byID map[string]usersmodel.User
}

func (s *stubUsersRepo) GetByID(ctx context.Context, id string) (usersmodel.User, error) {
	if u, ok := s.byID[id]; ok {
		return u, nil
	}
	return usersmodel.Nil, usersmodel.ErrNotFound
}

func (s *stubUsersRepo) GetByPlatformID(_ context.Context, _ platform.Platform, _ string) (usersmodel.User, error) {
	return usersmodel.Nil, usersmodel.ErrNotFound
}
func (s *stubUsersRepo) GetManyByIDS(_ context.Context, _ usersrepository.GetManyInput) ([]usersmodel.User, error) {
	return nil, nil
}
func (s *stubUsersRepo) Update(_ context.Context, _ string, _ usersrepository.UpdateInput) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}
func (s *stubUsersRepo) GetRandomOnlineUser(_ context.Context, _ usersrepository.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, nil
}
func (s *stubUsersRepo) GetOnlineUsersWithFilters(_ context.Context, _ usersrepository.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	return nil, nil
}
func (s *stubUsersRepo) GetByApiKey(_ context.Context, _ string) (usersmodel.User, error) {
	return usersmodel.Nil, usersmodel.ErrNotFound
}
func (s *stubUsersRepo) Create(_ context.Context, _ usersrepository.CreateInput) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func newTestManager(httpClient *http.Client, apiBaseURL string, usersRepo usersrepository.Repository) (*SubscriptionManager, redismock.ClientMock) {
	db, mock := redismock.NewClientMock()
	m := &SubscriptionManager{
		config:     cfg.Config{SiteBaseUrl: "http://localhost:3005"},
		redis:      db,
		httpClient: httpClient,
		logger:     slog.Default(),
		apiBaseURL: apiBaseURL,
		usersRepo:  usersRepo,
	}
	return m, mock
}

func TestSubscribeAll_SendsFourPostsAndStoresInRedis(t *testing.T) {
	callCount := 0
	subIDs := map[string]string{
		"chat.message.sent": "sub-id-1",
		"channel.follow":    "sub-id-2",
		"stream.online":     "sub-id-3",
		"stream.offline":    "sub-id-4",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "wrong method", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var req subscribeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}

		callCount++
		subID, ok := subIDs[req.Type]
		if !ok {
			http.Error(w, "unknown event type", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(subscribeResponse{
			Data: subscriptionData{
				ID:                subID,
				BroadcasterUserID: req.BroadcasterUserID,
				Type:              req.Type,
				Method:            req.Method,
				CallbackURL:       req.CallbackURL,
			},
		})
	}))
	defer server.Close()

	m, mock := newTestManager(server.Client(), server.URL, &stubUsersRepo{
		byID: map[string]usersmodel.User{
			"f47ac10b-58cc-4372-a567-0e02b2c3d479": {ID: "f47ac10b-58cc-4372-a567-0e02b2c3d479", PlatformID: "12345"},
		},
	})

	for _, eventType := range EventTypes {
		key := redisKey("f47ac10b-58cc-4372-a567-0e02b2c3d479", eventType)
		mock.ExpectSet(key, subIDs[eventType], 25*time.Hour).SetVal("OK")
	}

	ctx := context.Background()
	err := m.SubscribeAll(ctx, "f47ac10b-58cc-4372-a567-0e02b2c3d479", "test-token")
	if err != nil {
		t.Fatalf("SubscribeAll returned unexpected error: %v", err)
	}

	if callCount != 4 {
		t.Errorf("expected 4 POST requests, got %d", callCount)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestSubscribeAll_Returns401Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized"}`))
	}))
	defer server.Close()

	m, _ := newTestManager(server.Client(), server.URL, &stubUsersRepo{
		byID: map[string]usersmodel.User{
			"f47ac10b-58cc-4372-a567-0e02b2c3d479": {ID: "f47ac10b-58cc-4372-a567-0e02b2c3d479", PlatformID: "12345"},
		},
	})

	ctx := context.Background()
	err := m.SubscribeAll(ctx, "f47ac10b-58cc-4372-a567-0e02b2c3d479", "bad-token")
	if err == nil {
		t.Fatal("expected error from SubscribeAll with 401 response, got nil")
	}
}
