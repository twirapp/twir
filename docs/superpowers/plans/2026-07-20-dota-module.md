# Dota 2 Module Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Dotabod-like Dota 2 module in Twir: GSI ingestion, chat commands, auto Twitch predictions, overlays, chat events.

**Architecture:** New microservice `apps/dota` (fx, stdlib HTTP mux) ingests GSI events, runs a match state machine, talks to Stratz/OpenDota, and communicates via bus-core (NATS): request/reply for parser commands, events for chat alerts, `api.dota.state` for overlay GraphQL subscriptions. api-gql gets a `dota` domain (settings, Steam OpenID linking, GSI config, subscription bridging). Dashboard lives in `web/layers/dashboard` (Nuxt), overlays in `frontend/overlays`.

**Tech Stack:** Go (fx, pgx, goose, NATS), gqlgen, Vue 3 (Nuxt dashboard layer + Vite overlays app), urql graphql-ws subscriptions.

**Spec:** `docs/superpowers/specs/2026-07-20-dota-module-design.md`

**Key codebase facts (verified by research):**
- Dashboard is `web/layers/dashboard` (Nuxt layer, file-based routing), NOT `frontend/dashboard`.
- Overlays get live data via GraphQL subscriptions through api-gql WsRouter (kappagen pattern) — no changes to `apps/websockets` needed.
- Parser reads module settings via `parseCtx.Services.Gorm` + `libs/gomodels` (e.g. `apps/parser/internal/commands/games/russian_roulette.go:44`).
- Per-channel Twitch Helix client: `twitch.NewUserClientWithContext(ctx, twitchUserID, config, bus)` (libs/twitch) — used in `apps/parser/internal/commands/predictions/start.go:68`. Scopes `channel:manage:predictions` already requested.
- Migrations: `bun cli m create --name X --db postgres --type sql`, files land in `libs/migrations/postgres/`, use `uuidv7()` for UUID defaults (Postgres 18).
- New app checklist: `go.work`, `cli/internal/goapp/apps.go`, `cli/internal/cmds/deploy/deploy.go` (serviceImages + releaseServices), `.github/workflows/dockerv3.yml` matrix, `docker-compose.stack.yml`, `apps/dota/Dockerfile`.
- Old dota gomodels exist (`channels_dota_accounts`, `dota_matches`, ...) and old dota parser commands are commented-out dead code — do NOT revive; build fresh table `channels_dota_settings`.
- Legacy commented code in `apps/parser/internal/commands/dota/` should be deleted as part of Task 10.

---

### Task 1: Database migration — `channels_dota_settings`

**Files:**
- Create: `libs/migrations/postgres/<timestamp>_channels_dota_settings.sql` (via CLI)

- [ ] **Step 1: Create migration**

Run: `bun cli m create --name channels_dota_settings --db postgres --type sql`

Fill the generated file:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS channels_dota_settings (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id TEXT NOT NULL UNIQUE REFERENCES channels(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT false,
    steam_account_id TEXT,
    gsi_token TEXT NOT NULL DEFAULT replace(uuidv7()::text, '-', ''),
    mmr INT NOT NULL DEFAULT 0,
    mmr_delta INT NOT NULL DEFAULT 25,
    session_wins INT NOT NULL DEFAULT 0,
    session_losses INT NOT NULL DEFAULT 0,
    prediction_settings JSONB NOT NULL DEFAULT '{"enabled": false, "titleTemplate": "Win this game?", "windowSeconds": 300}'::jsonb,
    chat_events JSONB NOT NULL DEFAULT '{}'::jsonb,
    commands_settings JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS channels_dota_settings_gsi_token_idx ON channels_dota_settings(gsi_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS channels_dota_settings;
-- +goose StatementEnd
```

- [ ] **Step 2: Apply migration**

Run: `bun cli m up` (check `cli/internal/cmds/migrations` for exact subcommand; the CLI exposes migrate commands)
Expected: migration applied without errors.

- [ ] **Step 3: Commit**

```bash
git add libs/migrations/postgres/*channels_dota_settings.sql
git commit -m "feat: add channels_dota_settings migration"
```

---

### Task 2: Entity, gomodel, repository

**Files:**
- Create: `libs/entities/dota/entity.go`
- Create: `libs/gomodels/channels_dota_settings.go`
- Create: `libs/repositories/dota/repository.go`
- Create: `libs/repositories/dota/model/model.go`
- Create: `libs/repositories/dota/pgx/pgx.go`

- [ ] **Step 1: Entity** (`libs/entities/dota/entity.go`) — follow `libs/entities/secret/entity.go` pattern (isNil/IsNil/Nil):

```go
package dota

import (
	"time"

	"github.com/google/uuid"
)

type PredictionSettings struct {
	Enabled       bool   `json:"enabled"`
	TitleTemplate string `json:"titleTemplate"`
	WindowSeconds int    `json:"windowSeconds"`
}

type ChatEventSettings struct {
	Enabled  bool   `json:"enabled"`
	Template string `json:"template"`
	Cooldown int    `json:"cooldown"`
}

type ChatEvents struct {
	MatchStarted ChatEventSettings `json:"matchStarted"`
	MatchEnded   ChatEventSettings `json:"matchEnded"`
	RoshanKilled ChatEventSettings `json:"roshanKilled"`
	AegisPickup  ChatEventSettings `json:"aegisPickup"`
}

type CommandsSettings struct {
	Mmr bool `json:"mmr"`
	Wl  bool `json:"wl"`
	Lg  bool `json:"lg"`
	Gm  bool `json:"gm"`
	Np  bool `json:"np"`
	Wp  bool `json:"wp"`
}

type ChannelDotaSettings struct {
	ID                 uuid.UUID
	ChannelID          string
	Enabled            bool
	SteamAccountID     *string
	GsiToken           string
	Mmr                int
	MmrDelta           int
	SessionWins        int
	SessionLosses      int
	PredictionSettings PredictionSettings
	ChatEvents         ChatEvents
	CommandsSettings   CommandsSettings
	CreatedAt          time.Time
	UpdatedAt          time.Time

	isNil bool
}

func (c ChannelDotaSettings) IsNil() bool { return c.isNil }

var Nil = ChannelDotaSettings{isNil: true}

func (c ChannelDotaSettings) Winrate() float64 {
	total := c.SessionWins + c.SessionLosses
	if total == 0 {
		return 0
	}
	return float64(c.SessionWins) / float64(total) * 100
}
```

- [ ] **Step 2: gomodel** (`libs/gomodels/channels_dota_settings.go`) — gorm tags style like `libs/gomodels/channel_games_russian_roulette.go`, `TableName() = "channels_dota_settings"`. JSONB columns as `[]byte`/`json.RawMessage` or typed structs with Scan/Value like `ChannelModulesSettingsColumn` in `libs/gomodels/channel_modules_settings.go:10-38`. Used by parser-side reads.

- [ ] **Step 3: Repository** — follow `libs/repositories/overlays_kappagen/` exactly:
  - `repository.go`: interface `GetByChannelID`, `GetByGsiToken`, `Create`, `Update`, inputs, `ErrNotFound`.
  - `model/model.go`: model struct with `isNil`/`IsNil()`/`var Nil` (new convention, see `libs/repositories/AGENTS.md`).
  - `pgx/pgx.go`: `Opts{PgxPool}`, `New`, `NewFx`, `var _ dota.Repository = (*Pgx)(nil)`, `trmpgx.DefaultCtxGetter`, raw SQL, `errors.Is(err, pgx.ErrNoRows)` → `ErrNotFound`, JSONB via `goccy/go-json`.

- [ ] **Step 4: Compile check**

Run: `cd libs/repositories && go build ./... && cd ../entities && go build ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add libs/entities/dota libs/gomodels/channels_dota_settings.go libs/repositories/dota
git commit -m "feat: dota settings entity, gomodel and pgx repository"
```

---

### Task 3: bus-core dota subjects

**Files:**
- Create: `libs/bus-core/dota/dota.go`
- Modify: `libs/bus-core/bus-services.go` (add `dotaBus` struct + field on apiBus)
- Modify: `libs/bus-core/bus.go` (Bus struct field + NewNatsBus wiring)
- Modify: `libs/bus-core/api/kappagen.go`-sibling new file `libs/bus-core/api/dota.go`

- [ ] **Step 1: Define subjects** (`libs/bus-core/dota/dota.go`):

```go
package dota

const (
	GetDataSubject       = "dota.get_data"
	MatchStartedSubject  = "dota.match_started"
	MatchEndedSubject    = "dota.match_ended"
	RoshanKilledSubject  = "dota.roshan_killed"
	AegisPickupSubject   = "dota.aegis_pickup"
)

type HeroInfo struct {
	Name string
}

type PlayerInfo struct {
	AccountID int64
	HeroID    int
	Kills     int
	Deaths    int
	Assists   int
	IsStreamer bool
}

type GetDataRequest struct {
	ChannelID      string // internal channel uuid
	TwitchUserID   string
	SteamAccountID string
}

type GetDataResponse struct {
	Enabled        bool
	Linked         bool
	InGame         bool
	Mmr            int
	SessionWins    int
	SessionLosses  int
	HeroName       string
	MatchID        int64
	TeamIsRadiant  bool
	RadiantScore   int
	DireScore      int
	GameTime       int // seconds
	WinProbability float64
	NotablePlayers []string
	LastGame       *LastGameInfo
}

type LastGameInfo struct {
	HeroName  string
	Kills     int
	Deaths    int
	Assists   int
	Win       bool
	DurationS int
}

type MatchStartedMessage struct {
	ChannelID      string
	TwitchUserID   string
	SteamAccountID string
	HeroName       string
}

type MatchEndedMessage struct {
	ChannelID      string
	TwitchUserID   string
	SteamAccountID string
	Win            bool
	HeroName       string
	Mmr            int
	SessionWins    int
	SessionLosses  int
}

type RoshanKilledMessage struct {
	ChannelID    string
	TwitchUserID string
	Team         string // "radiant" | "dire"
	GameTime     int
}

type AegisPickupMessage struct {
	ChannelID    string
	TwitchUserID string
	PlayerName   string
	GameTime     int
}
```

- [ ] **Step 2: api overlay subject** (`libs/bus-core/api/dota.go`):

```go
package api

const DotaStateUpdateSubject = "api.dota.state_update"

type DotaStateUpdateMessage struct {
	ChannelID      string  `json:"channelId"`
	InGame         bool    `json:"inGame"`
	Mmr            int     `json:"mmr"`
	SessionWins    int     `json:"sessionWins"`
	SessionLosses  int     `json:"sessionLosses"`
	WinProbability float64 `json:"winProbability"`
	HeroName       string  `json:"heroName"`
	MatchID        int64   `json:"matchId"`
}
```

- [ ] **Step 3: Wire into bus** — `bus-services.go` add:

```go
type dotaBus struct {
	GetData       Queue[dota.GetDataRequest, dota.GetDataResponse]
	MatchStarted  Queue[dota.MatchStartedMessage, struct{}]
	MatchEnded    Queue[dota.MatchEndedMessage, struct{}]
	RoshanKilled  Queue[dota.RoshanKilledMessage, struct{}]
	AegisPickup   Queue[dota.AegisPickupMessage, struct{}]
}
```
`apiBus` gets `DotaStateUpdate Queue[api.DotaStateUpdateMessage, struct{}]`.
`bus.go`: add `Dota *dotaBus` field to `Bus` struct (near `bus.go:30-53`) and `NewNatsQueue` initializers in `NewNatsBus` (mirror the PredictionBegin block at `bus.go:414-419`, DotaStateUpdate mirrors TriggerKappagen at `bus.go:546-549`).

- [ ] **Step 4: Compile check**

Run: `cd libs/bus-core && go build ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add libs/bus-core
git commit -m "feat: dota bus subjects and queues"
```

---

### Task 4: `apps/dota` scaffolding

**Files:**
- Create: `apps/dota/cmd/main.go`
- Create: `apps/dota/app/app.go`
- Create: `apps/dota/go.mod`
- Create: `apps/dota/Dockerfile`
- Modify: `go.work` (add `./apps/dota`)
- Modify: `cli/internal/goapp/apps.go` (add `{Name: "dota", DebugPort: 2360}`)
- Modify: `cli/internal/cmds/deploy/deploy.go` (serviceImages + releaseServices)
- Modify: `.github/workflows/dockerv3.yml` (matrix entry)
- Modify: `docker-compose.stack.yml` (dota service, copy events block)
- Modify: `libs/config/config.go` (add dota env fields)

- [ ] **Step 1: config fields** (`libs/config/config.go`, near `EventsubHttpPort` at `:93`):

```go
	DotaHttpPort     int    `required:"false" default:"3031"  envconfig:"DOTA_HTTP_PORT"`
	DotaStratzToken  string `required:"false" envconfig:"DOTA_STRATZ_TOKEN"`
	DotaSteamAPIKey  string `required:"false" envconfig:"DOTA_STEAM_API_KEY"`
```

- [ ] **Step 2: app skeleton** — copy structure from `apps/events`: `cmd/main.go` (`fx.New(app.App).Run()`), `app/app.go` using `baseapp.CreateBaseApp(baseapp.Opts{AppName: "dota"})`. `go.mod`: copy module pattern from `apps/events/go.mod` (module `github.com/twirapp/twir/apps/dota`, same go version, same replace directives). Dockerfile: copy `apps/events/Dockerfile`, change binary name to `twir-dota`.

- [ ] **Step 3: registration** — edits in `go.work`, `cli/internal/goapp/apps.go`, `deploy.go`, `dockerv3.yml`, `docker-compose.stack.yml` (see header checklist).

- [ ] **Step 4: Build**

Run: `bun cli build app dota`
Expected: binary in `apps/dota/.out/twir-dota`

- [ ] **Step 5: Commit**

```bash
git add apps/dota go.work cli .github docker-compose.stack.yml libs/config
git commit -m "feat: scaffold apps/dota service"
```

---

### Task 5: GSI HTTP server

**Files:**
- Create: `apps/dota/internal/gsi/types.go` (GSI payload structs)
- Create: `apps/dota/internal/gsi/server.go` (HTTP mux, token auth, rate limit)
- Create: `apps/dota/internal/gsi/server_test.go`
- Create: `apps/dota/internal/gsi/testdata/*.json` (GSI fixtures)

- [ ] **Step 1: GSI payload types** — Valve GSI JSON schema: top-level `auth.token`, `provider`, `map` (matchid, game_state, game_time, radiant_score, dire_score, win_team), `player` (account_id, team_name, kills, deaths, assists), `hero` (name, id), `events` array, `previously`/`added`. Game states: `DOTA_GAMERULES_STATE_HERO_SELECTION`, `STRATEGY_TIME`, `PRE_GAME`, `GAME_IN_PROGRESS`, `POST_GAME`. Events include `roshan_killed` (killer_team), `aegis_picked_up` (player).

- [ ] **Step 2: HTTP server** — stdlib `net/http` ServeMux like `apps/eventsub/internal/http/server.go:19-31`. `POST /gsi/{token}`: lookup settings by token via repository (Task 2), 401 on miss, per-token sliding-window rate limit (e.g. 5 rps) in-memory, parse JSON, pass to state machine. `GET /health`. Wired into fx lifecycle with `Start()`/`Stop(ctx)`.

- [ ] **Step 3: Tests** — `server_test.go`: table tests with fixture JSON files: valid token+payload → 200; bad token → 401; malformed JSON → 400; rate limit exceeded → 429. Mock repository (interface from Task 2).

- [ ] **Step 4: Run tests**

Run: `cd apps/dota && go test ./internal/gsi/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add apps/dota/internal/gsi
git commit -m "feat: dota GSI http server"
```

---

### Task 6: Match state machine

**Files:**
- Create: `apps/dota/internal/match/state_machine.go`
- Create: `apps/dota/internal/match/state_machine_test.go`

- [ ] **Step 1: Failing tests** — cover: idle→in_game on `GAME_IN_PROGRESS`; ignore events when idle; single MatchStarted on reconnect (dedupe by matchid); post_game win detection via `map.win_team` + player team; MMR ±delta and session W/L increments on MatchEnded; roshan/aegis event emission from `added.events`.

- [ ] **Step 2: Run tests, verify FAIL** — `go test ./internal/match/...`

- [ ] **Step 3: Implement state machine** — per-channel struct keyed by channelID, transitions on `map.game_state`, emits domain events via a callback interface (injected; bus wiring in Task 8). Snapshot to Redis (`kv.KV`, key `cache:twir:dota:matchstate:<channelID>`, TTL 6h) so commands work after service restart mid-game.

- [ ] **Step 4: Run tests, verify PASS**

- [ ] **Step 5: Commit**

```bash
git add apps/dota/internal/match
git commit -m "feat: dota match state machine"
```

---

### Task 7: Stratz + OpenDota clients

**Files:**
- Create: `libs/integrations/stratz/client.go` (GraphQL: win probability for live match via `live` endpoints or match data; notable/pro players list)
- Create: `libs/integrations/opendota/client.go` (`GET /players/{account_id}/recentMatches`, `GET /players/{account_id}/heroes`, `GET /constants/heroes`)
- Create: `apps/dota/internal/stats/stats.go` (facade + Redis caching, TTL 30s for WP, 5min for notable players)

Note: Stratz GraphQL endpoint `https://api.stratz.com/graphql`, header `Authorization: Bearer <token>` (config `DotaStratzToken`). OpenDota is keyless with rate limits.

- [ ] **Step 1: OpenDota client** — recent matches + hero name constants (cache constants 24h in Redis).

- [ ] **Step 2: Stratz client** — win probability query; handle missing token (client disabled → WP = 0 / skip feature).

- [ ] **Step 3: Facade with kv caching** (`twirapp/kv`, pattern from `libs/cache/generic-cacher/db-generic-cacher.go`).

- [ ] **Step 4: Compile + unit tests for parsing (mock HTTP via httptest)**

Run: `cd apps/dota && go test ./internal/stats/... && go build ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add libs/integrations/stratz libs/integrations/opendota apps/dota/internal/stats
git commit -m "feat: stratz and opendota clients"
```

---

### Task 8: Bus wiring in apps/dota

**Files:**
- Create: `apps/dota/internal/buslistener/bus_listener.go`
- Modify: `apps/dota/app/app.go`

- [ ] **Step 1: GetData responder** — `bus.Dota.GetData.SubscribeGroup("dota", ...)`: load settings by channel, read state machine snapshot, if in game fetch WP from stats facade, else fetch last game from OpenDota; return `GetDataResponse`. Unlinked → `Linked: false`.

- [ ] **Step 2: Event publishers** — state machine callbacks (Task 6) publish `MatchStarted/MatchEnded/RoshanKilled/AegisPickup` and `bus.Api.DotaStateUpdate.Publish(...)` on every state change (throttle: only on meaningful change — game state, score, WP ±5%).

- [ ] **Step 3: Commit**

```bash
git add apps/dota
git commit -m "feat: dota bus listener and event publishers"
```

---

### Task 9: Chat events → bots messages

**Files:**
- Create: `apps/dota/internal/chatalerts/chat_alerts.go`

- [ ] **Step 1: Implement** — on MatchStarted/MatchEnded/RoshanKilled/AegisPickup: check `ChatEvents` settings (enabled, cooldown per event — track last-sent in Redis), template rendering (variables: `{hero}`, `{mmr}`, `{wins}`, `{losses}`, `{team}`, `{player}`, `{time}`), then `bus.Bots.SendMessage.Publish` (pattern from `apps/events/internal/chat_alerts/follow.go:105-115`, request type `bots.SendMessageRequest`).

- [ ] **Step 2: Compile + commit**

Run: `cd apps/dota && go build ./...`

```bash
git add apps/dota/internal/chatalerts
git commit -m "feat: dota chat alerts"
```

---

### Task 10: Parser commands

**Files:**
- Create: `apps/parser/internal/commands/dota2/mmr.go`, `wl.go`, `lg.go`, `gm.go`, `np.go`, `wp.go`, `dota2.go` (registration)
- Delete: `apps/parser/internal/commands/dota/` (commented-out legacy)
- Modify: wherever default commands are registered (find `predictions.Start` registration — `apps/parser/internal/commands/commands.go`)

- [ ] **Step 1: Commands** — each is a `types.DefaultCommand` (pattern: `apps/parser/internal/commands/games/russian_roulette.go`). Handler: read `model.ChannelsDotaSettings` via `parseCtx.Services.Gorm` (check enabled + per-command toggle) → `parseCtx.Services.Bus.Dota.GetData.Request(...)` → format response via i18n locales. `!mmr set N` subcommand restricted to broadcaster/mods (roles). Medal name from MMR: Herald <770, Guardian <1540, Crusader <2310, Archon <3080, Legend <3850, Ancient <4620, Divine <5420+, Immortal top (display "Immortal").

- [ ] **Step 2: i18n** — add locale entries under `apps/parser/locales` (en + ru, follow existing structure).

- [ ] **Step 3: Register commands** in the default commands list.

- [ ] **Step 4: Build + commit**

Run: `cd apps/parser && go build ./...`

```bash
git add apps/parser
git commit -m "feat: dota chat commands (!mmr !wl !lg !gm !np !wp)"
```

---

### Task 11: Auto Twitch predictions

**Files:**
- Create: `apps/dota/internal/predictions/predictions.go`

- [ ] **Step 1: Implement** — subscribe state machine callbacks: on `in_game` and `prediction_settings.enabled` → `twitch.NewUserClientWithContext(ctx, twitchUserID, config, bus)` (see `apps/parser/internal/commands/predictions/start.go:68-102`) → `CreatePrediction` with title from template + outcomes `Yes/No` (localized) + configured window. Store prediction ID + outcome IDs in Redis. On `post_game` → fetch prediction (`GetPredictions`), `EndPrediction` with winning outcome. On abandoned/error → `EndPrediction` status CANCELED. Dedupe by matchID (Redis key). Guards: only when streamer team known; skip if a prediction already active (API error 400 "already active" → ignore).

- [ ] **Step 2: Compile + commit**

Run: `cd apps/dota && go build ./...`

```bash
git add apps/dota/internal/predictions
git commit -m "feat: dota auto twitch predictions"
```

---

### Task 12: api-gql dota domain

**Files:**
- Create: `apps/api-gql/internal/delivery/gql/schema/dota/dota.graphql`
- Create: `apps/api-gql/internal/services/dota/service.go` (+ steam openid helper)
- Create: `apps/api-gql/internal/delivery/gql/mappers/dota.go`
- Create: `apps/api-gql/internal/di/dota.go`
- Modify: `apps/api-gql/cmd/main.go` (register DI module)
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/resolver.go` (Deps field)
- Generated: `resolvers/dota.resolver.go` via `bun cli build gql`

- [ ] **Step 1: Schema** — pattern from `schema/overlays/overlays-kappagen.graphql:1-13`:

```graphql
extend type Query {
	dotaSettings: DotaSettings! @isAuthenticated @hasChannelRolesDashboardPermission(permission: VIEW_GAMES)
	dotaGsiConfig: String! @isAuthenticated @hasChannelRolesDashboardPermission(permission: VIEW_GAMES)
	dotaSteamAuthLink: String! @isAuthenticated @hasChannelRolesDashboardPermission(permission: VIEW_GAMES)
}

extend type Mutation {
	dotaUpdateSettings(input: DotaSettingsInput!): DotaSettings! @isAuthenticated @hasChannelRolesDashboardPermission(permission: VIEW_GAMES)
	dotaRegenerateGsiToken: DotaSettings! @isAuthenticated @hasChannelRolesDashboardPermission(permission: VIEW_GAMES)
	dotaUnlinkSteam: Boolean! @isAuthenticated @hasChannelRolesDashboardPermission(permission: VIEW_GAMES)
	dotaResetSession: DotaSettings! @isAuthenticated @hasChannelRolesDashboardPermission(permission: VIEW_GAMES)
}

extend type Subscription {
	dotaState(apiKey: String!): DotaState!
}
```
Types: `DotaSettings` (all entity fields), `DotaSettingsInput`, `DotaState` (mirrors `DotaStateUpdateMessage`).

- [ ] **Step 2: Service** — `GetOrCreate` by channel (default settings), `Update`, `RegenerateToken`, `ResetSession`, `GsiConfig(channelID)` (renders cfg file text with `https://<host>/dota/gsi/<token>`... note: GSI endpoint is on apps/dota port — expose via public URL config or route through api-gql; simplest: cfg points to a dedicated public endpoint `DOTA_GSI_PUBLIC_URL` config value), Steam OpenID: `GetAuthLink` (openid mode=checkid_setup, realm = site URL, return_to = dashboard page `/dashboard/games/dota`), `SteamCallback(ctx, channelID, queryParams)` — verify `openid.mode=check_authentication` POST to `https://steamcommunity.com/openid/login`, parse steamid from `claimed_id`, convert SteamID64 → accountID (subtract 76561197960265728), save. Subscription bridging: subscribe `bus.Api.DotaStateUpdate` → `wsRouter.Publish("api.dota.state."+channelID, msg)` (kappagen pattern `apps/api-gql/internal/services/overlays/kappagen/service.go:37-45`).

- [ ] **Step 3: Regenerate resolvers**

Run: `bun cli build gql`
Then implement resolver bodies (pattern: `resolvers/overlays-kappagen.resolver.go`).

- [ ] **Step 4: Build + commit**

Run: `cd apps/api-gql && go build ./...`

```bash
git add apps/api-gql
git commit -m "feat: dota graphql api"
```

---

### Task 13: Dashboard module page (web/layers/dashboard)

**Files:**
- Create: `web/layers/dashboard/pages/dashboard/dota.vue`
- Create: `web/layers/dashboard/features/dota/dota.vue` (+ `dota-form-schema.ts`, steam-link card, gsi-config card)
- Create: `web/layers/dashboard/api/dota/dota.ts` (urql composable)
- Modify: `web/layers/dashboard/config/navigation.ts` (nav item)
- Modify: `web/layers/dashboard/locales/{en,ru,...}.json` (keys; en required, others fallback)

- [ ] **Step 1: API composable** — `createGlobalState` + `useQuery`/`useMutation` with `graphql()` codegen tag (pattern: `web/layers/dashboard/api/games/games.ts`). Run graphql codegen for dashboard if required (check `web/layers/dashboard` package.json scripts — gql.tada or graphql-codegen; follow how games.ts documents are typed).

- [ ] **Step 2: Page + feature** — sections: enable toggle; Steam link (button opens auth link in new tab; Steam OpenID returns to this page with query params → call mutation); GSI config download (fetch `dotaGsiConfig` string, download as `gamestate_integration_twir.cfg`, show install instructions); MMR input + delta; session W/L display + reset; commands toggles; predictions settings; chat event templates. Forms: vee-validate + zod (pattern: `features/overlays/kappagen/kappagen-form-schema.ts`).

- [ ] **Step 3: Navigation + locales.**

- [ ] **Step 4: Typecheck + lint**

Run: `cd web && bun run typecheck` (or per AGENTS.md command; verify with package.json)
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add web/layers/dashboard
git commit -m "feat: dota dashboard settings page"
```

---

### Task 14: Overlay pages (frontend/overlays)

**Files:**
- Create: `frontend/overlays/src/pages/overlays/dota-medal.vue`, `dota-wl.vue`, `dota-wp.vue`
- Create: `frontend/overlays/src/composables/dota/use-dota-state.ts`
- Modify: `frontend/overlays/src/plugins/router.ts` (3 routes)

- [ ] **Step 1: Composable** — urql `useSubscription` to `dotaState(apiKey: $apiKey)` (pattern: `composables/kappagen/use-kappagen-socket.ts:48-114`), `connect(apiKey)`/`destroy()`.

- [ ] **Step 2: Pages** — minimal styled displays: medal (name + MMR), WL (wins/losses/session), WP (percent bar). Transparent background, OBS-ready. Routes: `/:apiKey/dota/medal`, `/:apiKey/dota/wl`, `/:apiKey/dota/wp`.

- [ ] **Step 3: Build**

Run: `cd frontend/overlays && bun run build`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add frontend/overlays
git commit -m "feat: dota overlays (medal, wl, wp)"
```

---

### Task 15: Final verification

- [ ] **Step 1:** `bun cli build` (all Go apps) — PASS
- [ ] **Step 2:** `bun lint` — PASS
- [ ] **Step 3:** `go test ./...` in apps/dota — PASS
- [ ] **Step 4:** web typecheck — PASS
- [ ] **Step 5:** Update `docs/superpowers/specs/2026-07-20-dota-module-design.md` status → implemented (amend if deviations).

---

## Explicitly out of scope (v1)

go-dota2 GC client (spec lists as optional fallback — requires Steam credentials + 2FA handling; defer), OBS scene switching, minimap overlay, hero bets, smurf detection, multiple accounts per channel, `libs/bus-core/src` TS mirror for dota subjects (no JS consumers in v1).
