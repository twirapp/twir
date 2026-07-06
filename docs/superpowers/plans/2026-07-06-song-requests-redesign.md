# Song Requests Redesign — Implementation Plan

> **For agentic workers:** Use superpowers:executing-plans or superpowers:subagent-driven-development to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rewrite YouTube song requests from raw melody WebSocket to GQL subscriptions + mutations, with multi-client sync, OBS overlay, and popup widget.

**Architecture:** Server-authoritative master state in Redis. Clients receive ready-to-use position via GQL subscriptions. Inter-service comms via bus-core (NATS), client pubsub via wsRouter (NATS) → GQL subscriptions.

**Tech Stack:** Go (pgx, gqlgen, NATS, Redis), Vue 3 + Nuxt (dashboard, overlays, widgets), YouTube IFrame API

**Spec:** `docs/superpowers/specs/2026-07-06-song-requests-redesign-design.md`

**NOTE:** Do NOT commit any files during implementation.

---

## File Map

| # | Area | Action | Path |
|---|------|--------|------|
| 1 | Migration | Create | `libs/migrations/postgres/YYYYMMDDHHMMSS_channel_api_key_and_hide_on_pause.sql` |
| 2 | Bus-core | Create | `libs/bus-core/api/song_requests.go` |
| 3 | Bus-core | Modify | `libs/bus-core/bus-services.go` |
| 4 | Bus-core | Modify | `libs/bus-core/bus.go` |
| 5 | Auth | Modify | `apps/api-gql/internal/auth/api_key.go` |
| 6 | Playback | Create | `apps/api-gql/internal/services/song_requests/playback_state.go` |
| 7 | Bridge | Create | `apps/api-gql/internal/services/song_requests/bridge.go` |
| 8 | GQL Schema | Modify | `apps/api-gql/internal/delivery/gql/schema/song-requests.graphql` |
| 9 | GQL Resolvers | Create/Modify | `apps/api-gql/internal/delivery/gql/resolvers/song-requests.resolver.go` |
| 10 | GQL Mappers | Create | `apps/api-gql/internal/delivery/gql/mappers/song_requests.go` |
| 11 | Entity | Modify | `apps/api-gql/internal/entity/song_request.go` |
| 12 | Parser sr | Modify | `apps/parser/internal/commands/songrequest/youtube/sr.go` |
| 13 | Parser skip | Modify | `apps/parser/internal/commands/songrequest/youtube/skip.go` |
| 14 | Parser wrong | Modify | `apps/parser/internal/commands/songrequest/youtube/wrong.go` |
| 15 | EventSub | Modify | `apps/eventsub/internal/handler/redemption.go` |
| 16 | Events | Modify | `apps/events/internal/song_request/song_request.go` |
| 17 | Dashboard | Create | `web/layers/dashboard/composables/useSongRequestGql.ts` |
| 18 | Dashboard | Modify | `web/layers/dashboard/components/songRequests/player.vue` |
| 19 | Dashboard | Modify | `web/layers/dashboard/components/songRequests/queue.vue` |
| 20 | Dashboard | Modify | `web/layers/dashboard/composables/useGlobalYoutubePlayer.ts` |
| 21 | Dashboard | Modify | `web/layers/dashboard/layout/sidebar/sidebar-mini-player.vue` |
| 22 | Dashboard | Delete | `web/layers/dashboard/components/songRequests/hook.ts` |
| 23 | Overlay | Create | `web/layers/overlays/pages/o/[apiKey]/song-requests.vue` |
| 24 | Widget | Create | `web/layers/widgets/nuxt.config.ts` |
| 25 | Widget | Create | `web/layers/widgets/pages/w/[channelApiKey]/song-requests.vue` |
| 26 | Websockets | Delete | `apps/websockets/internal/namespaces/youtube/` (entire dir) |
| 27 | Websockets | Delete | `apps/websockets/internal/grpc_impl/youtube.go` |

---

### Task 1: Database Migrations

**Files:**
- Create: `libs/migrations/postgres/20260706120000_channel_api_key_and_hide_on_pause.sql`

- [ ] **Step 1: Create migration file**

```sql
-- Channel API key
ALTER TABLE channels ADD COLUMN api_key TEXT DEFAULT uuidv7();
CREATE UNIQUE INDEX channels_api_key_idx ON channels(api_key) WHERE api_key IS NOT NULL;

-- Hide on pause setting for song requests
ALTER TABLE channels_song_requests_settings ADD COLUMN hide_on_pause BOOL DEFAULT true;
```

Save as `libs/migrations/postgres/20260706120000_channel_api_key_and_hide_on_pause.sql`.

---

### Task 2: Bus-Core Message Types

**Files:**
- Create: `libs/bus-core/api/song_requests.go`
- Modify: `libs/bus-core/bus-services.go` (add fields to `apiBus`)
- Modify: `libs/bus-core/bus.go` (instantiate queues in `NewNatsBus`)

- [ ] **Step 1: Create message types file**

Create `libs/bus-core/api/song_requests.go`:

```go
package api

const (
	SongRequestAddToQueueSubject      = "api.songRequests.addToQueue"
	SongRequestRemoveFromQueueSubject = "api.songRequests.removeFromQueue"
	SongRequestPlaybackStateSubject   = "api.songRequests.playbackState"
)

type SongRequestAddToQueue struct {
	ChannelID   string
	SongRequest SongRequestData
}

type SongRequestRemoveFromQueue struct {
	ChannelID string
	VideoID   string
}

type SongRequestPlaybackState struct {
	ChannelID string
	VideoID   string
	Title     string
	Position  float64
	IsPlaying bool
	Volume    int
	UpdatedAt int64
}

type SongRequestData struct {
	ID                   string
	Title                string
	VideoID              string
	SongLink             string
	DurationSeconds      int
	OrderedByName        string
	OrderedByDisplayName string
	QueuePosition        int
	CreatedAt            string
}
```

- [ ] **Step 2: Register in apiBus struct**

In `libs/bus-core/bus-services.go`, add to `apiBus` struct (after `TriggerObsCommand`):

```go
SongRequestAddToQueue      Queue[api.SongRequestAddToQueue, struct{}]
SongRequestRemoveFromQueue Queue[api.SongRequestRemoveFromQueue, struct{}]
SongRequestPlaybackState   Queue[api.SongRequestPlaybackState, struct{}]
```

- [ ] **Step 3: Instantiate in NewNatsBus**

In `libs/bus-core/bus.go`, inside the `Api: &apiBus{...}` block (after `TriggerObsCommand`), add:

```go
SongRequestAddToQueue: NewNatsQueue[api.SongRequestAddToQueue, struct{}](
	nc,
	api.SongRequestAddToQueueSubject,
	time.Second,
	GobEncoder,
),
SongRequestRemoveFromQueue: NewNatsQueue[api.SongRequestRemoveFromQueue, struct{}](
	nc,
	api.SongRequestRemoveFromQueueSubject,
	time.Second,
	GobEncoder,
),
SongRequestPlaybackState: NewNatsQueue[api.SongRequestPlaybackState, struct{}](
	nc,
	api.SongRequestPlaybackStateSubject,
	time.Second,
	GobEncoder,
),
```

- [ ] **Step 4: Verify build**

Run: `cd libs/bus-core && go build ./...`
Expected: no errors

---

### Task 3: Auth Middleware — Channel API Key Backward Compat

**Files:**
- Modify: `apps/api-gql/internal/auth/api_key.go`

- [ ] **Step 1: Read current implementation**

Current `GetAuthenticatedUserByApiKey` queries `users WHERE "apiKey" = ?`. Need to add channel API key check first.

- [ ] **Step 2: Add channel API key resolution method**

Add a new method to `apps/api-gql/internal/auth/api_key.go`:

```go
func (s *Auth) GetChannelByApiKey(ctx context.Context) (*model.Channels, error) {
	var apiKey string

	wsApiKey, _ := s.getWsAuthenticatedApiKey(ctx)
	if wsApiKey != "" {
		apiKey = wsApiKey
	} else {
		ginCtx, err := gincontext.GetGinContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get gin context: %w", err)
		}
		apiKey = ginCtx.GetHeader("api-key")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	channel := model.Channels{}
	if err := s.gorm.Where(`"api_key" = ?`, apiKey).First(&channel).Error; err != nil {
		return nil, fmt.Errorf("cannot get channel from db: %w", err)
	}

	return &channel, nil
}
```

- [ ] **Step 3: Update GetAuthenticatedUserByApiKey with fallback**

Modify the existing `GetAuthenticatedUserByApiKey` to try channel key first:

```go
func (s *Auth) GetAuthenticatedUserByApiKey(ctx context.Context) (*model.Users, error) {
	var apiKey string

	wsApiKey, _ := s.getWsAuthenticatedApiKey(ctx)
	if wsApiKey != "" {
		apiKey = wsApiKey
	} else {
		ginCtx, err := gincontext.GetGinContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get gin context: %w", err)
		}
		apiKey = ginCtx.GetHeader("api-key")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	// Try channel API key first
	channel := model.Channels{}
	if err := s.gorm.Where(`"api_key" = ?`, apiKey).First(&channel).Error; err == nil {
		// Found channel by API key, resolve owner user
		var userID string
		if channel.TwitchUserID != nil {
			userID = *channel.TwitchUserID
		} else if channel.KickUserID != nil {
			userID = *channel.KickUserID
		}
		if userID != "" {
			user := model.Users{}
			if err := s.gorm.Where("id = ?", userID).First(&user).Error; err == nil {
				return &user, nil
			}
		}
	}

	// Fallback to user API key
	user := model.Users{}
	if err := s.gorm.Where(`"apiKey" = ?`, apiKey).First(&user).Error; err != nil {
		return nil, fmt.Errorf("cannot get user from db: %w", err)
	}

	return &user, nil
}
```

- [ ] **Step 4: Verify build**

Run: `cd apps/api-gql && go build ./...`
Expected: no errors

---

### Task 4: Playback State Service

**Files:**
- Create: `apps/api-gql/internal/services/song_requests/playback_state.go`

- [ ] **Step 1: Create playback state service**

```go
package song_requests

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const playbackKeyPrefix = "songrequests:playback:"

type PlaybackState struct {
	VideoID   string  `json:"videoId"`
	Title     string  `json:"title"`
	Position  float64 `json:"position"`
	IsPlaying bool    `json:"isPlaying"`
	Volume    int     `json:"volume"`
	UpdatedAt int64   `json:"updatedAt"`
}

type PlaybackStateService struct {
	redis *redis.Client
}

func NewPlaybackStateService(redis *redis.Client) *PlaybackStateService {
	return &PlaybackStateService{redis: redis}
}

func (s *PlaybackStateService) key(channelID string) string {
	return playbackKeyPrefix + channelID
}

func (s *PlaybackStateService) SetPlaying(ctx context.Context, channelID string, videoID string, title string) error {
	state := PlaybackState{
		VideoID:   videoID,
		Title:     title,
		Position:  0,
		IsPlaying: true,
		Volume:    s.getVolumeOrDefault(ctx, channelID),
		UpdatedAt: time.Now().UnixMilli(),
	}
	return s.save(ctx, channelID, state)
}

func (s *PlaybackStateService) SetPaused(ctx context.Context, channelID string) error {
	state, err := s.GetState(ctx, channelID)
	if err != nil {
		return err
	}
	if state == nil {
		return nil
	}

	// Calculate current position
	if state.IsPlaying {
		elapsed := float64(time.Now().UnixMilli()-state.UpdatedAt) / 1000.0
		state.Position = state.Position + elapsed
	}
	state.IsPlaying = false
	state.UpdatedAt = time.Now().UnixMilli()

	return s.save(ctx, channelID, *state)
}

func (s *PlaybackStateService) SetVolume(ctx context.Context, channelID string, volume int) error {
	state, err := s.GetState(ctx, channelID)
	if err != nil {
		return err
	}
	if state == nil {
		return nil
	}

	// Update position if playing
	if state.IsPlaying {
		elapsed := float64(time.Now().UnixMilli()-state.UpdatedAt) / 1000.0
		state.Position = state.Position + elapsed
	}
	state.Volume = volume
	state.UpdatedAt = time.Now().UnixMilli()

	return s.save(ctx, channelID, *state)
}

func (s *PlaybackStateService) GetState(ctx context.Context, channelID string) (*PlaybackState, error) {
	data, err := s.redis.Get(ctx, s.key(channelID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("get playback state: %w", err)
	}

	var state PlaybackState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unmarshal playback state: %w", err)
	}

	// Compute live position if playing
	if state.IsPlaying {
		elapsed := float64(time.Now().UnixMilli()-state.UpdatedAt) / 1000.0
		state.Position = state.Position + elapsed
		state.UpdatedAt = time.Now().UnixMilli()
	}

	return &state, nil
}

func (s *PlaybackStateService) ClearState(ctx context.Context, channelID string) error {
	return s.redis.Del(ctx, s.key(channelID)).Err()
}

func (s *PlaybackStateService) getVolumeOrDefault(ctx context.Context, channelID string) int {
	// Try to get saved volume, default to 100
	data, err := s.redis.Get(ctx, s.key(channelID)).Bytes()
	if err == nil {
		var state PlaybackState
		if json.Unmarshal(data, &state) == nil {
			return state.Volume
		}
	}
	return 100
}

func (s *PlaybackStateService) save(ctx context.Context, channelID string, state PlaybackState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal playback state: %w", err)
	}
	return s.redis.Set(ctx, s.key(channelID), data, 0).Err()
}
```

- [ ] **Step 2: Verify build**

Run: `cd apps/api-gql && go build ./...`
Expected: no errors

---

### Task 5: GraphQL Schema

**Files:**
- Modify: `apps/api-gql/internal/delivery/gql/schema/song-requests.graphql`

- [ ] **Step 1: Add new types, subscriptions, mutations, and queries**

Append to `apps/api-gql/internal/delivery/gql/schema/song-requests.graphql`:

```graphql
extend type Subscription {
	songRequestPlaybackState(channelId: UUID!): SongRequestPlaybackState
	songRequestQueueUpdated(channelId: UUID!): [SongRequestQueueItem!]!
}

type SongRequestPlaybackState {
	videoId: String!
	title: String!
	position: Float!
	isPlaying: Boolean!
	volume: Int!
	updatedAt: Time!
}

type SongRequestQueueItem {
	id: String!
	title: String!
	songLink: String!
	durationSeconds: Int!
	orderedByName: String!
	orderedByDisplayName: String!
	queuePosition: Int!
	createdAt: Time!
}
```

Also add to the existing `SongRequestsSettings` type:

```graphql
hideOnPause: Boolean!
```

And add to `SongRequestsSettingsOpts` input:

```graphql
hideOnPause: Boolean!
```

Add new mutations:

```graphql
extend type Mutation {
	songRequestPlay(channelId: UUID!, videoId: String!): Boolean! @isAuthenticated
	songRequestPause(channelId: UUID!): Boolean! @isAuthenticated
	songRequestSkip(channelId: UUID!): Boolean! @isAuthenticated
	songRequestSetVolume(channelId: UUID!, volume: Int!): Boolean! @isAuthenticated
	songRequestReorder(channelId: UUID!, videoIds: [String!]!): Boolean! @isAuthenticated
	songRequestDeleteFromQueue(channelId: UUID!, videoId: String!): Boolean! @isAuthenticated
	songRequestClearQueue(channelId: UUID!): Boolean! @isAuthenticated
}
```

Add new queries:

```graphql
extend type Query {
	songRequestWidgetData(channelId: UUID!): SongRequestWidgetData
	channelByApiKey(apiKey: String!): ChannelByApiKeyResult
}

type SongRequestWidgetData {
	playbackState: SongRequestPlaybackState
	queue: [SongRequestQueueItem!]!
}

type ChannelByApiKeyResult {
	id: UUID!
	twitchUserId: UUID
	kickUserId: UUID
}
```

- [ ] **Step 2: Regenerate GQL resolvers**

Run: `bun cli build gql`
Expected: new resolver stubs generated

---

### Task 6: Bridge Service (bus-core → wsRouter)

**Files:**
- Create: `apps/api-gql/internal/services/song_requests/bridge.go`

- [ ] **Step 1: Create bridge service**

```go
package song_requests

import (
	"context"
	"log/slog"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type BridgeOpts struct {
	fx.In
	LC      fx.Lifecycle
	WsRouter wsrouter.WsRouter
	TwirBus  *buscore.Bus
	Logger   *slog.Logger
}

type Bridge struct {
	wsRouter wsrouter.WsRouter
	twirBus  *buscore.Bus
	logger   *slog.Logger
}

func NewBridge(opts BridgeOpts) *Bridge {
	b := &Bridge{
		wsRouter: opts.WsRouter,
		twirBus:  opts.TwirBus,
		logger:   opts.Logger,
	}

	opts.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			b.twirBus.Api.SongRequestAddToQueue.SubscribeGroup("api",
				func(ctx context.Context, data api.SongRequestAddToQueue) (struct{}, error) {
					return struct{}{}, b.wsRouter.Publish(
						"api.songRequestQueue."+data.ChannelID, data,
					)
				},
			)
			b.logger.Info("Subscribed to SongRequestAddToQueue events")

			b.twirBus.Api.SongRequestRemoveFromQueue.SubscribeGroup("api",
				func(ctx context.Context, data api.SongRequestRemoveFromQueue) (struct{}, error) {
					return struct{}{}, b.wsRouter.Publish(
						"api.songRequestQueueRemove."+data.ChannelID, data,
					)
				},
			)
			b.logger.Info("Subscribed to SongRequestRemoveFromQueue events")

			b.twirBus.Api.SongRequestPlaybackState.SubscribeGroup("api",
				func(ctx context.Context, data api.SongRequestPlaybackState) (struct{}, error) {
					return struct{}{}, b.wsRouter.Publish(
						"api.songRequestPlayback."+data.ChannelID, data,
					)
				},
			)
			b.logger.Info("Subscribed to SongRequestPlaybackState events")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			b.twirBus.Api.SongRequestAddToQueue.Unsubscribe()
			b.twirBus.Api.SongRequestRemoveFromQueue.Unsubscribe()
			b.twirBus.Api.SongRequestPlaybackState.Unsubscribe()
			return nil
		},
	})

	return b
}
```

- [ ] **Step 2: Verify build**

Run: `cd apps/api-gql && go build ./...`
Expected: no errors

---

### Task 7: GQL Resolvers and Entity Updates

**Files:**
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/song-requests.resolver.go`
- Modify: `apps/api-gql/internal/entity/song_request.go`
- Create: `apps/api-gql/internal/delivery/gql/mappers/song_requests_playback.go`

- [ ] **Step 1: Add entity types**

Add to `apps/api-gql/internal/entity/song_request.go`:

```go
type SongRequestPlaybackState struct {
	VideoID   string
	Title     string
	Position  float64
	IsPlaying bool
	Volume    int
	UpdatedAt int64
}

type SongRequestQueueItem struct {
	ID                   string
	Title                string
	SongLink             string
	DurationSeconds      int
	OrderedByName        string
	OrderedByDisplayName string
	QueuePosition        int
	CreatedAt            string
}

type SongRequestWidgetData struct {
	PlaybackState *SongRequestPlaybackState
	Queue         []SongRequestQueueItem
}

type ChannelByApiKeyResult struct {
	ID           string
	TwitchUserID *string
	KickUserID   *string
}
```

- [ ] **Step 2: Implement subscription resolvers**

After `bun cli build gql`, implement the generated resolver stubs. The pattern follows TTS (`apps/api-gql/internal/services/overlays/tts/subscription.go`):

Subscriptions subscribe to wsRouter keys and return channels. The resolver goroutine reads from wsRouter, unmarshals, maps to GQL model, and sends to output channel.

For `SongRequestPlaybackState`: subscribe to `api.songRequestPlayback.{channelId}`, also send initial state from Redis.

For `SongRequestQueueUpdated`: subscribe to `api.songRequestQueue.{channelId}` and `api.songRequestQueueRemove.{channelId}`, send full queue from DB on subscribe.

- [ ] **Step 3: Implement mutation resolvers**

Mutations call PlaybackStateService methods and publish to bus-core:

- `SongRequestPlay` → `playbackStateService.SetPlaying()` → publish to bus-core
- `SongRequestPause` → `playbackStateService.SetPaused()` → publish to bus-core
- `SongRequestSkip` → soft-delete in DB, `playbackStateService.ClearState()`, publish to bus-core
- `SongRequestSetVolume` → `playbackStateService.SetVolume()` → publish to bus-core
- `SongRequestReorder` → update queue positions in DB, publish to bus-core
- `SongRequestDeleteFromQueue` → soft-delete in DB, publish to bus-core
- `SongRequestClearQueue` → soft-delete all in DB, `playbackStateService.ClearState()`, publish to bus-core

- [ ] **Step 4: Implement query resolvers**

- `SongRequestWidgetData` → reads playback state from Redis + queue from DB
- `ChannelByApiKey` → queries `channels WHERE api_key = ?`

- [ ] **Step 5: Regenerate and verify**

Run: `bun cli build gql && cd apps/api-gql && go build ./...`
Expected: no errors

---

### Task 8: Parser/EventSub/Events → bus-core

**Files:**
- Modify: `apps/parser/internal/commands/songrequest/youtube/sr.go`
- Modify: `apps/parser/internal/commands/songrequest/youtube/skip.go`
- Modify: `apps/parser/internal/commands/songrequest/youtube/wrong.go`
- Modify: `apps/eventsub/internal/handler/redemption.go`
- Modify: `apps/events/internal/song_request/song_request.go`

- [ ] **Step 1: Update parser sr.go**

Replace lines 200-208 (the gRPC `YoutubeAddSongToQueue` loop) with bus-core publish:

```go
for _, song := range requested {
	songLink := ""
	if song.SongLink.Valid {
		songLink = song.SongLink.String
	}

	parseCtx.Services.Bus.Api.SongRequestAddToQueue.Publish(
		ctx,
		api.SongRequestAddToQueue{
			ChannelID: parseCtx.Channel.DBChannelID,
			SongRequest: api.SongRequestData{
				ID:                   song.ID,
				Title:                song.Title,
				VideoID:              song.VideoID,
				SongLink:             songLink,
				DurationSeconds:      int(song.Duration),
				OrderedByName:        song.OrderedByName,
				OrderedByDisplayName: song.OrderedByDisplayName.String,
				QueuePosition:        song.QueuePosition,
				CreatedAt:            song.CreatedAt.UTC().String(),
			},
		},
	)
}
```

Add import for `"github.com/twirapp/twir/libs/bus-core/api"`.

- [ ] **Step 2: Update parser skip.go**

Replace the gRPC `YoutubeRemoveSongToQueue` call (lines 120-133) with:

```go
parseCtx.Services.Bus.Api.SongRequestRemoveFromQueue.Publish(
	ctx,
	api.SongRequestRemoveFromQueue{
		ChannelID: parseCtx.Channel.DBChannelID,
		VideoID:   currentSong.VideoID,
	},
)
```

Add import for `"github.com/twirapp/twir/libs/bus-core/api"`.

- [ ] **Step 3: Update parser wrong.go**

Replace the gRPC `YoutubeRemoveSongToQueue` call (lines 95-107) with:

```go
parseCtx.Services.Bus.Api.SongRequestRemoveFromQueue.Publish(
	ctx,
	api.SongRequestRemoveFromQueue{
		ChannelID: parseCtx.Channel.DBChannelID,
		VideoID:   choosedSong.VideoID,
	},
)
```

Add import for `"github.com/twirapp/twir/libs/bus-core/api"`.

- [ ] **Step 4: Update eventsub redemption.go**

The eventsub handler delegates to parser via bus-core (`Parser.ProcessMessageAsCommand`), so it doesn't directly call gRPC for song requests. No changes needed here — the parser handles the actual song request creation.

- [ ] **Step 5: Update events song_request.go**

The events service also delegates to parser via bus-core (`Parser.ProcessMessageAsCommand`). No changes needed — same reason as eventsub.

- [ ] **Step 6: Verify build**

Run: `cd apps/parser && go build ./...`
Expected: no errors

---

### Task 9: Playback Ticker

**Files:**
- Modify: `apps/api-gql/internal/services/song_requests/playback_state.go`

- [ ] **Step 1: Add ticker method**

Add to `PlaybackStateService`:

```go
func (s *PlaybackStateService) StartTicker(ctx context.Context, wsRouter wsrouter.WsRouter) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.tickAllChannels(ctx, wsRouter)
			}
		}
	}()
}

func (s *PlaybackStateService) tickAllChannels(ctx context.Context, wsRouter wsrouter.WsRouter) {
	keys, err := s.redis.Keys(ctx, playbackKeyPrefix+"*").Result()
	if err != nil {
		return
	}

	for _, key := range keys {
		channelID := key[len(playbackKeyPrefix):]
		state, err := s.GetState(ctx, channelID)
		if err != nil || state == nil || !state.IsPlaying {
			continue
		}

		// Publish updated state to wsRouter
		wsRouter.Publish("api.songRequestPlayback."+channelID, api.SongRequestPlaybackState{
			ChannelID: channelID,
			VideoID:   state.VideoID,
			Title:     state.Title,
			Position:  state.Position,
			IsPlaying: state.IsPlaying,
			Volume:    state.Volume,
			UpdatedAt: state.UpdatedAt,
		})
	}
}
```

- [ ] **Step 2: Start ticker in bridge service**

In `apps/api-gql/internal/services/song_requests/bridge.go`, in the `OnStart` hook, add:

```go
playbackService.StartTicker(ctx, opts.WsRouter)
```

This requires passing `PlaybackStateService` to the bridge. Update `BridgeOpts` and `NewBridge` accordingly.

- [ ] **Step 3: Verify build**

Run: `cd apps/api-gql && go build ./...`
Expected: no errors

---

### Task 10: Dashboard — New Composable

**Files:**
- Create: `web/layers/dashboard/composables/useSongRequestGql.ts`

- [ ] **Step 1: Create composable**

```typescript
import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { computed } from 'vue'
import { graphql } from '~/gql'

const SongRequestPlaybackStateSubscription = graphql(`
  subscription SongRequestPlaybackState($channelId: UUID!) {
    songRequestPlaybackState(channelId: $channelId) {
      videoId
      title
      position
      isPlaying
      volume
      updatedAt
    }
  }
`)

const SongRequestQueueUpdatedSubscription = graphql(`
  subscription SongRequestQueueUpdated($channelId: UUID!) {
    songRequestQueueUpdated(channelId: $channelId) {
      id
      title
      songLink
      durationSeconds
      orderedByName
      orderedByDisplayName
      queuePosition
      createdAt
    }
  }
`)

const SongRequestPlayMutation = graphql(`
  mutation SongRequestPlay($channelId: UUID!, $videoId: String!) {
    songRequestPlay(channelId: $channelId, videoId: $videoId)
  }
`)

const SongRequestPauseMutation = graphql(`
  mutation SongRequestPause($channelId: UUID!) {
    songRequestPause(channelId: $channelId)
  }
`)

const SongRequestSkipMutation = graphql(`
  mutation SongRequestSkip($channelId: UUID!) {
    songRequestSkip(channelId: $channelId)
  }
`)

const SongRequestSetVolumeMutation = graphql(`
  mutation SongRequestSetVolume($channelId: UUID!, $volume: Int!) {
    songRequestSetVolume(channelId: $channelId, volume: $volume)
  }
`)

const SongRequestReorderMutation = graphql(`
  mutation SongRequestReorder($channelId: UUID!, $videoIds: [String!]!) {
    songRequestReorder(channelId: $channelId, videoIds: $videoIds)
  }
`)

const SongRequestDeleteFromQueueMutation = graphql(`
  mutation SongRequestDeleteFromQueue($channelId: UUID!, $videoId: String!) {
    songRequestDeleteFromQueue(channelId: $channelId, videoId: $videoId)
  }
`)

const SongRequestClearQueueMutation = graphql(`
  mutation SongRequestClearQueue($channelId: UUID!) {
    songRequestClearQueue(channelId: $channelId)
  }
`)

export function useSongRequestGql(channelId: Ref<string>) {
  const playbackStateSub = useSubscription({
    query: SongRequestPlaybackStateSubscription,
    variables: computed(() => ({ channelId: channelId.value })),
  })

  const queueSub = useSubscription({
    query: SongRequestQueueUpdatedSubscription,
    variables: computed(() => ({ channelId: channelId.value })),
  })

  const playbackState = computed(() => playbackStateSub.data.value?.songRequestPlaybackState ?? null)
  const queue = computed(() => queueSub.data.value?.songRequestQueueUpdated ?? [])

  const { mutate: play } = useMutation(SongRequestPlayMutation)
  const { mutate: pause } = useMutation(SongRequestPauseMutation)
  const { mutate: skip } = useMutation(SongRequestSkipMutation)
  const { mutate: setVolume } = useMutation(SongRequestSetVolumeMutation)
  const { mutate: reorder } = useMutation(SongRequestReorderMutation)
  const { mutate: deleteFromQueue } = useMutation(SongRequestDeleteFromQueueMutation)
  const { mutate: clearQueue } = useMutation(SongRequestClearQueueMutation)

  return {
    playbackState,
    queue,
    play: (videoId: string) => play({ channelId: channelId.value, videoId }),
    pause: () => pause({ channelId: channelId.value }),
    skip: () => skip({ channelId: channelId.value }),
    setVolume: (volume: number) => setVolume({ channelId: channelId.value, volume }),
    reorder: (videoIds: string[]) => reorder({ channelId: channelId.value, videoIds }),
    deleteFromQueue: (videoId: string) => deleteFromQueue({ channelId: channelId.value, videoId }),
    clearQueue: () => clearQueue({ channelId: channelId.value }),
  }
}
```

---

### Task 11: Dashboard — Player Refactoring

**Files:**
- Modify: `web/layers/dashboard/components/songRequests/player.vue`

- [ ] **Step 1: Replace useYoutubeSocket with useSongRequestGql**

Replace the import of `useYoutubeSocket` (from `hook.ts`) with `useSongRequestGql`.

The player component should:
1. Get `channelId` from the selected dashboard context
2. Use `useSongRequestGql(channelId)` for state and actions
3. Use `useGlobalYoutubePlayer()` for the YouTube IFrame API
4. When `playbackState` changes from subscription → apply to player: `seekTo(position)`, `setVolume(volume)`, `playVideo()` or `pauseVideo()`
5. Remove all local queue management logic (server is authority)
6. Wire skip/play/pause/volume buttons to mutation functions

- [ ] **Step 2: Update useGlobalYoutubePlayer.ts**

Remove the queue management and auto-advance logic. The player should:
- Still create and manage the YouTube IFrame API player instance
- Expose `seekTo`, `playVideo`, `pauseVideo`, `setVolume`, `loadVideoById`, `cueVideoById`
- NOT decide when to switch tracks — server sends new videoId via subscription
- When `playbackState.videoId` changes → `loadVideoById(videoId)` + `seekTo(position)`

- [ ] **Step 3: Verify dashboard compiles**

Run: `cd web && bun run typecheck` (or equivalent)
Expected: no errors

---

### Task 12: Dashboard — Queue and Mini-Player

**Files:**
- Modify: `web/layers/dashboard/components/songRequests/queue.vue`
- Modify: `web/layers/dashboard/layout/sidebar/sidebar-mini-player.vue`

- [ ] **Step 1: Update queue.vue**

Replace raw WebSocket data with `useSongRequestGql().queue`. Wire delete/reorder/clear buttons to mutation functions.

- [ ] **Step 2: Update sidebar-mini-player.vue**

Connect to `useSongRequestGql` for playback state and controls. Show current track title, progress, play/pause/skip/volume.

- [ ] **Step 3: Delete old hook.ts**

Delete `web/layers/dashboard/components/songRequests/hook.ts`.

- [ ] **Step 4: Verify dashboard compiles**

Run: `cd web && bun run typecheck`
Expected: no errors

---

### Task 13: OBS Overlay

**Files:**
- Create: `web/layers/overlays/pages/o/[apiKey]/song-requests.vue`

- [ ] **Step 1: Create overlay page**

```vue
<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery, useSubscription } from '@urql/vue'
import { graphql } from '~/gql'

const route = useRoute()
const apiKey = computed(() => route.params.apiKey as string)

// Resolve channel ID from API key
const channelQuery = useQuery({
  query: graphql(`
    query ChannelByApiKey($apiKey: String!) {
      channelByApiKey(apiKey: $apiKey) {
        id
      }
    }
  `),
  variables: computed(() => ({ apiKey: apiKey.value })),
  pause: computed(() => !apiKey.value),
})

const channelId = computed(() => channelQuery.data.value?.channelByApiKey?.id ?? '')

// Settings
const settingsQuery = useQuery({
  query: graphql(`
    query SongRequestOverlaySettings($channelId: UUID!) {
      songRequests(channelId: $channelId) {
        hideOnPause
      }
    }
  `),
  variables: computed(() => ({ channelId: channelId.value })),
  pause: computed(() => !channelId.value),
})

const hideOnPause = computed(() => settingsQuery.data.value?.songRequests?.hideOnPause ?? true)

// Playback state subscription
const playbackSub = useSubscription({
  query: graphql(`
    subscription SongRequestPlaybackState($channelId: UUID!) {
      songRequestPlaybackState(channelId: $channelId) {
        videoId
        title
        position
        isPlaying
        volume
        updatedAt
      }
    }
  `),
  variables: computed(() => ({ channelId: channelId.value })),
})

const playbackState = computed(() => playbackSub.data.value?.songRequestPlaybackState ?? null)

// YouTube player
const playerReady = ref(false)
let player: any = null

onMounted(() => {
  // Load YouTube IFrame API
  const tag = document.createElement('script')
  tag.src = 'https://www.youtube.com/iframe_api'
  document.head.appendChild(tag)

  ;(window as any).onYouTubeIframeAPIReady = () => {
    player = new (window as any).YT.Player('yt-player', {
      height: '100%',
      width: '100%',
      playerVars: {
        autoplay: 0,
        controls: 0,
        modestbranding: 1,
        rel: 0,
      },
      events: {
        onReady: () => { playerReady.value = true },
      },
    })
  }
})

// Apply state from subscription
watch(playbackState, (state) => {
  if (!state || !playerReady.value) return

  player.seekTo(state.position, true)
  player.setVolume(state.volume)

  if (state.isPlaying) {
    player.playVideo()
  } else {
    player.pauseVideo()
  }
})

// Visibility logic
const isVisible = computed(() => {
  if (!playbackState.value) return false
  if (playbackState.value.isPlaying) return true
  return !hideOnPause.value
})

// Format time
function formatTime(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>

<template>
  <div v-if="isVisible" class="song-request-overlay">
    <div class="player-container">
      <div id="yt-player" class="yt-player" />
    </div>
    <div class="track-info">
      <div class="track-title">{{ playbackState?.title }}</div>
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: /* calculated */ '0%' }" />
      </div>
      <div class="progress-time">{{ formatTime(playbackState?.position ?? 0) }}</div>
    </div>
  </div>
</template>

<style scoped>
.song-request-overlay {
  position: relative;
  width: 100%;
  max-width: 640px;
  background: rgba(0, 0, 0, 0.85);
  border-radius: 8px;
  overflow: hidden;
}

.player-container {
  position: relative;
  width: 100%;
  padding-top: 56.25%; /* 16:9 */
}

.yt-player {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.track-info {
  padding: 12px 16px;
}

.track-title {
  color: white;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  margin-bottom: 4px;
}

.progress-fill {
  height: 100%;
  background: #8b5cf6;
  border-radius: 2px;
  transition: width 0.5s linear;
}

.progress-time {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
}
</style>
```

- [ ] **Step 2: Verify overlay compiles**

Run: `cd web && bun run typecheck`
Expected: no errors

---

### Task 14: Widget Layer

**Files:**
- Create: `web/layers/widgets/nuxt.config.ts`
- Create: `web/layers/widgets/pages/w/[channelApiKey]/song-requests.vue`

- [ ] **Step 1: Create layer config**

`web/layers/widgets/nuxt.config.ts`:

```typescript
export default defineNuxtConfig({
  // Minimal config, inherits from root
})
```

- [ ] **Step 2: Create widget page**

`web/layers/widgets/pages/w/[channelApiKey]/song-requests.vue`:

Similar to OBS overlay but with:
- YouTube player (visible, 16:9)
- Track title + author + progress bar
- Play/pause/skip/volume controls
- Queue list with delete buttons + clear all
- Compact layout (~400x600)

Uses same pattern: `channelByApiKey` query → `channelId` → subscriptions + mutations.

- [ ] **Step 3: Register layer in root nuxt.config**

Add `widgets` to the `extends` or `layers` array in `web/nuxt.config.ts` if not auto-discovered.

- [ ] **Step 4: Verify widget compiles**

Run: `cd web && bun run typecheck`
Expected: no errors

---

### Task 15: Delete Old Melody YouTube Namespace

**Files:**
- Delete: `apps/websockets/internal/namespaces/youtube/` (entire directory)
- Delete: `apps/websockets/internal/grpc_impl/youtube.go`

- [ ] **Step 1: Delete youtube namespace directory**

Delete all files in `apps/websockets/internal/namespaces/youtube/`.

- [ ] **Step 2: Delete youtube gRPC impl**

Delete `apps/websockets/internal/grpc_impl/youtube.go`.

- [ ] **Step 3: Remove gRPC registration**

In `apps/websockets/internal/grpc_impl/impl.go`, remove the YouTube-related gRPC service registration.

- [ ] **Step 4: Remove youtube namespace from app wiring**

In `apps/websockets/app/app.go`, remove the YouTube namespace module.

- [ ] **Step 5: Verify websockets builds**

Run: `cd apps/websockets && go build ./...`
Expected: no errors

---

### Task 16: i18n Keys

**Files:**
- Modify: `web/locales/` (add keys for overlay and widget UI)

- [ ] **Step 1: Add i18n keys**

Add keys for:
- `songRequests.overlay.hideOnPause` — setting label
- `songRequests.widget.queue` — "Queue" header
- `songRequests.widget.clearAll` — "Clear All" button
- `songRequests.widget.empty` — "Queue is empty" message

---

## Verification Checklist

After all tasks:

1. `cd libs/bus-core && go build ./...`
2. `cd apps/api-gql && go build ./...`
3. `cd apps/parser && go build ./...`
4. `cd apps/websockets && go build ./...`
5. `cd web && bun run typecheck`
6. `bun lint`
7. Manual test: open dashboard → play song → open OBS overlay → verify sync
8. Manual test: open popup widget → verify controls work
9. Manual test: pause → verify overlay behavior (hideOnPause setting)
