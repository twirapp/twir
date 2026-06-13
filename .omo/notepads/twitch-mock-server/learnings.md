## 2026-04-12

- `TWITCH_MOCK_API_URL` must point to the mock host root (`http://localhost:7777`), not `/helix`, because the helix client adds that prefix itself.
- EventSub websocket debug output should use structured `slog`/logger helpers instead of `fmt.Println` or `pretty.Println`.
- Twitch mock admin trigger docs now use `/admin/trigger/channel.*` routes and the Nuxt dev server port is `3010`.
- Mock EventSub broadcasts should target only websocket session IDs currently registered on conduit shards; otherwise stale or extra clients receive duplicate notifications.
- Local mock mode should create EventSub conduits with a single shard when `TwitchMockEnabled` is true, while production keeps the existing 3-shard default.
