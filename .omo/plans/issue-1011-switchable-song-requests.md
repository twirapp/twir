# issue-1011-switchable-song-requests - Work Plan

## TL;DR (For humans)

Реализуем переключаемый режим song requests: сейчас только YouTube, добавляем Spotify. Spotify режим работает по схеме Songify: Twir ведёт собственные записи заявок, а песни добавляются в нативную очередь Spotify через Web API. Удаление конкретной позиции из Spotify невозможно, поэтому `sr wrong` делает отложенный skip (deferred skip), как Songify. Браузер/overlay не требуется — управление идёт с сервера через Spotify Connect device стримера. YouTube режим остаётся без изменений.

## Scope

### Внутри scope
- Поддержка режима `YOUTUBE` и `SPOTIFY` в настройках канала.
- Сохранение режима в `channels_song_requests_settings` с backfill `YOUTUBE` для всех существующих каналов.
- YouTube режим: существующий `!sr`/`sr wrong`/queue/overlay/widget без изменений.
- Spotify режим: поиск трека, валидация, добавление в `POST /me/player/queue`, ведение Twir-записи заявки с requester/status.
- Spotify device discovery и кеширование активного устройства.
- Worker-реконсилятор, который опрашивает current playback и Spotify queue.
- `sr wrong` для Spotify: удаление логической записи + отложенный skip, когда этот трек дойдёт до current playback.
- Счётчики и лимиты (`MaxRequests`, `UserMaxRequests`, per-user cooldown) работают в обоих режимах.
- Donation-based song requests маршрутизируются по выбранному режиму.
- OAuth Spotify: добавить scope `user-modify-playback-state`, UI reconnect для старых подключений.
- Dashboard: переключатель режима, индикатор Spotify connection/device, Spotify search, список заявок со статусами.
- GraphQL: mode в settings, Spotify track search, request lifecycle/status query.
- Миграции и rollback-safe additive schema changes.
- Fixture-based тесты для Spotify HTTP клиента, parser, events и Vue settings.

### За пределами scope (Must NOT have)
- Web Playback SDK в overlay или backend.
- Произвольное удаление/перестановка элементов внутри Spotify native queue.
- Реальный browser playback device внутри Twir.
- Кластерная синхронизация playback state между несколькими dashboard tabs ( optimism обычных мутаций достаточно).
- Автоматическая передача playback на Twir device — только выбор существующего Spotify Connect device стримера.
- Поддержка Spotify podcast episodes, albums, artists, playlists как song request.
- Режим simultaneous YouTube + Spotify queues.

## Verification strategy

- `go test` для новых repository/service packages.
- `bun run test` в `web/layers/dashboard/features/songRequests` для Vue UI и composables.
- `bun cli build gql` + `go run ./apps/api-gql/cmd` smoke test.
- `bun dev` ручной end-to-end: Spotify OAuth reconnect, `!sr`, `sr wrong`, donation, переключение режима, device unavailable.
- Fixtures: Spotify token refresh, `AddToQueue` 204, 401/403/429, no active device, restricted device, duplicate request, queue polling snapshots.
- Regression test на `sr wrong 6` panic при 5 заявках (найденный баг).

## Execution strategy

### Волна 1 — Domain + persistence
1. Добавить `SongRequestMode` enum (`YOUTUBE`, `SPOTIFY`) в `libs/entities/song_requests_settings/entity.go` и `libs/entities/song_request_mode`.
2. Миграция: `channels_song_requests_settings.mode` с default `YOUTUBE`, backfill существующих строк.
3. Новая таблица `spotify_song_requests` с полями: id, channel_id, track_id, track_uri, title, artist, album, duration_ms, requester_user_id, requester_name, requester_display_name, source, queue_position, status, queued_at, playing_observed_at, played_observed_at, cancelled_pending_skip_at, skipped_by_twir_at, removed_or_reconciled_at, unknown_at, created_at, updated_at. Soft-delete нет — статусная машина.
4. Entity + repository + pgx для `spotify_song_requests` с методами Create, GetByID, GetActiveByChannel, GetActiveByRequester, MarkStatus, CancelPendingSkip, CountActiveByChannel, CountActiveByRequester, ListByChannel.
5. Обновить `song_requests_settings` repository/entity upsert и mapper для mode.
6. Тесты: entity validation, migration default, upsert round-trip, active counts, status transitions.

### Волна 2 — Spotify HTTP client + OAuth scopes
7. В `libs/integrations/spotify` добавить:
   - `SearchTracks(ctx, query string, limit int) ([]SpotifyTrack, error)`;
   - `AddToQueue(ctx, trackURI string, deviceID string) error`;
   - `GetCurrentlyPlaying(ctx) (*CurrentlyPlaying, error)`;
   - `GetQueue(ctx) ([]SpotifyTrack, error)`;
   - `GetDevices(ctx) ([]Device, error)`;
   - `SkipNext(ctx, deviceID string) error`.
8. Типизированные ошибки: `ErrNotPremium`, `ErrNoActiveDevice`, `ErrRestrictedDevice`, `ErrRateLimited`, `ErrNotConnected`, `ErrTrackNotFound`, `ErrInsufficientScope`, `ErrQueueTimeout`.
9. Обновить `apps/api-gql/internal/services/spotify_integration/spotify_integration.go`: добавить `user-modify-playback-state` в scopes. Persisted scopes уже хранятся, UI проверит наличие scope и потребует reconnect.
10. Тесты: HTTP fixtures для всех методов и ошибок, token refresh flow, `Retry-After` backoff.

### Волна 3 — Dispatcher + reconciler
11. В `apps/api-gql/internal/services/song_requests` (или новый `spotify_song_requests`) добавить `SpotifySongRequestsService`:
    - `CreateRequest(ctx, channelID, requester, query)` — поиск, валидация, AddToQueue, запись в таблицу;
    - `CancelRequest(ctx, channelID, requesterName, optionalPosition)` — логика `sr wrong`: newest-first active request → `cancelled_pending_skip`, вернуть название;
    - `GetActiveQueue(ctx, channelID)` — логический список;
    - `SelectAndCacheDevice(ctx, channelID)` — device discovery/pick first active non-restricted.
12. Reconciler lifecycle hook в api-gql: фоновый worker с adaptive backoff (idle 5s → 60s cap), который:
    - для каждого активного Spotify-канала опрашивает `GetCurrentlyPlaying` + `GetQueue`;
    - сопоставляет current track URI с `spotify_song_requests` по occurrence/FIFO;
    - обновляет status → `playing_observed` / `played_observed`;
    - если current track — `cancelled_pending_skip`, вызывает `SkipNext` и переводит в `skipped_by_twir`;
    - если active request исчез из queue без playback → `removed_or_reconciled`;
    - обрабатывает rate-limit, device unavailable, token refresh.
13. Тесты: state machine transitions, deferred skip, duplicate URI matching, device disappearance, polling backoff.

### Волна 4 — Parser + donation mode dispatch
14. Переименовать/расширить `apps/parser/internal/commands/songrequest/youtube` в `songrequest` dispatch:
    - `!sr` читает settings.Mode;
    - YouTube: существующий handler;
    - Spotify: вызывает `SpotifySongRequestsService.CreateRequest` через bus/API-gql gRPC (см. существующий `parseCtx.Services.Bus.Api.SongRequestAddToQueue` — нужно добавить bus endpoint или gRPC call).
15. `sr wrong` dispatch по режиму: YouTube — текущий; Spotify — `CancelRequest`.
16. Обновить `apps/events/internal/song_request/song_request.go`: donation flow читает mode и вызывает соответствующий сервис, не генерирует YouTube URL в Spotify режиме.
17. Тесты: parser mode dispatch, donation dispatch, limits shared, validation error messages, `sr wrong` bounds regression.

### Волна 5 — GraphQL + dashboard UI
18. Обновить `apps/api-gql/internal/delivery/gql/schema/song-requests.graphql`:
    - `enum SongRequestMode { YOUTUBE, SPOTIFY }`;
    - `songRequests.mode: SongRequestMode!`;
    - `SongRequestsSettingsOpts.mode: SongRequestMode`;
    - `spotifySongRequestSearch(query: String!): [SpotifySongRequestSearchResult!]!`;
    - `spotifySongRequests(channelID: UUID!): [SpotifySongRequest!]!`;
    - `spotifyIntegration.canQueue: Boolean!` (true when connected + scope present);
    - `spotifyIntegration.missingScopes: [String!]!`.
19. Запустить `bun cli build gql`.
20. Обновить resolvers/mappers для новых полей и статусов.
21. `web/layers/dashboard/api/song-requests.ts`: добавить mode, Spotify search, Spotify request list.
22. `web/layers/dashboard/components/songRequests/settings.vue`: добавить mode selector, disabled-state для Spotify, reconnect prompt, device status indicator.
23. Dashboard song request page: условный список YouTube queue vs Spotify logical queue со статусами.
24. Vue/Vitest tests: settings mode save, disabled Spotify without integration, reconnect flow, search result rendering, status chips.

### Волна 6 — Rollout + monitoring + final checks
25. `bun cli build gql` и `bun cli build` после всех изменений.
26. `bun lint` и `go test ./...` по затронутым пакетам.
27. Добавить runtime metrics: spotify_request_created, spotify_request_status_changed, spotify_device_missing, spotify_rate_limited, spotify_deferred_skip_executed.
28. Обновить AGENTS.md или `docs/` (если существует) описанием режима и ограничений Spotify.
29. Ручной QA на dev-стенде.
30. Feature release plan: additive migration → deploy backend → deploy dashboard → communicate reconnect requirement.

## Todos

- [ ] 1. Add `SongRequestMode` enum to settings entity and repository migration with `YOUTUBE` backfill.
- [ ] 2. Add `spotify_song_requests` table, entity, repository, and mapper with status machine.
- [ ] 3. Update Spotify HTTP client with search/queue/current/queue/devices/skip and typed errors.
- [ ] 4. Add `user-modify-playback-state` to Spotify OAuth scopes and expose capability flags in GraphQL.
- [ ] 5. Implement Spotify song request creation service with validation, device handling, and AddToQueue.
- [ ] 6. Implement `sr wrong` deferred-skip cancellation service and newest-first selection.
- [ ] 7. Implement background polling reconciler for current playback, queue matching, and status transitions.
- [ ] 8. Refactor parser `!sr` and `sr wrong` to dispatch by mode while preserving YouTube path.
- [ ] 9. Update donation song-request event to route by mode.
- [ ] 10. Extend GraphQL schema with mode, Spotify search, and request lifecycle; regenerate and implement resolvers.
- [ ] 11. Add dashboard mode selector, Spotify connection/device status, and Spotify logical queue view.
- [ ] 12. Add backend fixture tests for Spotify client, service, and reconciler; add regression test for `sr wrong` bounds.
- [ ] 13. Add frontend Vitest tests for settings and queue components.
- [ ] 14. Run `bun cli build gql`, `bun cli build`, `bun lint`, and targeted `go test`/`bun test`.
- [ ] 15. Add metrics, docs, and rollout checklist; perform manual end-to-end QA.

## Final verification wave

- [ ] F1. Schema/codegen/build: `bun cli build gql` exits 0, `bun cli build` exits 0, `bun lint` exits 0.
- [ ] F2. Backend tests: `go test ./libs/integrations/spotify/... ./libs/repositories/song_requests_settings/... ./apps/api-gql/internal/services/song_requests/... ./apps/parser/internal/commands/songrequest/... ./apps/events/internal/song_request/...` pass.
- [ ] F3. Frontend tests: `bun run test -- layers/dashboard/features/songRequests` pass.
- [ ] F4. Manual QA: YouTube mode unchanged, Spotify mode OAuth reconnect, `!sr` success, `sr wrong` deferred skip, donation request, device unavailable error, mode switch preserves existing data.

## Commit strategy

- One commit per wave: `feat(song-requests): add spotify mode enum and migration`, `feat(spotify): add queue/control client methods`, `feat(song-requests): add spotify request service and reconciler`, `feat(parser): dispatch song requests by mode`, `feat(events): route donation song requests by mode`, `feat(api-gql): expose song request mode and spotify lifecycle`, `feat(dashboard): add spotify mode selector and queue`, `test(song-requests): add fixtures and regression tests`, `chore(song-requests): metrics and rollout docs`.
- Conventional Commits, no co-authored trailers, author = user identity.

## Success criteria

- YouTube mode behaves identically to pre-change baseline.
- Spotify mode allows `!sr <query/URI>` to enqueue a track in the streamer's Spotify native queue.
- `sr wrong` removes the logical request from Twir and skips it when it becomes current playback.
- Spotify mode requires connected Spotify integration with `user-modify-playback-state` and surfaces a reconnect prompt otherwise.
- Dashboard can switch modes and view Spotify request statuses.
- All backend/frontend tests and build/lint pass.
- Manual QA covers happy path, deferred skip, and error paths (no device, no scope, rate limit).
