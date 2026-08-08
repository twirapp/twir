# Dota 2 модуль (dotabod-like) — дизайн

Дата: 2026-07-20
Статус: approved

## Цель

Встроить в Twir систему уровня dotabod для стримеров Dota 2: чат-команды
(`!mmr`, `!wl`, `!lg`, `!gm`, `!np`, `!wp`), автоматические Twitch-ставки на
матч, оверлеи (медаль, W/L счёт, win probability), события в чат (Рошан,
эгис, начало/конец матча). Бесплатно, на Go.

## Решённые вопросы (брейншторм)

- **Live-данные**: гибрид — GSI (Game State Integration, cfg-файл пушит
  HTTP-события из клиента Dota) как основной источник + `go-dota2`
  (paralin) с центральным Steam бот-аккаунтом для постматчевых деталей
  как фоллбек, если Stratz/OpenDota лагают. go-dota2 опционален: стартует
  только если заданы env-креденшелы.
- **Статистика**: Stratz GraphQL API — notable players, win probability;
  OpenDota — матчи/fallback. MMR стример вводит вручную, модуль сам
  пересчитывает по результатам матчей (±25, конфигурируемо).
- **Привязка аккаунта**: Steam OpenID OAuth в дашборде → SteamID64 →
  Dota accountID (младшие 32 бита).
- **Архитектура**: отдельный микросервис `apps/dota` (подход A), обмен
  через bus-core (NATS) с parser, bots, websockets, api-gql.

## Архитектура

### Сервис `apps/dota`

Новый Go-сервис на uber/fx по образцу `apps/events`:

- **GSI HTTP server**: `POST /gsi/:token` — приём событий от клиента Dota.
  Токен — per-channel secret из БД, регенерируемый. Rate-limit по токену.
- **Match state machine**: in-memory + snapshot в Redis:
  `idle → hero_selected → strategy → pre_game → in_game → post_game → idle`.
  Отслеживает: герой стримера, команда (radiant/dire), счёт, аегисы, Рошан,
  win/loss, match_id. Дедупликация по match_id (перезаход в матч).
- **Stats clients**: Stratz GraphQL (notable players, WP) и OpenDota
  (матчи, фоллбек). Общий кэш ответов в Redis.
- **go-dota2 клиент**: центральный Steam аккаунт (env: логин/пароль/2FA
  seed), постматчевые детали через GC как фоллбек.
- **Bus**: публикация событий матча (`dota.MatchStarted`, `dota.MatchEnded`,
  `dota.RoshanKilled`, `dota.AegisPickup`, `dota.GameStateUpdate`),
  ответы на запросы команд (`dota.GetData` — request/reply).

### Команды (apps/parser)

Default commands, toggleable, с кастомными шаблонами как у остальных команд:

- `!mmr` — расчётный MMR (ручной ввод ±25/игра). Мод-команда `!mmr set N`.
- `!wl` — счёт побед/поражений за сессию/стрим.
- `!lg` — последняя игра: герой, KDA, результат, длительность.
- `!gm` — медаль/ранг (из расчётного MMR).
- `!np` — notable players в текущем матче (Stratz).
- `!wp` — win probability (Stratz, данные как у dotabod).

Parser запрашивает `apps/dota` через bus request/reply, ответ кэшируется
(Redis, TTL ~10s). Если матча нет — ответ "нет активного матча".

### Авто-ставки на Twitch

- Триггер: переход в `in_game` (драфт закончен) → создание prediction
  через Helix API (переиспользование кода из
  `apps/parser/internal/commands/predictions/start.go`), скоуп
  `channel:manage:predictions` уже запрашивается.
- Шаблон: "Победит стример? Да/Нет", окно 5 минут (конфигурируемо).
- Триггер `post_game` (win radiant/dire) → resolve по команде стримера
  (accountID из GSI).
- Защита: одна активная ставка на канал; обрыв матча → cancel prediction.

### Оверлеи

В `frontend/overlays` новые страницы:

- `/overlays/dota/:apiKey/medal` — медаль/ранг + расчётный MMR.
- `/overlays/dota/:apiKey/wl` — счёт W/L сессии.
- `/overlays/dota/:apiKey/wp` — win probability текущего матча.

Данные — через существующий `apps/websockets` (pub/sub по channelID),
обновления пушит `apps/dota` при смене состояния. Позиционирование —
стандартное для OBS browser source.

### События в чат

`apps/dota` публикует события в bus → `apps/bots` отправляет сообщения по
шаблонам из настроек модуля (enable/disable и cooldown per event, как у
chat-alerts).

### БД

Миграция: таблица `channels_dota_settings`:

- `id uuid pk`, `channel_id text fk channels(id) unique`
- `enabled bool default false`
- `steam_account_id text null` (из Steam OpenID)
- `gsi_token text` (генерируется, регенерируемый)
- `mmr int default 0`, `mmr_delta int default 25`
- `session_wins int default 0`, `session_losses int default 0`
- `prediction_settings jsonb` (enabled, title template, window seconds)
- `chat_events jsonb` (per-event enabled + templates + cooldowns)
- `commands_settings jsonb` (per-command enabled + templates)
- `created_at`, `updated_at`

### api-gql

- Сервис `internal/services/dota` + GraphQL queries/mutations по паттерну
  kappagen overlay (entity mapper, repository pgx, entity в libs/entities).
- Steam OpenID flow: `GET /auth/steam` → callback → сохранение
  steam_account_id в настройки канала.
- Endpoint генерации/скачивания GSI cfg-файла с токеном канала.

### Дашборд (frontend/dashboard)

Страница модуля Dota:

- Привязка Steam (кнопка OpenID), отображение привязанного аккаунта.
- Скачивание GSI-конфига, инструкция по установке, регенерация токена.
- Ручной ввод/коррекция MMR, mmr_delta, сброс сессии W/L.
- Тоглы команд + шаблоны ответов.
- Настройки авто-ставок (enable, шаблон тайтла, окно).
- Шаблоны и тоглы чат-событий.

## Edge cases

- Dota закрыта/оффлайн → команды отвечают "нет активного матча" из кэша.
- GSI шлёт ~1 req/с на канал → rate-limit по токену, валидация auth до
  парсинга body.
- Перезаход в матч → дедупликация по match_id, ставка не дублируется.
- GSI-события без матча (меню) → state machine в `idle`, события не
  публикуются.
- Стример без привязанного Steam → команды/ставки/оверлеи неактивны,
  команда отвечает "аккаунт не привязан".
- go-dota2 без env-креденшелов → фоллбек отключён, сервис работает на
  Stratz/OpenDota.

## Тестирование

- Юнит-тесты стейт-машины на фикстурах GSI JSON (записанные потоки
  реальных событий).
- Моки Stratz/OpenDota клиентов.
- Интеграционный тест GSI endpoint (auth, rate-limit, дедупликация).
- Тесты расчёта MMR и W/L сессии.
- Фронтенд — по существующим паттернам дашборда/оверлеев.

## YAGNI (не входит в v1)

- OBS scene switching, minimap overlay, bets на hero pick, smurf detection,
  pro-player shoutouts, Dota-чат через GC, несколько Dota-аккаунтов на
  канал, турбо/unranked-специфика MMR.
