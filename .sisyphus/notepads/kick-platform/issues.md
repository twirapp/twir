# Issues — kick-platform

## Known Gotchas

- `users.id` is currently a TEXT column holding the Twitch ID (e.g., "12345678"). After migration it becomes a UUID.
- `channels.id` also equals the Twitch ID. After migration it stays a TEXT FK referencing the new users.id UUID.
- The auth handler at `apps/api-gql/internal/delivery/http/routes/auth/post-code.go` uses GORM + `libs/gomodels` — leave Twitch handler untouched; new Kick handler uses pgx only.
- Session stores `model.Users{ID: "twitch_id"}` + `helix.User` in Redis via SCS — ALL sessions must be invalidated at migration time via `FLUSHDB`.
- `tokens` repo `GetByUserID(userID string)` currently queries by Twitch ID — must be updated to use internal UUID post-migration.
- EventSub app currently uses Twitch WebSocket conduits — does NOT have an HTTP endpoint for webhooks yet.
- FK ordering in migration is critical: cannot add FKs to `users.internal_id` until UNIQUE constraint added.
- go vet on apps/api-gql failed on pre-existing slog.Error argument patterns and one invalid json tag on an unexported field; fixed them alongside auth work to satisfy required verification.
- `go build ./libs/bus-core/...` succeeded after adding the new Kick bus package; no extra bus-core wiring was needed for compile coverage.
- Keyword execution lives in `apps/bots/internal/messagehandler`, not parser, so platform filtering for keywords had to be applied there too.
