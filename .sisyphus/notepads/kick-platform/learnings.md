- Goose SQL migrations in this repo use `-- +goose Up/Down` with `-- +goose StatementBegin/End` wrappers.
- GraphQL schema files for api-gql are generated from `apps/api-gql/internal/delivery/gql/schema/`.
- Platform enum foundation is now shared across Go, Postgres, and GraphQL with matching `twitch`/`kick` values.
- `users.id` can be swapped to internal UUID in-place by staging a temp `internal_id`, converting all direct user FKs plus `channels.id`, then renaming old `users.id` to `twitch_id` for backward compatibility.
- `user_platform_accounts.platform` should use the existing Postgres `platform` enum and the Go entity should use `libs/entities/platform.Platform`.
- Deploying the user UUID migration requires a Redis `FLUSHDB` maintenance step because stored sessions still contain legacy Twitch IDs.
- `kick_bots` is a dedicated bot-token table with no platform column; migration uses Goose SQL wrappers and the entity follows the repo's nil-pattern with `isNil` + `IsNil()` + `Nil` singleton.

- Added platforms[] columns for commands/timers/keywords using the shared platform enum; migration runner completed successfully.
- pgx scanning/writing worked with the new platforms field in repository code, and the repo-level build checks passed in the affected Go modules.
- Channels multi-platform migration can safely reuse `channels.id` as the legacy-to-new UUID mapping source after renaming it to `user_id`, then re-point all child FKs by updating UUID values before re-adding the original FK definitions.
- QA for the new `(user_id, platform)` uniqueness must copy required non-null channel fields like `botId` from an existing Twitch row; bare `INSERT (user_id, platform)` will fail on existing table constraints unrelated to the new schema.
- Tokens repository callers now need to parse legacy string user IDs into `uuid.UUID` at the boundary before hitting the pgx repository API.
- Bus-core generic chat messages now live in `libs/bus-core/generic/chat-message.go`; `ChannelID` is the internal surrogate `channels.id` and `UserID` is the internal `users.id`.
- `libs/bus-core/bus.go` wires `ChatMessagesGeneric` to `chat.messages.generic` and `Parser.ProcessGenericMessage` to `parser.process_generic_message` without touching the existing Twitch queues.
- `apps/parser/internal/types/ParseContext` now carries a `Platform` string, and parser constructors should populate it from the originating message platform.
- Shared platform filtering is easiest to keep consistent with a small helper on `libs/entities/platform`, then reuse it in parser/timer/keyword execution paths.

## [T5] PlatformProvider interface + Twitch implementation
- `platform.PlatformProvider` interface lives in `apps/api-gql/internal/platform/provider.go`; Twitch impl in `internal/platform/twitch/provider.go`.
- Twitch provider creates a plain `helix.Client` (no bus/app-token) using only ClientID, ClientSecret, RedirectURI — sufficient for OAuth exchange flows.
- `GetAuthURL` ignores `codeChallenge` (Twitch doesn't use PKCE); scopes are embedded in the provider.
- `ExchangeCode` ignores `codeVerifier` param — pass `_` to satisfy interface.
- Registered in `cmd/main.go` via `fx.Annotate(twitchplatform.New, fx.As(new(platform.PlatformProvider)))`.
- `TwitchMockEnabled` → `apiBaseURL` override pattern mirrors libs/twitch/twitch.go.

- Kick provider uses exported PKCE helpers so the authorize route can store a session-only verifier and return a Kick OAuth URL with S256 challenge params.
- Kick config now supports KICK_CLIENT_ID, KICK_CLIENT_SECRET, and optional KICK_REDIRECT_URL; default callback falls back to {SITE_BASE_URL}/login/kick.
- Kick /auth/kick/code needs session code verifier retrieval via Auth.Get plus explicit Auth.Commit after multiple Put calls to avoid per-key helper proliferation.
- Kick bus-core topics live in a dedicated `libs/bus-core/kick` package with raw event structs and subject constants; `go build ./libs/bus-core/...` passes after adding them.
- `apps/eventsub/internal/bus-listener/bus-listener.go` can branch on `channel.Platform` after `channelsRepo.GetByID`: Twitch keeps legacy topic-based EventSub subscription flow, while Kick uses `user_platform_accounts.GetByUserIDAndPlatform` plus `kick.SubscriptionManager.SubscribeAll`.
- Twitch EventSub chat handling can dual-publish into generic queues by mapping `eventsub.ChannelChatMessageEvent` into `generic.ChatMessage`; guard `data.Message` before reading `.Text` and keep legacy Twitch bus publishes intact.

## T20: Kick Resubscribe Job

- `GetAllByPlatform` added to `user_platform_accounts.Repository` interface + pgx impl using `WHERE platform = $1::platform` pattern
- `SubscriptionLister` interface defined in `kick` package so `*SubscriptionManager` satisfies it for tests (Option A from plan)
- fx wiring: `ResubscribeJobOpts.SubManager` uses concrete `*SubscriptionManager` so fx can resolve it; the `ResubscribeJob.subManager` field holds the interface
- `ResubscribeJob` lifecycle: OnStart fires `go j.Start(ctx)` — goroutine exits on ctx cancellation (ticker.Stop + select on ctx.Done)
- Existing `handlers_test.go` mock needed `GetAllByPlatform` added when the interface was extended
- `allEventTypesPresent` uses a map for O(n) lookup over all 4 EventTypes

## T22: Dual-subscribe ProcessMessageAsCommand + ProcessGenericMessage with Redis dedup

### Redis dedup pattern
- Key: `parser:dedup:{messageID}`
- `redis.SetNX(ctx, key, "1", 60*time.Second)` — returns `true` if key was SET (new msg), `false` if existed (dup)
- So: `isDuplicate = !set` 
- On Redis error: log and proceed (don't drop message)
- Empty messageID: skip dedup (return false, nil)

### CommandsBus.redis field
- Added `redis *redis.Client` field to `CommandsBus` struct
- Populated from `s.Redis` in `New()` — no signature change needed

### generic.ChatMessage field mapping to TwitchChatMessage
- `msg.ChannelID` → `BroadcasterUserId`
- `msg.PlatformChannelID` → `BroadcasterUserLogin` / `BroadcasterUserName`
- `msg.UserID` → `ChatterUserId`
- `msg.SenderLogin` → `ChatterUserLogin`
- `msg.SenderDisplayName` → `ChatterUserName`
- `msg.MessageID` → `MessageId`
- `msg.Badges[].SetID` → `twitch.ChatMessageBadge.SetId`
- `msg.Badges[].Text` → `twitch.ChatMessageBadge.Info`
- `msg.Text` → `Message.Text`

### Important caveat
- `ProcessChatMessage` checks `data.EnrichedData.DbUser == nil` and returns error if nil
- For Kick messages (no enriched data), `ProcessChatMessage` will return an error — full Kick command execution requires a separate enrichment step (future task)
- The subscription + dedup infrastructure is in place; Kick command execution needs enrichment pipeline

### Unsubscribe
- Must add `c.bus.Parser.ProcessGenericMessage.Unsubscribe()` to `Unsubscribe()` method

### TODO comment required by spec
```go
// TODO(Phase-2): remove ProcessMessageAsCommand subscription once all consumers migrated off it
```

## T22-fix: Propagate platform from generic message into ParseContext

### Problem
`ParseCommandResponses` hardcoded `Platform: "twitch"` in the `ParseContext` it builds at line 339-360 of `commands.go`. All calls via `ProcessGenericMessage` for Kick messages would incorrectly get `Platform: "twitch"`.

### Fix (Option A — explicit parameter)
- Added `platform string` parameter to both `ParseCommandResponses` and `ProcessChatMessage` signatures
- `ParseCommandResponses` uses `platform` instead of hardcoded `"twitch"` 
- `ProcessChatMessage` receives `platform` and forwards it to `ParseCommandResponses`
- Call sites in `commands-bus.go`:
  - `GetCommandResponse` handler: passes `"twitch"`
  - `ProcessMessageAsCommand` handler: passes `"twitch"`
  - `ProcessGenericMessage` handler: passes `msg.Platform` (the correct value from the generic message)

### Note on `GetCommandResponse`
This handler (used by variable parsing path) also calls `ProcessChatMessage` — updated to pass `"twitch"` since it only receives `TwitchChatMessage`.
