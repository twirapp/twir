# Multi-platform header stats — план реализации

> Дизайн-мокап (4 варианта, открывается в браузере): `.omo/plans/multiplatform-header-mockup.html`
> Реализуем варианты **2, 3 и 4(модифицированный)** с переключателем для пользователя.

## Проблема

Хедер дашборда (`web/layers/dashboard/layout/header/header.vue`) показывает стату одного стрима
(Twitch-only). У юзера может быть несколько подключённых платформ (Twitch, Kick, VK) — у каждой
свой title/category, свои viewers/followers/uptime.

## Что есть сегодня (разведка 2026-07-29)

### Backend

- `channels_streams` — **per-platform**: `ON CONFLICT ("userId", platform) DO UPDATE`, параллельные
  стримы Twitch+Kick хранятся. Репо: `libs/repositories/streams` (`GetByChannelID(ctx, id, platform)`,
  `Save`, `Update`, `DeleteByChannelID` — всё platform-scoped).
- Писатели стримов: Twitch — `apps/eventsub/internal/handler/stream_online.go` / `stream_offline.go`;
  Twitch+Kick поллинг — `apps/scheduler/internal/timers/streams.go`; Kick metadata update —
  `apps/eventsub/internal/kick/handlers.go`. **VK писателя стримов нет** (isLive всегда false).
- Биндинги платформ: `channel_platforms` (`libs/repositories/channel_platforms`, entity
  `libs/entities/channel` → `Bindings`, helpers `Binding(platform)` / `TwitchBinding()`).
  Поля: `ChannelID, Platform, UserID, PlatformChannelID, Enabled, BotUserID, BotConfig`.
- Текущая стата: `apps/api-gql/internal/services/dashboard/dashboard.go` — читает только Twitch-стрим;
  followers — Twitch API `GetChannelFollows`; subs — Twitch API (cache); usedEmotes — ClickHouse
  `channels_emotes_usages` (уже умеет фильтр по platform); messages — Redis
  `stream:parsedMessages:<streamID>` (per stream row → per-platform); requestedSongs — Postgres,
  channel-level (не platform-scoped).
- GraphQL: `apps/api-gql/internal/delivery/gql/schema/dashboard.graphql` —
  `type DashboardStats { categoryId categoryName viewers startedAt title chatMessages followers
  usedEmotes requestedSongs subs }`, подписка `dashboardStats` = поллинг раз в 5с
  (`resolvers/dashboard.resolver.go`).
- Редактирование title/category: `twitchSetChannelInformation(title, categoryId)` — Twitch-only,
  без platform-аргумента.

### Frontend

- `useRealtimeDashboardStats()` — `web/layers/dashboard/api/dashboard.ts`, поля см. выше, без platform.
- Подключённые платформы фронт уже знает: `channelPlatformBindings { id platform enabled platformLogin
  platformDisplayName platformAvatar capabilities { name } }` —
  `web/layers/dashboard/features/channel-platforms/api.ts`.
- Редактор title: `web/layers/dashboard/layout/stream-info-editor.vue` — захардкожен на Twitch
  (`isTwitchDashboard` guard), мутация в `web/layers/dashboard/api/twitch.ts`.
- Хедер уже имеет: конфиг виджетов в localStorage (`twirHeaderStatsWidgetsv1`), edit-mode с
  drag&drop, кнопку Edit. Uptime считается локально тикером (`useIntervalFn` + `intervalToDuration`).

## Дизайны (зафиксированы, мокап см. в HTML)

### Layout `rows` (мокап, вариант 2) — «платформы внутри виджетов»

- Viewers/Followers виджеты многострочные: строка `иконка-платформы + число` на каждую платформу
  (offline — приглушена).
- Title-виджет → поповер: на каждую платформу строка с иконкой, title + category, свой карандаш
  (disabled/`canEditInfo=false` для VK).
- Глобальные виджеты (Messages/Subs/Emotes/Songs) — как сейчас, существующий widget-config
  (reorder/hide) продолжает работать.

### Layout `aggregate` (мокап, вариант 3) — «агрегат + поповер» — **DEFAULT**

- Навбар: «33 Viewers», «686 Followers» (сумма по live-платформам) + стек мини-иконок платформ +
  chevron. Followers без источника данных (VK) не входят в сумму — в поповере показываем `—`.
- Клик → поповер: таблица Platform | Viewers | Followers | Uptime (offline — dimmed + бейдж),
  ниже — per-platform строки title/category с редактированием.
- Title-виджет: показывает «primary» платформу (первая live, иначе первая enabled), иконка-платформы
  как бейдж; клик → тот же поповер.
- Глобальные виджеты — как сейчас.

### Layout `grid` (мокап, вариант 4, **модифицированный по запросу**)

- **Без глобальных статистик вообще.** Навбар = только карточки платформ.
- Сетка **2 столбца** (`grid-cols-2`): при 3 платформах — 2 в первом столбце, 1 во втором
  (vertical flow, пока ок — зафиксировано пользователем).
- Карточка платформы: иконка + LIVE/Offline бейдж, title (truncate) + inline карандаш
  (если `canEditInfo`), category, uptime, viewers, followers.
- Справа остаются CommandMenu/BotStatus/Profile.

## Архитектура

### Backend (Phase 1)

1. `dashboard.graphql`: добавить
   ```graphql
   type DashboardPlatformStats {
     platform: Platform!
     isLive: Boolean!
     title: String
     categoryId: ID
     categoryName: String
     viewers: Int
     followers: Int      # null — нет источника (VK)
     startedAt: Time
     chatMessages: Int!
     usedEmotes: Int!
     canEditInfo: Boolean! # twitch+kick: true, vk: false
   }
   # в DashboardStats добавить:
   platforms: [DashboardPlatformStats!]!
   ```
   Существующие поля `DashboardStats` не трогаем (обратная совместимость; глобальные виджеты
   продолжают жить на них).
2. `services/dashboard/dashboard.go`: загрузить enabled `channel_platforms` биндинги канала;
   для каждой платформы: stream row (`streams.GetByChannelID`), viewers/title/category/startedAt
   из неё; messages — Redis counter по streamID; usedEmotes — ClickHouse с platform-фильтром
   (уже есть); followers — Twitch: существующий `GetChannelFollows`; Kick: проверить followers
   в Kick API (skill `kick-platform`, `GET /public/v1/channels`), если нет — null; VK — null.
3. Entity + mapper: `internal/entity/dashboard.go`, `internal/delivery/gql/mappers/dashboard.go`.
4. После схемы: `bun cli build gql`, потом `bun run graphql-codegen` в `web/`.

### Frontend (Phase 2)

Переиспользование — ядро дизайна кода:

```
web/layers/dashboard/layout/header/
  header.vue                       # shell: рендерит активный layout по useHeaderLayout()
  header-bot-status.vue            # без изменений
  header-profile.vue               # без изменений
  composables/
    use-header-layout.ts           # 'rows' | 'aggregate' | 'grid' + localStorage
                                   # (`twirHeaderStatsLayoutV1`) + дефолт 'aggregate'
    use-platform-stats.ts          # нормализация stats.platforms: sortedPlatforms
                                   # (live first), totalViewers, totalFollowers (skip null),
                                   # primaryPlatform, per-platform uptime ticker
                                   # (обобщение существующего тикера из header.vue)
  ui/
    platform-icon.vue              # simple-icons:* + brand color + optional live-dot
    stat-widget.vue                # shell виджета (существующие .header-widget стили)
    platform-stats-rows.vue        # строки «иконка+число» (V2 виджеты) И таблица
                                   # разбивки (V3 поповер) — один компонент, два size
    platform-title-row.vue         # title+category+карандаш на платформу
                                   # (V2/V3 поповеры, V4 карточки)
  layouts/
    stats-layout-rows.vue          # V2
    stats-layout-aggregate.vue     # V3
    stats-layout-grid.vue          # V4: grid-cols-2, только карточки платформ
```

- Глобальные виджеты (uptime/messages/subs/emotes/songs) + существующий widget-config + edit-mode
  выделяются в `global-stats-widgets.vue`, используется в `rows` и `aggregate`; в `grid` не
  рендерится.
- Переключатель: кнопка `lucide:layout-grid` рядом с существующей кнопкой Edit → Popover с
  радио-списком из 3 вариантов (название + однострочное описание), выбор пишется в localStorage,
  применяется мгновенно. i18n-ключи `dashboard.statsLayouts.{rows,aggregate,grid}` во все локали
  (en/ru/uk/de/es/pt/sk/ja).
- Иконки платформ: `simple-icons:twitch` (#9146FF), `simple-icons:kick` (#53FC18),
  `simple-icons:vk` (#0077FF) через `<Icon />` (конвенция проекта).

### Редактирование title/category per-platform (Phase 3)

1. `stream-info-editor.vue` принимает `platform` prop; guard меняется с `isTwitchDashboard` на
   `canEditInfo(platform)`.
2. Новая мутация `channelSetStreamInformation(platform: Platform!, title: String, categoryId: String)`:
   сервис диспатчит по платформе → Twitch: существующий путь; Kick: Kick API PATCH channel
   (проверить scopes по skill `kick-platform`); VK: ошибка `unsupported`.
   Старая `twitchSetChannelInformation` остаётся (deprecate позже).
3. Поиск категорий per-platform (Twitch categories search существует; Kick categories — Kick API).

## Фазы и границы

| Phase | Содержание | Критерий готовности |
|---|---|---|
| 1 | Backend: `platforms` в `dashboardStats` | QA-1 пройден |
| 2 | Frontend: 3 layout + переключатель | QA-2 пройден |
| 3 | Per-platform редактирование title/category | QA-3 пройден |

MVP-разрез, если нужно быстрее: Phase 1 (только Twitch+Kick followers, VK null) + Phase 2 без
редактирования (карандаш скрыт до Phase 3).

## QA-сценарии (исполнимые, по фазам)

### QA-1 (Phase 1, backend)

1. **Go unit-тест сервиса** (инструмент: `go test ./apps/api-gql/internal/services/dashboard/...`).
   Собрать `GetDashboardStats` с моками репозиториев (паттерн моков — как в соседних сервисах):
   - Кейс A: у канала enabled-биндинги Twitch+Kick, в `channels_streams` есть live-строки обеих
     платформ → ожидаем `platforms` длиной 2+, у Twitch `isLive=true, viewers/title/categoryName/
     startedAt` из строки стрима, у Kick аналогично; `chatMessages`/`usedEmotes` per-platform
     из моков Redis/ClickHouse.
   - Кейс B: enabled VK-биндинг без строки стрима → элемент с `platform=VK_VIDEO_LIVE,
     isLive=false, viewers=null, followers=null, canEditInfo=false`.
   - Кейс C: канал вообще без стримов → `platforms` содержит все enabled-биндинги с `isLive=false`,
     существующие top-level поля не изменили поведение (Twitch-фолбэк как сейчас).
   Ожидаемый результат: все кейсы зелёные.
2. **Ручная проверка подписки** (инструмент: `bun dev` + браузер, DevTools → Network →
   GraphQL subscription `dashboardStats`, дашборд через Caddy URL из `.env` `SITE_BASE_URL`,
   по умолчанию `http://localhost:3005`). Шаги: открыть дашборд канала с Twitch-биндингом →
   в payload подписки проверить поле `platforms`. Ожидаемый результат: массив содержит объект с
   `platform: "TWITCH"` и непустыми `title/viewers/followers`, старые поля (`title`, `viewers`,
   `followers` top-level) присутствуют как раньше.
3. **Регрессия codegen**: после `bun cli build gql` и `bun run graphql-codegen` (в `web/`) —
   `bun cli build` завершается с exit 0, типы фронта обновлены.

### QA-2 (Phase 2, frontend)

1. **Vitest composables** (инструмент: `cd web && bun run test -- layers/dashboard`):
   - `use-platform-stats.spec.ts`: суммирование viewers (только live), `totalFollowers` пропускает
     `null` (VK не входит в сумму), `primaryPlatform` = первая live, иначе первая enabled;
     сортировка live-first.
   - `use-header-layout.spec.ts`: дефолт `'aggregate'` при пустом localStorage, запись/чтение
     ключа `twirHeaderStatsLayoutV1`, игнорирование невалидного значения из localStorage.
   Ожидаемый результат: все спеки зелёные.
2. **Ручной UI-прогон через Playwright** (инструмент: Playwright MCP, дашборд через Caddy URL из
   `.env` `SITE_BASE_URL`; НЕ localhost:3000 напрямую):
   - Шаг 1: открыть дашборд → дефолтный layout `aggregate`: виджет Viewers показывает сумму
     зрителей и стек иконок платформ; клик открывает поповер с таблицей Platform/Viewers/
     Followers/Uptime.
   - Шаг 2: кнопка переключателя (иконка layout рядом с Edit) → выбрать `rows` → виджет Viewers
     стал многострочным (иконка+число на платформу), title-виджет открывает поповер со строками
     платформ.
   - Шаг 3: выбрать `grid` → глобальные виджеты (Messages/Subs/Emotes/Songs/Uptime) исчезли,
     остались только карточки платформ в `grid-cols-2` (при 3 платформах: 2+1), VK-карточка
     dimmed с Offline.
   - Шаг 4: reload страницы → выбранный layout сохранился (localStorage).
   Ожидаемый результат: каждый шаг соответствует мокапу
   `.omo/plans/multiplatform-header-mockup.html` (варианты 2/3/4), нет горизонтального overflow
   в хедере на 1500px.
3. **Lint/типы**: `bun lint` и сборка web чистые; существующие спеки дашборда не сломаны
   (`bun run test` полный прогон зелёный).

### QA-3 (Phase 3, редактирование)

1. **Go unit-тест диспатча мутации** (`go test` resolver/service): `channelSetStreamInformation`
   с `platform=TWITCH` → вызван Twitch-путь (мок); `platform=KICK` → Kick-путь; `platform=
   VK_VIDEO_LIVE` → ошибка `unsupported platform`. Ожидаемый результат: тесты зелёные.
2. **Ручной UI-прогон (Playwright MCP, Caddy URL)**:
   - Twitch: в layout `rows` открыть title-поповер → карандаш у Twitch → редактор открывается с
     текущими title/category → сохранить новый title → в течение ~5с (период поллинга) title в
     хедере обновился.
   - Kick (если в dev окружении есть Kick-креды; иначе только пункт 1 + код-ревью): то же самое
     для Kick-строки.
   - VK: карандаш у VK-строки disabled, hover показывает тултип «редактирование недоступно»,
     клик ничего не открывает.
   Ожидаемый результат: мутация уходит с корректным `platform`, ошибок в консоли нет.
3. **Регрессия**: старая мутация `twitchSetChannelInformation` продолжает работать
   (существующие вызовы не сломаны) — покрыто компиляцией + существующими спеками.

## Открытые вопросы (решения зафиксированы, можно пересмотреть)

- **Followers Kick/VK**: Kick — проверить API; VK — нет источника → null/`—`. Решено: null, не блокер.
- **VK live-стата**: писателя стримов нет → VK всегда offline в хедере. Решено: показываем биндинг
  как offline, отдельная задача если понадобится.
- **Subs per-platform**: сейчас Twitch-only и остаётся глобальным виджетом. Не трогаем.
- **Requested songs / messages**: остаются глобальными (channel-level). В `rows`/`aggregate` — как
  сейчас; в `grid` отсутствуют по дизайну.
