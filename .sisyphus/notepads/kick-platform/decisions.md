# Decisions — kick-platform

## Architecture Decisions

- **Platform as enum everywhere**: `"twitch"` / `"kick"` are typed enums at all layers:
  - **DB**: `CREATE TYPE platform AS ENUM ('twitch', 'kick')` — all `platform` columns use this Postgres type
  - **Go**: `type Platform string` with consts in `libs/entities/platform/platform.go`; no bare strings
  - **GraphQL**: `enum Platform { TWITCH KICK }` in `apps/api-gql/schema/platform.graphql`
  - **TypeScript**: generated `Platform` enum used throughout frontend; no raw `"twitch"`/`"kick"` literals
- **Auth model**: Universal login screen with Twitch / Kick platform selection. First linked platform = primary; others are "linked accounts"
- **users.id migration**: Big-bang with maintenance window (not zero-downtime). All sessions flushed (Redis FLUSHDB) at deploy time.
- **user_platform_accounts**: New table storing platform identity + tokens per platform. FK → users.internal_id (→ renamed users.id)
- **channels schema**: One row per platform per user (e.g., Twitch channel + Kick channel = 2 rows per user)
- **Platform switchers**: `platforms platform[] DEFAULT '{}'` on commands/timers/keywords. Empty = all platforms.
- Execution-layer filtering should be applied after loading models, not in SQL, so empty platform lists stay backward compatible.
- **Kick EventSub**: Official Kick EventSub Webhooks via HTTP. Broadcaster must authorize Twir; broadcaster's own token used for subscriptions.
- **Kick bot**: Single `kick_bots` table; bot account sends via `POST /public/v1/chat` with bot token.
- **NATS bus**: Parallel queue strategy — keep `twitch.TwitchChatMessage` alongside new generic `ChatMessage`. Do NOT remove old queue until ALL consumers migrated.
- **Token storage**: `user_platform_accounts` for user tokens; `tokens` table repurposed for bot tokens only.
- **7TV Kick**: Only `GetProfileByKickId` REST call (`https://7tv.io/v3/users/kick/{channelID}`) in scope. No SSE.
- **UUID generation**: Always use `uuidv7()` (NOT `gen_random_uuid()`) for new UUID PKs in migrations. Project runs on Postgres 18+ which ships `uuidv7()` natively. UUIDv7 is time-sortable (better index locality). This applies to all new tables and columns added in this branch and any future work.
- Implemented Kick OAuth code exchange as a dedicated /auth/kick/code route while keeping legacy Twitch /auth handler untouched for backward compatibility.
- Kick bus topics should mirror the existing Twitch package pattern: one subject constant plus one exported message struct per file, under `libs/bus-core/kick/`.
- For 7TV Kick profile lookups, use the existing GraphQL `userByConnection(platform: KICK, id: $id)` pattern in both v2 and v3 clients instead of the REST URL mentioned in the plan.
