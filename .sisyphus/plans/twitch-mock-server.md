# Twitch Mock Server for Local Development

## TL;DR

> **Quick Summary**: Build a new `apps/twitch-mock/` Go service that fully emulates the Twitch HTTP API (`id.twitch.tv` + `api.twitch.tv/helix/*`) and EventSub WebSocket server, with minimal surgical changes to 7 existing services behind a `TWITCH_MOCK_ENABLED` flag so the entire dev stack works without a real Twitch account.
>
> **Deliverables**:
>
> - New Go service `apps/twitch-mock/` (HTTP :7777, WebSocket :8081, Admin UI :3333)
> - Config changes in `libs/config/config.go`
> - Mock redirect in `libs/twitch/twitch.go` (shared `mockRoundTripper`)
> - Redirects in `apps/tokens`, `apps/eventsub` to use mock URLs
> - `apps/api-gql` `authLink` resolver — returns mock OAuth URL when enabled (NEW)
> - `libs/migrations/seeds/defaultBot.go` — configurable validate URL
> - `docker-compose.dev.yml` — new `twitch-mock` service with healthcheck
> - `.env.mock.example` + `apps/twitch-mock/README.md`
> - Unit tests + E2E verification via MCP Playwright
>
> **Estimated Effort**: Large
> **Parallel Execution**: YES — 4 waves
> **Critical Path**: Task 1 (Config) → Task 2 (mock server) → Task 11 (authLink) → Task 10 (E2E)

---

## Context

### Original Request

Разработчик перманентно забанен на Твиче и не может залогиниться. Весь стек (`apps/api-gql`, `apps/tokens`, `apps/bots`, `apps/eventsub`) завязан на реальный Twitch OAuth и Helix API. Нужен полноценный dev-стенд — собственный мок-сервер, эмулирующий все API Твича, чтобы разрабатывать и тестировать бота без реального аккаунта. Одна команда должна поднять всё.

### Interview Summary

**Key Discussions**:

- Цель: войти в дашборд + тестировать логику бота + EventSub события + полноценный dev-стенд
- Уровень мокинга: собственный HTTP/WS сервер (не флаги в коде)
- База данных: полностью автоматизированный seed, одна команда

**Research Findings** (от 4 агентов):

- `helix.AuthBaseURL` в библиотеке `nicklaw5/helix/v2` — **hardcoded package constant**, НЕ перекрывается `APIBaseURL`. Нужен кастомный `http.RoundTripper` который перезаписывает hostname в запросах.
- `kvizyx/twitchy/eventsub` — поддерживает кастомный WS URL через `eventsub.WebsocketWithServerURL(url)`.
- `apps/eventsub/internal/manager/` делает прямые `http.DefaultClient.Do()` вызовы (не через helix) — нужно передать `httpClient *http.Client` в `Manager` struct.
- `scs` Redis сессия использует gob-кодирование — seed должен создать сессию через реальный OAuth flow мока, не вручную.
- `apps/eventsub/internal/manager/on_start.go` содержит **3** hardcoded `api.twitch.tv` URL (не 1), плюс ещё 1 в `subscribe.go`.
- `frontend/dashboard/src/plugins/router.ts` делает GraphQL `authenticatedUser` запрос — если null, редирект на `/`. Нужно чтобы весь OAuth flow работал.
- `web/layers/landing/pages/login.client.vue` открывает OAuth URL в браузере — этот URL должен указывать на мок.

### Metis Review

**Identified Gaps** (resolved):

- `helix.AuthBaseURL` не перекрывается: **resolved** — используем `mockRoundTripper` в `http.Client`.
- Нет критерия приёмки для каждой задачи: **resolved** — добавлены curl/grep команды для каждой задачи.
- `http.DefaultClient` в eventsub: **resolved** — добавляем `httpClient *http.Client` в `Manager` struct.
- Redis сессия для seed: **resolved** — seed делает полный OAuth flow через мок, не создаёт сессию вручную.
- `defaultBot.go` хардкодит `id.twitch.tv`: **resolved** — configurable URL via config.
- Frontend OAuth URL (браузер → мок): **resolved** — `TWITCH_MOCK_AUTH_URL` конфигурирует URL для авторизации в frontend.
- `BOT_ACCESS_TOKEN` / `BOT_REFRESH_TOKEN` обязательны в config: **resolved** — mock values в `.env.mock.example`.
- Нет Docker healthcheck: **resolved** — добавляем `/health` endpoint + `healthcheck` в compose.
- `twitch-mock` в production: **resolved** — ТОЛЬКО в `docker-compose.dev.yml` с `# DEV ONLY` комментарием.
- **[Momus]** `authLink` GraphQL resolver возвращает реальный Twitch URL: **resolved** — Task 11 модифицирует resolver для возврата mock OAuth URL когда `TWITCH_MOCK_ENABLED`.
- **[Momus]** Task 10 использует `@playwright/test` которого нет в репо: **resolved** — Task 10 переписан для использования MCP Playwright (без repo-hosted test файла).
- **[Momus]** WebSocket `event` поле вложено в `payload`: **resolved** — исправлен JSON в Task 2 (`event` теперь на верхнем уровне, как в `kvizyx/twitchy` structs).

---

## Work Objectives

### Core Objective

Создать полноценный Twitch mock server, который позволяет запустить весь dev-стек без реального аккаунта Твич — одной командой `docker compose -f docker-compose.dev.yml up -d && bun dev`.

### Concrete Deliverables

- `apps/twitch-mock/` — новый Go сервис (HTTP + WebSocket + Admin UI + Dockerfile)
- `libs/config/config.go` — 4 новых поля `TwitchMock*`
- `libs/twitch/twitch.go` + `libs/twitch/mock_roundtripper.go` — shared `mockRoundTripper`
- `apps/tokens/internal/bus_listener/bus_listener.go` — redirect to mock
- `apps/eventsub/internal/manager/*.go` — configurable httpClient + URLs
- `apps/api-gql/internal/delivery/gql/resolvers/user.resolver.go` — mock `authLink` URL (NEW)
- `libs/migrations/seeds/defaultBot.go` — configurable validate URL
- `docker-compose.dev.yml` — `twitch-mock` service с healthcheck
- `.env.mock.example` — все mock переменные
- `apps/twitch-mock/README.md` — документация
- `apps/twitch-mock/*_test.go` — unit tests
- E2E verification via MCP Playwright (no test file written to disk)

### Definition of Done

- [ ] `curl http://localhost:7777/health` → 200
- [ ] `docker compose -f docker-compose.dev.yml up -d && bun dev` — все сервисы запускаются без ошибок
- [ ] Eventsub логи содержат "conduit ensured" + "websocket welcome received"
- [ ] Браузер: `http://localhost:3005` → Login → Dashboard → виден mock пользователь "MockStreamer"
- [ ] `go test ./apps/twitch-mock/...` — все тесты проходят
- [ ] E2E MCP Playwright сценарий — все шаги проходят, скриншот сохранён

### Must Have

- Мок сервер эмулирует OAuth (`/oauth2/token`, `/oauth2/validate`, `/oauth2/authorize`)
- Мок сервер эмулирует минимальный Helix API (users, streams, channels, eventsub/\*, chat/messages, moderation/bans + все остальные возвращают `{"data":[]}`)
- EventSub WebSocket сервер с правильным протоколом (session_welcome, keepalive, notifications)
- Admin UI для ручного тригегринга событий (plain HTML, embedded в Go binary)
- `TWITCH_MOCK_ENABLED=false` — НУЛЕВЫЕ изменения поведения в production
- Seed автоматический через OAuth flow мока (не вручную)
- Docker healthcheck + `depends_on: condition: service_healthy`
- Idempotent seed (check-before-insert)
- Токены возвращаются с `expires_in: 99999999` (никогда не истекают в dev)
- Все mock изменения за `if config.TwitchMockEnabled` гардом

### Must NOT Have (Guardrails)

- **НЕ** мутировать `helix.AuthBaseURL` package global — использовать RoundTripper
- **НЕ** добавлять `twitch-mock` в `docker-compose.stack.yml` (production)
- **НЕ** создавать stateful Helix эндпоинты (кроме conduits + subscriptions) — только stub 200 OK
- **НЕ** использовать frontend framework для Admin UI — только embedded plain HTML в Go binary
- **НЕ** хардкодить порты в application code — только из config/env
- **НЕ** помещать test helpers в production code (`libs/`, shared packages)
- **НЕ** затрагивать `docker-compose.stack.yml`
- **НЕ** менять поведение когда `TWITCH_MOCK_ENABLED=false`

---

## Verification Strategy

> **ZERO HUMAN INTERVENTION** - ALL verification is agent-executed.

### Test Decision

- **Infrastructure exists**: YES (Go tests в проекте)
- **Automated tests**: Tests-after (unit tests для mock сервера + E2E)
- **Framework**: Go built-in `testing` + Playwright для E2E

### QA Policy

Каждая задача имеет curl/grep/go test команды для верификации. Все E2E тесты через Playwright.

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately — no dependencies):
├── Task 1: Config fields in libs/config/config.go [quick]
└── Task 2: Build apps/twitch-mock/ service [deep]

Wave 2 (After Wave 1):
├── Task 3:  libs/twitch mockRoundTripper + APIBaseURL redirect [quick]
├── Task 4:  apps/tokens redirect to mock [quick]
├── Task 5:  apps/eventsub redirect (httpClient + configurable URLs + WS URL) [unspecified-high]
├── Task 6:  defaultBot.go seed configurable validate URL [quick]
├── Task 7:  docker-compose.dev.yml add twitch-mock service [quick]
├── Task 8:  .env.mock.example + README.md [writing]
└── Task 11: apps/api-gql authLink resolver — return mock OAuth URL [quick]

Wave 3 (After Wave 2):
└── Task 9: Unit tests for apps/twitch-mock/ [unspecified-high]

Wave 4 (After Wave 3 + Task 11):
└── Task 10: E2E MCP Playwright verification (full login flow) [unspecified-high + playwright]

Wave FINAL (After ALL tasks — 4 parallel reviews):
├── Task F1: Plan compliance audit [oracle]
├── Task F2: Code quality review [unspecified-high]
├── Task F3: Real manual QA [unspecified-high + playwright]
└── Task F4: Scope fidelity check [deep]
→ Present results → Get explicit user okay

Critical Path: Task 1 → Task 2 → Task 11 → Task 10 → F1-F4
Parallel Speedup: ~65% faster than sequential
```

### Dependency Matrix

| Task | Depends On        | Blocks                |
| ---- | ----------------- | --------------------- |
| 1    | —                 | 3, 4, 5, 6, 7, 8, 11  |
| 2    | —                 | 3, 5, 6, 7, 9, 10, 11 |
| 3    | 1                 | 4, 10                 |
| 4    | 1, 3              | 10                    |
| 5    | 1, 2              | 10                    |
| 6    | 1, 2              | 10                    |
| 7    | 2                 | 10                    |
| 8    | 1                 | —                     |
| 9    | 2                 | —                     |
| 10   | 3, 4, 5, 6, 7, 11 | F1-F4                 |
| 11   | 1, 2              | 10                    |

### Agent Dispatch Summary

- **Wave 1**: 2 tasks — T1 → `quick`, T2 → `deep`
- **Wave 2**: 7 tasks — T3→`quick`, T4→`quick`, T5→`unspecified-high`, T6→`quick`, T7→`quick`, T8→`writing`, T11→`quick`
- **Wave 3**: 1 task — T9 → `unspecified-high`
- **Wave 4**: 1 task — T10 → `unspecified-high` (+ `playwright` skill)
- **FINAL**: 4 tasks — F1→`oracle`, F2→`unspecified-high`, F3→`unspecified-high`, F4→`deep`

---

## TODOs

---

- [x] 1. Add Mock Config Fields to `libs/config/config.go`

  **What to do**:
  - Open `libs/config/config.go` and add 4 new fields to the `Config` struct:
    ```go
    TwitchMockEnabled bool   `env:"TWITCH_MOCK_ENABLED"`
    TwitchMockApiUrl  string `env:"TWITCH_MOCK_API_URL"  envDefault:"http://twitch-mock:7777"`
    TwitchMockAuthUrl string `env:"TWITCH_MOCK_AUTH_URL" envDefault:"http://twitch-mock:7777"`
    TwitchMockWsUrl   string `env:"TWITCH_MOCK_WS_URL"   envDefault:"ws://twitch-mock:8081/ws"`
    ```
  - The fields use `caarlos0/env` tags (same pattern as existing fields in the file)
  - `TwitchMockEnabled` defaults to `false`
  - Default URLs point to Docker service name `twitch-mock` (for use inside Docker network)

  **Must NOT do**:
  - Do NOT change any existing fields
  - Do NOT add logic — just struct fields

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Single file, struct field additions only
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `supabase-postgres-best-practices`: no DB work

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 2)
  - **Blocks**: Tasks 3, 4, 5, 6, 7, 8
  - **Blocked By**: None

  **References**:
  - `libs/config/config.go` — read entire file, follow existing field pattern (env tags, naming convention)

  **Acceptance Criteria**:

  ```
  Scenario: Config fields added correctly
    Tool: Bash
    Steps:
      1. grep -c "TwitchMock" libs/config/config.go
    Expected Result: output is "4"
    Evidence: .sisyphus/evidence/task-1-config-fields.txt

  Scenario: Config compiles
    Tool: Bash
    Steps:
      1. go build ./libs/config/...
    Expected Result: exit code 0, no errors
    Evidence: .sisyphus/evidence/task-1-build.txt
  ```

  **Commit**: YES (group alone)
  - Message: `feat(config): add TWITCH_MOCK_* config fields`
  - Files: `libs/config/config.go`

---

- [x] 2. Build `apps/twitch-mock/` — Full Mock Server

  **What to do**:
  Create the entire `apps/twitch-mock/` Go service. Use the same project conventions: `fx` DI, context propagation, structured logging via `libs/logger`. Structure:

  ```
  apps/twitch-mock/
  ├── cmd/main.go           # fx.New() entrypoint
  ├── app/app.go            # fx module registrations
  ├── internal/
  │   ├── config/config.go  # local config (mock port, fake user IDs)
  │   ├── state/state.go    # in-memory state (conduits, subscriptions, stream state)
  │   ├── handlers/
  │   │   ├── auth.go       # POST /oauth2/token, GET /oauth2/authorize, GET /oauth2/validate
  │   │   ├── helix.go      # all /helix/* stub endpoints
  │   │   ├── conduits.go   # GET/POST /helix/eventsub/conduits, PATCH /helix/eventsub/conduits/shards
  │   │   ├── subscriptions.go # POST/DELETE /helix/eventsub/subscriptions
  │   │   └── health.go     # GET /health
  │   ├── websocket/
  │   │   └── eventsub.go   # WS server :8081, session_welcome/keepalive/notification
  │   └── admin/
  │       ├── admin.go      # GET /admin (embedded HTML), POST /admin/trigger/{event}
  │       └── templates/
  │           └── index.html  # plain HTML admin panel
  └── Dockerfile
  ```

  **Fake user constants** (hardcoded, stable):

  ```go
  const (
      MockBroadcasterID    = "12345"
      MockBroadcasterLogin = "mockstreamer"
      MockBroadcasterName  = "MockStreamer"
      MockBotID            = "67890"
      MockBotLogin         = "mockbot"
      MockBotName          = "MockBot"
      MockAppToken         = "mock-app-token"
      MockBotToken         = "mock-bot-token"
      MockUserToken        = "mock-user-token"
  )
  ```

  **HTTP Server (port 7777)** — handles BOTH `id.twitch.tv` and `api.twitch.tv` paths on same port:

  Auth endpoints (`/oauth2/*`):
  - `POST /oauth2/token` — returns `{"access_token":"mock-app-token","token_type":"bearer","expires_in":99999999,"scope":""}` for `client_credentials`; for `authorization_code` grant returns user token; for `refresh_token` returns refreshed token
  - `GET /oauth2/authorize` — redirects to `{SITE_BASE_URL}/login?code=mock_code&state={state}` (state preserved from query param)
  - `GET /oauth2/validate` — returns `{"client_id":"mock","login":"mockstreamer","user_id":"12345","expires_in":99999999}` for bearer token; `{"login":"mockbot","user_id":"67890","expires_in":99999999}` for `mock-bot-token`

  Helix stubs (`/helix/*`) — ALL return 200 with `{"data":[]}` EXCEPT:
  - `GET /helix/users` → returns mock user array (broadcaster or bot based on token/params)
  - `GET /helix/streams` → returns mock stream (online/offline based on in-memory state)
  - `GET /helix/channels` → returns mock channel info
  - `GET /helix/eventsub/conduits` → returns in-memory conduits
  - `POST /helix/eventsub/conduits` → creates conduit, stores in memory
  - `PATCH /helix/eventsub/conduits/shards` → updates conduit shard session_id
  - `POST /helix/eventsub/subscriptions` → stores subscription, returns 202
  - `DELETE /helix/eventsub/subscriptions/{id}` → removes subscription, returns 204
  - `POST /helix/chat/messages` → logs message, returns `{"data":[{"message_id":"mock-msg-id","is_sent":true}]}`
  - `POST /helix/moderation/bans` → returns `{"data":[{"broadcaster_id":"12345","moderator_id":"67890","user_id":"...","end_time":null}]}`
  - All other `/helix/*` → `{"data":[], "total": 0, "pagination": {}}`

  **WebSocket Server (port 8081, path `/ws`)**:
  - On connect: immediately send `session_welcome` with `session.id = UUID`, `keepalive_timeout_seconds: 30`
  - Every 25s: send `session_keepalive`
  - Support multiple simultaneous WS connections (multiple shards)
  - Expose internal channel `chan EventPayload` for Admin UI to send events
  - When event received on channel: broadcast `notification` message to all connected clients

  **Admin UI (port 3333)**:
  - `GET /admin` → serves embedded `index.html` (HTML forms for each event type)
  - `POST /admin/trigger/{event_type}` → parses JSON body, sends to WS broadcast channel
  - Event types supported: `channel.follow`, `channel.subscribe`, `channel.raid`, `stream.online`, `stream.offline`, `channel.chat.message`, `channel.ban`, `channel.unban`, `channel.cheer`, `channel.prediction.begin`, `channel.prediction.end`, `channel.poll.begin`, `channel.poll.end`, `channel.channel_points_custom_reward_redemption.add`

  **WebSocket notification payload format** — must match `kvizyx/twitchy/eventsub` library structs exactly.

  The library defines `WebsocketNotificationMessage` as:

  ```go
  type WebsocketNotificationMessage[Event any, C Condition] struct {
      WebsocketMessage[WebsocketNotificationMetadata, WebsocketNotificationPayload[C]]
      Event Event `json:"event"`
  }
  ```

  This means `event` is a **top-level field** (not inside `payload`). The correct JSON shape is:

  ```json
  {
  	"metadata": {
  		"message_id": "<uuid>",
  		"message_type": "notification",
  		"message_timestamp": "<RFC3339>",
  		"subscription_type": "channel.follow",
  		"subscription_version": "2"
  	},
  	"payload": {
  		"subscription": {
  			"id": "<uuid>",
  			"status": "enabled",
  			"type": "channel.follow",
  			"version": "2",
  			"cost": 0,
  			"condition": { "broadcaster_user_id": "12345", "moderator_user_id": "67890" },
  			"created_at": "<RFC3339>",
  			"transport": { "method": "websocket", "session_id": "<session_id>" }
  		}
  	},
  	"event": {
  		/* event-specific fields — top level, NOT nested in payload */
  	}
  }
  ```

  **CRITICAL**: Do NOT put `event` inside `payload`. It must be at the top level of the JSON object alongside `metadata` and `payload`. This is how `kvizyx/twitchy/eventsub` parses notifications via its embedded struct pattern.

  **Dockerfile** — multi-stage build, copy binary only, minimal image.

  **SITE_BASE_URL**: Read from env var `SITE_BASE_URL` (same as rest of project). Used in OAuth authorize redirect.

  **Must NOT do**:
  - Do NOT implement actual business logic — mock only
  - Do NOT persist to DB (except conduits/subscriptions in-memory)
  - Do NOT use a frontend framework for Admin UI — plain HTML with `<form>` elements only
  - Do NOT hardcode ports in Go code — use env vars or config

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Complex new service with multiple components (HTTP, WS, Admin), Go protocol implementation, must match exact library struct shapes
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: plain HTML only, no design system
    - `supabase-postgres-best-practices`: no DB work

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Task 1)
  - **Blocks**: Tasks 3, 5, 6, 7, 9, 10
  - **Blocked By**: None

  **References**:
  - `apps/eventsub/internal/manager/websocket.go` — study the event handler registration to understand what WS message shapes twitchy expects
  - `apps/eventsub/internal/mappers/chat_message.go` — study event payload shapes
  - `apps/api-gql/cmd/main.go` — follow fx entrypoint pattern
  - `libs/logger/` — use project logger
  - `github.com/kvizyx/twitchy/eventsub` — check `WebsocketWelcomeMessage` and `WebsocketNotificationMetadata` struct definitions to match exact JSON shapes

  **Acceptance Criteria**:

  ```
  Scenario: OAuth token endpoint returns correct shape
    Tool: Bash (curl)
    Steps:
      1. cd apps/twitch-mock && go run ./cmd/main.go &
      2. curl -s -X POST http://localhost:7777/oauth2/token -d "grant_type=client_credentials&client_id=mock&client_secret=mock"
    Expected Result: JSON with access_token="mock-app-token", expires_in=99999999
    Evidence: .sisyphus/evidence/task-2-oauth-token.json

  Scenario: /oauth2/validate returns mock user
    Tool: Bash (curl)
    Steps:
      1. curl -s -H "Authorization: OAuth mock-bot-token" http://localhost:7777/oauth2/validate
    Expected Result: JSON with login="mockbot", user_id="67890", expires_in=99999999
    Evidence: .sisyphus/evidence/task-2-validate.json

  Scenario: /helix/users returns mock user
    Tool: Bash (curl)
    Steps:
      1. curl -s -H "Authorization: Bearer mock-user-token" -H "Client-Id: mock" "http://localhost:7777/helix/users"
    Expected Result: JSON with data[0].login="mockstreamer", data[0].id="12345"
    Evidence: .sisyphus/evidence/task-2-helix-users.json

  Scenario: EventSub WebSocket welcome
    Tool: Bash
    Steps:
      1. Install wscat: npm install -g wscat
      2. Timeout 3 wscat -c ws://localhost:8081/ws 2>&1 | head -5
    Expected Result: output contains "session_welcome" and "session"
    Evidence: .sisyphus/evidence/task-2-ws-welcome.txt

  Scenario: Conduit CRUD
    Tool: Bash (curl)
    Steps:
      1. curl -s -X POST http://localhost:7777/helix/eventsub/conduits -H "Authorization: Bearer mock-app-token" -H "Client-Id: mock" -H "Content-Type: application/json" -d '{"shard_count":1}'
      2. curl -s http://localhost:7777/helix/eventsub/conduits -H "Authorization: Bearer mock-app-token" -H "Client-Id: mock"
    Expected Result: Step 1 returns conduit with id; Step 2 returns data array containing that conduit
    Evidence: .sisyphus/evidence/task-2-conduits.json

  Scenario: Build succeeds
    Tool: Bash
    Steps:
      1. go build ./apps/twitch-mock/...
    Expected Result: exit code 0
    Evidence: .sisyphus/evidence/task-2-build.txt
  ```

  **Commit**: YES (group alone — large new service)
  - Message: `feat(twitch-mock): add Twitch API mock server (HTTP+WS+AdminUI)`
  - Files: `apps/twitch-mock/`

---

- [x] 3. Add `mockRoundTripper` to `libs/twitch/` + Redirect helix clients

  **What to do**:
  Create `libs/twitch/mock_roundtripper.go` with a `MockRoundTripper` that:

  ```go
  type MockRoundTripper struct {
      Base      http.RoundTripper
      ApiUrl    string  // replaces https://api.twitch.tv
      AuthUrl   string  // replaces https://id.twitch.tv
  }
  func (t *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
      req = req.Clone(req.Context())
      if strings.Contains(req.URL.Host, "id.twitch.tv") {
          // rewrite host+scheme to AuthUrl
      }
      if strings.Contains(req.URL.Host, "api.twitch.tv") {
          // rewrite host+scheme to ApiUrl
      }
      return t.Base.RoundTrip(req)
  }
  ```

  In `libs/twitch/twitch.go`, in ALL THREE factory functions (`NewAppClientWithContext`, `NewUserClientWithContext`, `NewBotClientWithContext`), when `cfg.TwitchMockEnabled == true`:
  - Add `APIBaseURL: cfg.TwitchMockApiUrl` to `helix.Options{}`
  - Wrap `createHttpClient(ctx)` result with `MockRoundTripper{Base: existing, ApiUrl: cfg.TwitchMockApiUrl, AuthUrl: cfg.TwitchMockAuthUrl}`
  - Set `helix.Options{ HTTPClient: wrappedClient }`

  The `MockRoundTripper` must parse the mock URL to extract scheme+host and rewrite the request URL scheme+host accordingly while preserving path+query.

  **Must NOT do**:
  - Do NOT mutate `helix.AuthBaseURL` package global
  - Do NOT change behavior when `cfg.TwitchMockEnabled == false`

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Two files, focused change, established pattern
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 4, 5, 6, 7, 8)
  - **Blocks**: Task 4, Task 10
  - **Blocked By**: Task 1

  **References**:
  - `libs/twitch/twitch.go` — read full file, understand `createHttpClient()` and `helix.Options{}`
  - `libs/twitch/http_client.go` — understand existing RoundTripper wrapping pattern
  - `libs/config/config.go` — `TwitchMockEnabled`, `TwitchMockApiUrl`, `TwitchMockAuthUrl` fields

  **Acceptance Criteria**:

  ```
  Scenario: Build succeeds
    Tool: Bash
    Steps:
      1. go build ./libs/twitch/...
    Expected Result: exit code 0
    Evidence: .sisyphus/evidence/task-3-build.txt

  Scenario: mockRoundTripper rewrites id.twitch.tv
    Tool: Bash
    Steps:
      1. go test ./libs/twitch/... -run TestMockRoundTripper -v
    Expected Result: PASS — test verifies that requests to id.twitch.tv are rewritten to mock URL
    Evidence: .sisyphus/evidence/task-3-roundtripper-test.txt
  ```

  **Commit**: YES
  - Message: `feat(libs/twitch): add mockRoundTripper, redirect helix to mock when enabled`
  - Files: `libs/twitch/mock_roundtripper.go`, `libs/twitch/twitch.go`

---

- [x] 4. Redirect `apps/tokens` service to mock

  **What to do**:
  In `apps/tokens/internal/bus_listener/bus_listener.go`:
  - Import `libs/twitch`
  - When `config.TwitchMockEnabled`, add `APIBaseURL: config.TwitchMockApiUrl` to the `helix.Options{}` in `NewTokens()` (line ~74-84)
  - Wrap the HTTP client with `twitch.NewMockRoundTripper(cfg)` (the shared helper from Task 3) when mock enabled
  - The `globalClient` created in `NewTokens` makes `RequestAppAccessToken` on startup — this call must hit the mock

  **Must NOT do**:
  - Do NOT change token storage logic
  - Do NOT change refresh logic
  - Do NOT change behavior when `TwitchMockEnabled == false`

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Single file, minimal targeted change
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 3, 5, 6, 7, 8)
  - **Blocks**: Task 10
  - **Blocked By**: Tasks 1, 3

  **References**:
  - `apps/tokens/internal/bus_listener/bus_listener.go` — read full file, understand `NewTokens()` helix.Options construction
  - `libs/twitch/mock_roundtripper.go` — use `NewMockRoundTripper(cfg)` function (from Task 3)
  - `libs/config/config.go` — `TwitchMockEnabled`, `TwitchMockApiUrl`, `TwitchMockAuthUrl`

  **Acceptance Criteria**:

  ```
  Scenario: Build succeeds
    Tool: Bash
    Steps:
      1. go build ./apps/tokens/...
    Expected Result: exit code 0
    Evidence: .sisyphus/evidence/task-4-build.txt

  Scenario: Tokens service hits mock on startup (integration)
    Tool: Bash
    Steps:
      1. Start mock server: docker compose -f docker-compose.dev.yml up twitch-mock -d
      2. TWITCH_MOCK_ENABLED=true TWITCH_MOCK_API_URL=http://localhost:7777 TWITCH_MOCK_AUTH_URL=http://localhost:7777 go run ./apps/tokens/cmd/main.go &
      3. sleep 3
      4. grep "POST /oauth2/token" <mock server log>
    Expected Result: mock server log shows POST /oauth2/token request from tokens service
    Evidence: .sisyphus/evidence/task-4-tokens-mock-hit.txt
  ```

  **Commit**: YES
  - Message: `feat(tokens): redirect helix client to mock when TWITCH_MOCK_ENABLED`
  - Files: `apps/tokens/internal/bus_listener/bus_listener.go`

---

- [x] 5. Redirect `apps/eventsub` service to mock

  **What to do**:
  This is the most invasive change. 4 files in `apps/eventsub/internal/manager/`:

  **`manager.go`**:
  - Add `httpClient *http.Client` field to `Manager` struct
  - In `NewManager()`: when `config.TwitchMockEnabled`, create `httpClient` with `MockRoundTripper`; otherwise use `http.DefaultClient`
  - When `config.TwitchMockEnabled`, pass `eventsub.WebsocketWithServerURL(config.TwitchMockWsUrl)` as option to `eventsub.New()` (currently `eventsub.New()` is called with no options)
  - Add `apiBaseUrl string` field to `Manager` struct, set to `config.TwitchMockApiUrl` if enabled, else `"https://api.twitch.tv"`

  **`on_start.go`** (3 hardcoded URLs):
  - Line ~84: `http.NewRequest("GET", "https://api.twitch.tv/helix/eventsub/conduits", nil)` → replace `"https://api.twitch.tv"` with `c.apiBaseUrl`
  - Line ~175: `http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/conduits", ...)` → same
  - Line ~256: `http.NewRequest("PATCH", "https://api.twitch.tv/helix/eventsub/conduits/shards", ...)` → same
  - All 3 `http.DefaultClient.Do(req)` → `c.httpClient.Do(req)`

  **`subscribe.go`** (1 hardcoded URL):
  - Line ~116: `http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/subscriptions", ...)` → `c.apiBaseUrl + "/helix/eventsub/subscriptions"`
  - `http.DefaultClient.Do(req)` → `c.httpClient.Do(req)`

  **`unsubscribe.go`**: uses helix client (not DefaultClient) — already handled by Task 3's RoundTripper, no changes needed.

  **Must NOT do**:
  - Do NOT change any event handling logic
  - Do NOT change subscription type registrations
  - Do NOT use `http.DefaultClient` directly anywhere in manager after this change
  - Do NOT break the existing reconnect loop

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: 4 files, careful refactoring, must not break existing behavior
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 3, 4, 6, 7, 8)
  - **Blocks**: Task 10
  - **Blocked By**: Tasks 1, 2

  **References**:
  - `apps/eventsub/internal/manager/manager.go` — read full, understand Manager struct and NewManager
  - `apps/eventsub/internal/manager/on_start.go` — read full, find all 3 hardcoded URLs
  - `apps/eventsub/internal/manager/subscribe.go` — read full, find hardcoded URL
  - `apps/eventsub/internal/manager/unsubscribe.go` — read to confirm no DefaultClient usage
  - `libs/twitch/mock_roundtripper.go` — `NewMockRoundTripper(cfg)` function
  - `github.com/kvizyx/twitchy/eventsub` — `WebsocketWithServerURL` option (exists in `websocket_option.go:25`)

  **Acceptance Criteria**:

  ```
  Scenario: Build succeeds
    Tool: Bash
    Steps:
      1. go build ./apps/eventsub/...
    Expected Result: exit code 0
    Evidence: .sisyphus/evidence/task-5-build.txt

  Scenario: Eventsub conduit ensured via mock
    Tool: Bash
    Steps:
      1. docker compose -f docker-compose.dev.yml up twitch-mock -d && sleep 5
      2. TWITCH_MOCK_ENABLED=true go run ./apps/eventsub/cmd/main.go 2>&1 | tee /tmp/eventsub.log &
      3. sleep 10
      4. grep "conduit" /tmp/eventsub.log
    Expected Result: log contains conduit creation/ensured message (no error)
    Evidence: .sisyphus/evidence/task-5-conduit.txt

  Scenario: Eventsub WS welcome received via mock
    Tool: Bash
    Steps:
      1. grep -i "welcome\|session" /tmp/eventsub.log
    Expected Result: log shows websocket session established (welcome received)
    Evidence: .sisyphus/evidence/task-5-ws-welcome.txt
  ```

  **Commit**: YES
  - Message: `feat(eventsub): inject httpClient and configurable API/WS URLs for mock mode`
  - Files: `apps/eventsub/internal/manager/manager.go`, `apps/eventsub/internal/manager/on_start.go`, `apps/eventsub/internal/manager/subscribe.go`

---

- [x] 6. Make `defaultBot.go` seed configurable for mock

  **What to do**:
  In `libs/migrations/seeds/defaultBot.go`:
  - The file contains `http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)` (approximately line 51)
  - Replace hardcoded URL with: `validateURL := "https://id.twitch.tv/oauth2/validate"` at top, then when `config.TwitchMockEnabled`, use `config.TwitchMockAuthUrl + "/oauth2/validate"` instead
  - Config is already available via dependency injection in the seed (check how it's passed)

  **Note**: This file may need to accept `config` as a parameter if it doesn't already. Follow the existing pattern in other seed files.

  **Must NOT do**:
  - Do NOT change any other logic in the seed
  - Do NOT change behavior when `TwitchMockEnabled == false`

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Single URL replacement in one file
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 3, 4, 5, 7, 8)
  - **Blocks**: Task 10
  - **Blocked By**: Tasks 1, 2

  **References**:
  - `libs/migrations/seeds/defaultBot.go` — read full file
  - `libs/config/config.go` — `TwitchMockEnabled`, `TwitchMockAuthUrl`
  - Other seed files in `libs/migrations/seeds/` — understand how config is accessed

  **Acceptance Criteria**:

  ```
  Scenario: Build succeeds
    Tool: Bash
    Steps:
      1. go build ./libs/migrations/...
    Expected Result: exit code 0
    Evidence: .sisyphus/evidence/task-6-build.txt

  Scenario: Seed uses mock URL when enabled
    Tool: Bash
    Steps:
      1. grep "TwitchMock\|oauth2/validate" libs/migrations/seeds/defaultBot.go
    Expected Result: file contains reference to TwitchMockEnabled and configurable URL (not only hardcoded id.twitch.tv)
    Evidence: .sisyphus/evidence/task-6-grep.txt
  ```

  **Commit**: YES
  - Message: `feat(seeds): make bot validate URL configurable for mock mode`
  - Files: `libs/migrations/seeds/defaultBot.go`

---

- [x] 7. Add `twitch-mock` service to `docker-compose.dev.yml`

  **What to do**:
  In `docker-compose.dev.yml`:
  - Add new service `twitch-mock` with:
    ```yaml
    twitch-mock:
      # DEV ONLY — NEVER IN PRODUCTION
      build:
        context: .
        dockerfile: apps/twitch-mock/Dockerfile
      ports:
        - "7777:7777" # HTTP (OAuth + Helix)
        - "8081:8081" # WebSocket EventSub
        - "3333:3333" # Admin UI
      environment:
        - SITE_BASE_URL=${SITE_BASE_URL:-http://localhost:3005}
      healthcheck:
        test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:7777/health"]
        interval: 5s
        timeout: 3s
        retries: 10
        start_period: 10s
      networks:
        - <existing_network_name> # use the same network as other services
    ```
  - Check what network name is used by other services in the file and use the same
  - Check if there's a `wget` or `curl` available in the Dockerfile (use whichever is available)

  **Must NOT do**:
  - Do NOT add `twitch-mock` to `docker-compose.stack.yml`
  - Do NOT add `depends_on: twitch-mock` to other services (this would break non-mock setups) — mock is opt-in

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: YAML file addition
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 3, 4, 5, 6, 8)
  - **Blocks**: Task 10
  - **Blocked By**: Task 2

  **References**:
  - `docker-compose.dev.yml` — read full file, match existing service format and network name
  - `apps/twitch-mock/Dockerfile` — will exist after Task 2

  **Acceptance Criteria**:

  ```
  Scenario: twitch-mock service starts and health check passes
    Tool: Bash
    Steps:
      1. docker compose -f docker-compose.dev.yml up twitch-mock -d
      2. sleep 15
      3. docker compose -f docker-compose.dev.yml ps twitch-mock
    Expected Result: status shows "healthy" (not "starting" or "unhealthy")
    Evidence: .sisyphus/evidence/task-7-healthcheck.txt

  Scenario: mock endpoints accessible from host
    Tool: Bash
    Steps:
      1. curl -f http://localhost:7777/health
    Expected Result: 200 OK
    Evidence: .sisyphus/evidence/task-7-health-curl.txt
  ```

  **Commit**: YES
  - Message: `feat(docker): add twitch-mock service to docker-compose.dev.yml`
  - Files: `docker-compose.dev.yml`

---

- [x] 8. Create `.env.mock.example` and `apps/twitch-mock/README.md`

  **What to do**:

  **`.env.mock.example`** — create at repo root:

  ```env
  # Copy to .env to use Twitch mock server for development
  # Enables local development without a real Twitch account

  # Mock server
  TWITCH_MOCK_ENABLED=true
  TWITCH_MOCK_API_URL=http://localhost:7777
  TWITCH_MOCK_AUTH_URL=http://localhost:7777
  TWITCH_MOCK_WS_URL=ws://localhost:8081/ws

  # Fake Twitch app credentials (any non-empty values)
  TWITCH_CLIENTID=mock-client-id
  TWITCH_CLIENTSECRET=mock-client-secret

  # Fake bot tokens (used by defaultBot seed)
  BOT_ACCESS_TOKEN=mock-bot-token
  BOT_REFRESH_TOKEN=mock-bot-refresh-token

  # Mock user IDs (fixed — do not change)
  # Broadcaster: ID=12345, login=mockstreamer
  # Bot: ID=67890, login=mockbot
  ```

  **`apps/twitch-mock/README.md`** — create with:
  - Setup instructions (copy .env.mock.example, docker compose up, bun dev)
  - Fake user IDs table (broadcaster, bot)
  - List of mocked endpoints (Phase 1 list)
  - How to trigger EventSub events (Admin UI at :3333 + REST API examples)
  - What is NOT mocked (Twitch player iframe, IRC)
  - Troubleshooting section

  **Must NOT do**:
  - Do NOT copy actual secrets into the file
  - Do NOT add more than the fields listed above

  **Recommended Agent Profile**:
  - **Category**: `writing`
    - Reason: Documentation + example config file
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 3, 4, 5, 6, 7)
  - **Blocks**: nothing
  - **Blocked By**: Task 1

  **References**:
  - `libs/config/config.go` — read to get all TwitchMock field names and env tags
  - `.env.example` at repo root — follow same format/style
  - `apps/twitch-mock/README.md` — will be a new file

  **Acceptance Criteria**:

  ```
  Scenario: env.mock.example exists with required keys
    Tool: Bash
    Steps:
      1. grep -c "TWITCH_MOCK_ENABLED\|TWITCH_MOCK_API_URL\|TWITCH_MOCK_AUTH_URL\|TWITCH_MOCK_WS_URL\|BOT_ACCESS_TOKEN\|BOT_REFRESH_TOKEN" .env.mock.example
    Expected Result: output is "6"
    Evidence: .sisyphus/evidence/task-8-env-keys.txt

  Scenario: README exists and has content
    Tool: Bash
    Steps:
      1. wc -l apps/twitch-mock/README.md
    Expected Result: line count > 50
    Evidence: .sisyphus/evidence/task-8-readme.txt
  ```

  **Commit**: YES
  - Message: `docs: add .env.mock.example and apps/twitch-mock/README.md`
  - Files: `.env.mock.example`, `apps/twitch-mock/README.md`

---

- [x] 9. Unit tests for `apps/twitch-mock/`

  **What to do**:
  Write Go unit tests in `apps/twitch-mock/` covering:
  1. `GET /health` → 200 OK
  2. `POST /oauth2/token` with `grant_type=client_credentials` → returns mock app token with `expires_in=99999999`
  3. `POST /oauth2/token` with `grant_type=authorization_code&code=mock_code` → returns mock user token
  4. `POST /oauth2/token` with `grant_type=refresh_token` → returns refreshed token
  5. `GET /oauth2/validate` with `Authorization: OAuth mock-bot-token` → returns mockbot user
  6. `GET /helix/users` → returns MockBroadcaster in data array
  7. `POST /helix/eventsub/conduits` → creates conduit, GET returns it
  8. `PATCH /helix/eventsub/conduits/shards` → updates shard session_id
  9. `POST /helix/eventsub/subscriptions` → returns 202, subscription stored
  10. `DELETE /helix/eventsub/subscriptions/{id}` → returns 204
  11. `GET /helix/chat/messages` (unknown endpoint) → returns 200 with `{"data":[]}`
  12. WebSocket: connect → receive `session_welcome` within 1s
  13. WebSocket: wait 26s → receive `session_keepalive`
  14. Admin trigger: `POST /admin/trigger/channel.follow` → WS client receives notification

  Use `net/http/httptest` for HTTP tests, `github.com/coder/websocket` or `gorilla/websocket` for WS tests.

  **Must NOT do**:
  - Do NOT write integration tests here (no real services needed)
  - Do NOT use any external test framework (only stdlib `testing`)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: WebSocket test harness + HTTP test coverage requires careful setup
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (alone)
  - **Blocks**: nothing
  - **Blocked By**: Task 2

  **References**:
  - `apps/twitch-mock/internal/handlers/` — test targets (from Task 2)
  - `apps/twitch-mock/internal/websocket/eventsub.go` — WS server to test
  - `apps/twitch-mock/internal/admin/admin.go` — admin trigger to test
  - Existing test files in project (e.g., `apps/api-gql/`) — follow test patterns

  **Acceptance Criteria**:

  ```
  Scenario: All unit tests pass
    Tool: Bash
    Steps:
      1. cd apps/twitch-mock && go test ./... -v -timeout 60s
    Expected Result: all tests PASS, 0 failures, 0 errors
    Evidence: .sisyphus/evidence/task-9-tests.txt

  Scenario: Test coverage meets minimum
    Tool: Bash
    Steps:
      1. cd apps/twitch-mock && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out | tail -1
    Expected Result: total coverage ≥ 60%
    Evidence: .sisyphus/evidence/task-9-coverage.txt
  ```

  **Commit**: YES
  - Message: `test(twitch-mock): add unit tests for mock server endpoints and WS protocol`
  - Files: `apps/twitch-mock/**/*_test.go`

---

- [ ] 10. E2E Verification — Full Login Flow with Mock (MCP Playwright)

  **What to do**:
  This task is executed by the agent using the **MCP Playwright browser tool** (not a repo-hosted test file). No test file is written to disk. The agent must directly drive the browser through the full login and event-trigger flow.

  **Why MCP Playwright (not a repo test file)**: The project has no `@playwright/test` dependency, no `playwright.config.*`, and no `test:e2e` script in `package.json`. Rather than add Playwright as a dev dependency (scope creep), use the MCP Playwright tool that the executing agent already has access to.

  **Preconditions** (executor must verify before starting):
  1. `docker compose -f docker-compose.dev.yml ps twitch-mock` → shows status `healthy`
  2. All Go services running with `TWITCH_MOCK_ENABLED=true`
  3. Frontend accessible at `http://localhost:3005`

  **Steps to execute**:
  1. Open browser at `http://localhost:3005`
  2. Find and click "Login with Twitch" button (selector: `button` with text matching `Login` or `Twitch`, or `a[href*="authLink"]`)
  3. Mock OAuth auto-redirects — wait for URL to match `**/dashboard/**` (timeout 15s)
  4. Assert: page URL contains `/dashboard`
  5. Assert: text "MockStreamer" is visible in the page (header or user menu)
  6. Take screenshot → save as `.sisyphus/evidence/task-10-e2e-screenshot.png`
  7. In browser console, assert: zero `console.error()` calls during the login flow
  8. Trigger follow event via curl:
     ```bash
     curl -s -X POST http://localhost:3333/admin/trigger/channel.follow \
       -H "Content-Type: application/json" \
       -d '{"from_user_id":"99999","from_user_login":"testfollower","followed_at":"2026-04-12T00:00:00Z"}'
     ```
  9. Assert: curl returns HTTP 200

  **Must NOT do**:
  - Do NOT create any `.spec.ts` or `e2e/` test files in the repository
  - Do NOT install `@playwright/test` or add any npm/bun packages
  - Do NOT require manual browser interaction — fully automated via MCP tool

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: E2E verification requiring full-stack orchestration and browser automation via MCP
  - **Skills**: [`playwright`]
    - `playwright`: MCP browser tool for navigation, click, assertion, screenshot

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (alone — needs all services wired including Task 11)
  - **Blocks**: Final Verification Wave
  - **Blocked By**: Tasks 3, 4, 5, 6, 7, 11

  **References**:
  - `frontend/dashboard/src/plugins/router.ts` — auth guard logic (what triggers redirect to `/`)
  - `web/layers/landing/pages/login.client.vue` — login button location, OAuth flow
  - `apps/twitch-mock/internal/handlers/auth.go` — mock OAuth auto-redirect behaviour

  **Acceptance Criteria**:

  ```
  Scenario: Full login flow succeeds with mock
    Tool: MCP Playwright (browser_navigate, browser_click, browser_snapshot, browser_screenshot)
    Preconditions:
      - docker compose -f docker-compose.dev.yml up -d (twitch-mock healthy)
      - TWITCH_MOCK_ENABLED=true bun dev (all services running, port 3005 reachable)
    Steps:
      1. browser_navigate to http://localhost:3005
      2. browser_snapshot — locate element with text "Login" or "Login with Twitch"
      3. browser_click that element
      4. browser_wait_for URL to contain "/dashboard" (timeout: 15000ms)
      5. browser_snapshot — assert text "MockStreamer" visible in DOM
      6. browser_screenshot → .sisyphus/evidence/task-10-e2e-screenshot.png
    Expected Result: dashboard page visible, "MockStreamer" text present in DOM
    Failure Indicators: URL still at localhost:3005/, "MockStreamer" not found, 500 error page
    Evidence: .sisyphus/evidence/task-10-e2e-screenshot.png

  Scenario: EventSub event trigger returns 200
    Tool: Bash (curl)
    Steps:
      1. curl -s -o /tmp/trigger.json -w "%{http_code}" -X POST http://localhost:3333/admin/trigger/channel.follow \
           -H "Content-Type: application/json" \
           -d '{"from_user_id":"99999","from_user_login":"testfollower","followed_at":"2026-04-12T00:00:00Z"}'
    Expected Result: HTTP status code 200
    Failure Indicators: 4xx or 5xx status code
    Evidence: .sisyphus/evidence/task-10-event-trigger.txt

  Scenario: No JavaScript console errors during login
    Tool: MCP Playwright (browser_console_messages)
    Steps:
      1. After completing the login flow, call browser_console_messages(level="error")
      2. Assert: returned list is empty (no error-level messages)
    Expected Result: empty error list
    Evidence: .sisyphus/evidence/task-10-console-errors.txt
  ```

  **Commit**: NO (this task produces only evidence files, not source code changes)

---

- [x] 11. Make `authLink` GraphQL resolver return mock auth URL when mock enabled

  **What to do**:
  In `apps/api-gql/internal/delivery/gql/resolvers/user.resolver.go`, find the `AuthLink` resolver function (line ~279). Currently it calls `twitchClient.GetAuthorizationURL(...)` which always builds a URL pointing to `https://id.twitch.tv/oauth2/authorize`.

  When `cfg.TwitchMockEnabled == true`, instead of building a real Twitch OAuth URL, construct and return the mock authorize URL directly:

  ```go
  if r.deps.Config.TwitchMockEnabled {
      // Build mock authorize URL — mock server will auto-redirect back to dashboard
      mockAuthUrl := fmt.Sprintf(
          "%s/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=&state=%s",
          r.deps.Config.TwitchMockAuthUrl,
          r.deps.Config.TwitchClientId,
          url.QueryEscape(u.JoinPath("login").String()),
          url.QueryEscape(state),
      )
      return mockAuthUrl, nil
  }
  ```

  This ensures that when the browser clicks "Login with Twitch", it goes to `http://localhost:7777/oauth2/authorize` (the mock server) rather than `https://id.twitch.tv`. The mock server's `GET /oauth2/authorize` handler (from Task 2) will redirect back to `{SITE_BASE_URL}/login?code=mock_code&state={state}`, completing the OAuth flow.

  **Must NOT do**:
  - Do NOT change the resolver when `TwitchMockEnabled == false` — keep existing behaviour
  - Do NOT skip the `state` parameter — it must be preserved for redirect handling
  - Do NOT modify the GraphQL schema — only the resolver implementation

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Single targeted change in one resolver function
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (can run alongside Tasks 3–8)
  - **Blocks**: Task 10
  - **Blocked By**: Tasks 1, 2 (needs config fields + mock server's authorize endpoint spec)

  **References**:
  - `apps/api-gql/internal/delivery/gql/resolvers/user.resolver.go:279-340` — `AuthLink` resolver, understand state encoding and redirect URI construction
  - `libs/config/config.go` — `TwitchMockEnabled`, `TwitchMockAuthUrl`
  - `apps/twitch-mock/internal/handlers/auth.go` — mock `GET /oauth2/authorize` handler (from Task 2) — must understand what it expects in query params and what it redirects to

  **Acceptance Criteria**:

  ```
  Scenario: Build succeeds
    Tool: Bash
    Steps:
      1. go build ./apps/api-gql/...
    Expected Result: exit code 0
    Evidence: .sisyphus/evidence/task-11-build.txt

  Scenario: authLink returns mock URL when TWITCH_MOCK_ENABLED
    Tool: Bash
    Steps:
      1. grep "TwitchMockEnabled\|TwitchMockAuthUrl\|oauth2/authorize" apps/api-gql/internal/delivery/gql/resolvers/user.resolver.go
    Expected Result: file contains reference to TwitchMockEnabled guard and mock auth URL construction (not only GetAuthorizationURL)
    Evidence: .sisyphus/evidence/task-11-grep.txt

  Scenario: Mock login redirect goes to mock server (integration)
    Tool: Bash (curl to GraphQL)
    Steps:
      1. Start mock server: docker compose -f docker-compose.dev.yml up twitch-mock -d
      2. TWITCH_MOCK_ENABLED=true go run ./apps/api-gql/cmd/main.go &
      3. sleep 5
      4. curl -s -X POST http://localhost:3009/graphql \
           -H "Content-Type: application/json" \
           -d '{"query":"{ authLink(redirectTo: \"/\") }"}'
    Expected Result: JSON response where `data.authLink` starts with `http://localhost:7777/oauth2/authorize`
    Evidence: .sisyphus/evidence/task-11-authlink.json
  ```

  **Commit**: YES
  - Message: `feat(api-gql): return mock OAuth authorize URL when TWITCH_MOCK_ENABLED`
  - Files: `apps/api-gql/internal/delivery/gql/resolvers/user.resolver.go`

---

## Final Verification Wave (MANDATORY — after ALL implementation tasks)

> 4 review agents run in PARALLEL. ALL must APPROVE. Present consolidated results to user and get explicit "okay" before completing.

- [ ] F1. **Plan Compliance Audit** — `oracle`
      Read the plan end-to-end. For each "Must Have": verify implementation exists (read file, curl endpoint, run command). For each "Must NOT Have": search codebase for forbidden patterns — reject with file:line if found. Check evidence files exist in `.sisyphus/evidence/`. Compare deliverables against plan.
      Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT: APPROVE/REJECT`

- [ ] F2. **Code Quality Review** — `unspecified-high`
      Run `go build ./...` + linter across all changed files. Review for: `as any`, empty catches, commented-out code, unused imports. Check: `if config.TwitchMockEnabled` guard is present in ALL 4 modified services. Check `twitch-mock` is NOT referenced in `docker-compose.stack.yml`.
      Output: `Build [PASS/FAIL] | Lint [PASS/FAIL] | Guard [N/N services] | VERDICT`

- [ ] F3. **Real Manual QA** — `unspecified-high` (+ `playwright` skill)
      Cold-start from scratch:

  ```
  docker compose -f docker-compose.dev.yml down -v
  cp .env.mock.example .env
  docker compose -f docker-compose.dev.yml up -d
  bun dev
  ```

  Execute ALL QA scenarios from ALL tasks. Save evidence. Test:
  - OAuth flow (curl `/oauth2/token`, `/oauth2/validate`)
  - Helix users endpoint
  - Conduit CRUD
  - WS welcome message
  - Admin UI loads and triggers events
  - E2E login flow
    Output: `Scenarios [N/N pass] | Cold start [PASS/FAIL] | Evidence files [N] | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`
      For each task: read "What to do", check actual diff (`git diff`). Verify 1:1 — everything in spec was built, nothing beyond spec. Check "Must NOT do" compliance (no mutations to `helix.AuthBaseURL`, no `twitch-mock` in production stack, no stateful endpoints beyond conduits). Flag unaccounted changes.
      Output: `Tasks [N/N compliant] | Must NOT violations [CLEAN/N] | Unaccounted [CLEAN/N] | VERDICT`

---

## Commit Strategy

```
feat(config): add TWITCH_MOCK_* config fields
feat(twitch-mock): add Twitch API mock server (HTTP+WS+AdminUI)
feat(libs/twitch): add mockRoundTripper, redirect helix to mock when enabled
feat(tokens): redirect helix client to mock when TWITCH_MOCK_ENABLED
feat(eventsub): inject httpClient and configurable API/WS URLs for mock mode
feat(seeds): make bot validate URL configurable for mock mode
feat(api-gql): return mock OAuth authorize URL when TWITCH_MOCK_ENABLED
feat(docker): add twitch-mock service to docker-compose.dev.yml
docs: add .env.mock.example and apps/twitch-mock/README.md
test(twitch-mock): add unit tests for mock server endpoints and WS protocol
```

---

## Success Criteria

### Verification Commands

```bash
# 1. Cold start
docker compose -f docker-compose.dev.yml down -v
cp .env.mock.example .env
docker compose -f docker-compose.dev.yml up -d
bun dev

# 2. Mock health
curl http://localhost:7777/health  # Expected: 200

# 3. OAuth
curl -X POST http://localhost:7777/oauth2/token \
  -d "grant_type=client_credentials&client_id=mock&client_secret=mock" \
  # Expected: {"access_token":"mock-app-token","expires_in":99999999,...}

# 4. Eventsub logs (within 30s of bun dev)
grep -i "conduit\|welcome" <eventsub-service-log>
# Expected: "conduit ensured" + websocket session established

# 5. Dashboard login
# Browser: http://localhost:3005 → Login → Dashboard → "MockStreamer" visible

# 6. All unit tests
go test ./apps/twitch-mock/... -v
# Expected: all PASS

# 7. E2E (run by executing agent via MCP Playwright — no repo test file)
# Steps: browser_navigate http://localhost:3005 → click Login → assert /dashboard → assert "MockStreamer"
# Evidence: .sisyphus/evidence/task-10-e2e-screenshot.png
```

### Final Checklist

- [ ] All "Must Have" present
- [ ] All "Must NOT Have" absent
- [ ] `TWITCH_MOCK_ENABLED=false` → no behavioral change
- [ ] `twitch-mock` NOT in `docker-compose.stack.yml`
- [ ] `helix.AuthBaseURL` package global NOT mutated
- [ ] All unit tests pass
- [ ] E2E login flow passes
- [ ] Cold start completes within 90 seconds
