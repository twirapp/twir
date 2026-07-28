# Wave 1: Bots Outbound

## Findings
- `MessageHandler.Handle` fans out in parallel and logs handler failures without returning them.
- Greetings discards send/shoutout errors; moderation has Redis/error swallowing; keywords nests goroutine fanout.
- Explicit dispatch may break after the first same-platform binding even if disabled; implicit dispatch silently ignores unsupported capability.
- Twitch rate-limit behavior can no-op, `SkipRateLimits` is unused, toxicity is disabled; Kick drops 429; VK no-ops without bot user.
- `/timeout` parsing can index beyond bounds.
- ClickHouse batch failures and VK provider 404s are log-only at upper layers.

## EXPAND
- Verify disabled-first-binding dispatch edge.
- Verify malformed timeout panic and nested fanout race/error behavior.
- Audit provider send contracts and typed errors.
