package state

import (
	"errors"
	"maps"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/twitch-mock/internal/config"
)

var ErrConduitNotFound = errors.New("conduit not found")

type ConduitShardTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id,omitempty"`
}

type ConduitShard struct {
	ID        int                   `json:"id"`
	Transport ConduitShardTransport `json:"transport"`
}

type Conduit struct {
	ID         string         `json:"id"`
	ShardCount int            `json:"shard_count"`
	Shards     []ConduitShard `json:"shards,omitempty"`
}

type SubscriptionTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id,omitempty"`
}

type Subscription struct {
	ID        string                `json:"id"`
	Status    string                `json:"status"`
	Type      string                `json:"type"`
	Version   string                `json:"version"`
	Cost      int                   `json:"cost"`
	Condition map[string]any        `json:"condition"`
	CreatedAt time.Time             `json:"created_at"`
	Transport SubscriptionTransport `json:"transport"`
}

type State struct {
	mu                sync.RWMutex
	conduits          map[string]Conduit
	subscriptions     map[string]Subscription
	streamOnline      bool
	streamStarted     time.Time
	streamTitle       string
	streamGameID      string
	streamGameName    string
	streamViewerCount int
	followersTotal    int
	followedUserIDs   map[string]time.Time
	moderatorIDs      map[string]bool
}

func New() *State {
	return &State{
		conduits:      map[string]Conduit{},
		subscriptions: map[string]Subscription{},
		streamTitle:   "Mock Stream",
		followedUserIDs: map[string]time.Time{
			config.MockBotID: time.Now().UTC(),
		},
		moderatorIDs: map[string]bool{
			config.MockBotID: true,
		},
	}
}

func (s *State) ListConduits() []Conduit {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Conduit, 0, len(s.conduits))
	for _, conduit := range s.conduits {
		result = append(result, cloneConduit(conduit))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func (s *State) CreateConduit(shardCount int) Conduit {
	if shardCount <= 0 {
		shardCount = 1
	}

	conduit := Conduit{
		ID:         uuid.NewString(),
		ShardCount: shardCount,
		Shards:     make([]ConduitShard, 0, shardCount),
	}

	for i := range shardCount {
		conduit.Shards = append(conduit.Shards, ConduitShard{
			ID: i,
			Transport: ConduitShardTransport{
				Method: "websocket",
			},
		})
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.conduits[conduit.ID] = conduit

	return cloneConduit(conduit)
}

func (s *State) UpdateConduitShards(conduitID string, shards []ConduitShard) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conduit, ok := s.conduits[conduitID]
	if !ok {
		return ErrConduitNotFound
	}

	indexByID := make(map[int]int, len(conduit.Shards))
	for idx, shard := range conduit.Shards {
		indexByID[shard.ID] = idx
	}

	for _, shard := range shards {
		if idx, exists := indexByID[shard.ID]; exists {
			conduit.Shards[idx].Transport = shard.Transport
			continue
		}

		conduit.Shards = append(conduit.Shards, shard)
	}

	s.conduits[conduitID] = conduit

	return nil
}

func (s *State) CreateSubscription(
	eventType string,
	version string,
	condition map[string]any,
	transport SubscriptionTransport,
) Subscription {
	if version == "" {
		version = defaultSubscriptionVersion(eventType)
	}

	if condition == nil {
		condition = defaultCondition()
	}

	if transport.Method == "" {
		transport.Method = "websocket"
	}

	subscription := Subscription{
		ID:        uuid.NewString(),
		Status:    "enabled",
		Type:      eventType,
		Version:   version,
		Cost:      0,
		Condition: cloneMap(condition),
		CreatedAt: time.Now().UTC(),
		Transport: transport,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscriptions[subscription.ID] = subscription

	return cloneSubscription(subscription)
}

func (s *State) DeleteSubscription(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subscriptions, id)
}

func (s *State) FindSubscriptionByType(eventType string) (Subscription, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, subscription := range s.subscriptions {
		if subscription.Type == eventType {
			return cloneSubscription(subscription), true
		}
	}

	return Subscription{}, false
}

func (s *State) SetStreamOnline(online bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if online && (!s.streamOnline || s.streamStarted.IsZero()) {
		s.streamStarted = time.Now().UTC()
	}

	if !online {
		s.streamStarted = time.Time{}
	}

	s.streamOnline = online
}

func (s *State) IsStreamOnline() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.streamOnline
}

func (s *State) SetStreamMeta(title, gameID, gameName string, viewerCount int) {
	if viewerCount < 0 {
		viewerCount = 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.streamTitle = title
	s.streamGameID = gameID
	s.streamGameName = gameName
	s.streamViewerCount = viewerCount
}

func (s *State) StreamMeta() (title, gameID, gameName string, viewerCount int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.streamTitle, s.streamGameID, s.streamGameName, s.streamViewerCount
}

func (s *State) SetFollowersTotal(n int) {
	if n < 0 {
		n = 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.followersTotal = n
}

func (s *State) FollowersTotal() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.followersTotal
}

func (s *State) SetUserFollowed(userID string, followed bool) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if followed {
		if _, exists := s.followedUserIDs[userID]; !exists {
			s.followedUserIDs[userID] = time.Now().UTC()
		}
		return
	}

	if userID == config.MockBotID {
		return
	}

	delete(s.followedUserIDs, userID)
	ensureDefaultStateLocked(s)
}

func (s *State) IsUserFollowed(userID string) (time.Time, bool) {
	userID = strings.TrimSpace(userID)

	s.mu.RLock()
	defer s.mu.RUnlock()

	followedAt, ok := s.followedUserIDs[userID]
	return followedAt, ok
}

func (s *State) ListFollowedUserIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, 0, len(s.followedUserIDs))
	for userID := range s.followedUserIDs {
		result = append(result, userID)
	}

	sort.Strings(result)
	return result
}

func (s *State) SetModerator(userID string, isMod bool) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if isMod || userID == config.MockBotID {
		s.moderatorIDs[userID] = true
		return
	}

	delete(s.moderatorIDs, userID)
	ensureDefaultStateLocked(s)
}

func (s *State) IsModerator(userID string) bool {
	userID = strings.TrimSpace(userID)

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.moderatorIDs[userID]
}

func (s *State) ListModerators() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, 0, len(s.moderatorIDs))
	for userID, isMod := range s.moderatorIDs {
		if !isMod {
			continue
		}

		result = append(result, userID)
	}

	sort.Strings(result)
	return result
}

func (s *State) StreamSnapshot() (map[string]any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.streamOnline {
		return nil, false
	}

	startedAt := s.streamStarted
	if startedAt.IsZero() {
		startedAt = time.Now().UTC()
	}

	return map[string]any{
		"id":            "mock-stream-id",
		"user_id":       config.MockBroadcasterID,
		"user_login":    config.MockBroadcasterLogin,
		"user_name":     config.MockBroadcasterName,
		"game_id":       s.streamGameID,
		"game_name":     s.streamGameName,
		"type":          "live",
		"title":         s.streamTitle,
		"viewer_count":  s.streamViewerCount,
		"started_at":    startedAt.Format(time.RFC3339),
		"language":      "en",
		"thumbnail_url": "",
		"is_mature":     false,
		"tags":          []string{},
		"tag_ids":       []string{},
		"delay":         0,
	}, true
}

func defaultSubscriptionVersion(eventType string) string {
	switch eventType {
	case "channel.follow":
		return "2"
	default:
		return "1"
	}
}

func defaultCondition() map[string]any {
	return map[string]any{
		"broadcaster_user_id": config.MockBroadcasterID,
		"moderator_user_id":   config.MockBotID,
	}
}

func cloneConduit(conduit Conduit) Conduit {
	clone := conduit
	clone.Shards = append([]ConduitShard(nil), conduit.Shards...)
	return clone
}

func cloneSubscription(subscription Subscription) Subscription {
	clone := subscription
	clone.Condition = cloneMap(subscription.Condition)
	return clone
}

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return nil
	}

	output := make(map[string]any, len(input))
	maps.Copy(output, input)

	return output
}

func ensureDefaultStateLocked(s *State) {
	if s.followedUserIDs == nil {
		s.followedUserIDs = map[string]time.Time{}
	}

	if _, ok := s.followedUserIDs[config.MockBotID]; !ok {
		s.followedUserIDs[config.MockBotID] = time.Now().UTC()
	}

	if s.moderatorIDs == nil {
		s.moderatorIDs = map[string]bool{}
	}

	s.moderatorIDs[config.MockBotID] = true
}
