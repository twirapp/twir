# Kick Token Centralization Design

**Goal:** Move Kick bot token refresh ownership into `apps/tokens` and make Kick chat message splitting match Kick API byte limit behavior.

## Current State

- `apps/bots/internal/kick/chat_client.go` reads default Kick bot from `kick_bots` repository.
- Bots service decrypts Kick access and refresh tokens locally.
- Bots service refreshes Kick token locally on `401` and persists updated token.
- Kick message splitting uses rune count, not UTF-8 byte length.
- Local `gokick.ValidateChatMessageContent` also validates by rune count.

## Problems

- Token lifecycle logic lives in more than one service.
- Bots service owns credential refresh and encryption logic it should not own.
- Kick API effectively enforces message content limit by bytes, so Cyrillic and other multibyte text fail earlier than current splitter expects.

## Decision

### Token Ownership

- `apps/tokens` becomes single authority for Kick bot token decrypt, refresh, and persistence.
- `apps/bots/internal/kick/chat_client.go` no longer decrypts or refreshes Kick bot tokens.
- Bots service asks tokens service for bot token through existing `RequestBotToken` bus queue.

### Bot Token Request Shape

- Extend `libs/bus-core/tokens.GetBotTokenRequest` with `Platform platformentity.Platform`.
- Empty platform defaults to Twitch inside tokens service, matching existing app token behavior.
- Kick callers pass `platformentity.PlatformKick` explicitly.

### Kick Bot Resolution

- For Kick bot token requests, tokens service loads default Kick bot from `kick_bots.Repository`.
- If stored token is expired, tokens service refreshes it with Kick OAuth, persists encrypted values, and returns decrypted access token.
- If refresh token response omits new refresh token, keep existing refresh token.

### Kick Chat Sending

- Bots service still sends messages as current default Kick bot identity.
- Chat client requests token once per message part send attempt.
- Unauthorized retry logic is removed from bots service. Refresh responsibility stays in tokens service.
- Rate-limit and forbidden handling stay in bots service because they are send-time concerns, not token lifecycle concerns.

### Byte-Based Splitting

- Replace rune-based Kick splitting with UTF-8 byte-based splitting.
- Preserve valid UTF-8 boundaries by splitting only on rune boundaries.
- Normalize `\n` to spaces before validation and split, preserving current behavior.
- Prefer largest prefix with `len([]byte(part)) <= 500`.

### Shared Kick Validation

- Update local `gokick.ValidateChatMessageContent` to enforce byte length instead of rune length so library behavior matches real API behavior.
- Keep public constant name for now to minimize churn, but use byte-based semantics in implementation and comments.

## Files To Change

- `libs/bus-core/tokens/tokens.go`
- `apps/tokens/internal/bus_listener/bus_listener.go`
- `apps/bots/internal/kick/chat_client.go`
- `apps/bots/app/app.go`
- `gokick/chat.go`
- New tests near changed packages

## Error Handling

- Tokens service wraps Kick bot lookup, decrypt, refresh, encrypt, and persist failures with source context.
- Bots service returns token request failures without retrying locally.
- Bots service still drops `429` send failures and logs `403` failures with bot and broadcaster context.

## Testing

- Add failing tests first for `gokick.ValidateChatMessageContent` byte-length enforcement.
- Add failing tests first for Kick splitter with ASCII, Cyrillic, exact fit, and UTF-8 boundary cases.
- Add failing tests first for tokens service Kick bot request path and Twitch default behavior when platform omitted.
- Add failing tests first proving bots Kick client requests token via bus instead of decrypting locally.

## Non-Goals

- No change to Twitch bot token flow beyond preserving default platform behavior.
- No new token transport protocol beyond existing bus queue.
- No refactor of unrelated Kick chat features.
