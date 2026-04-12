# Testing Guide

Этот документ описывает как тестировать различные части Twir локально с помощью встроенного
Twitch mock-сервера — без реального Twitch-аккаунта.

## Содержание

- [Быстрый старт](#быстрый-старт)
- [Fake-пользователи](#fake-пользователи)
- [Тестирование событий через Admin UI](#тестирование-событий-через-admin-ui)
- [Тестирование чат-сообщений и команд бота](#тестирование-чат-сообщений-и-команд-бота)
- [Тестирование EventSub-событий](#тестирование-eventsub-событий)
- [Тестирование OAuth / авторизации](#тестирование-oauth--авторизации)
- [Тестирование Helix API](#тестирование-helix-api)
- [Что НЕ замокано](#что-не-замокано)
- [Порты и сервисы](#порты-и-сервисы)

---

## Быстрый старт

```bash
# 1. Скопировать mock-конфиг
cp .env.mock.example .env

# 2. Поднять инфраструктуру (включая mock-сервер)
docker compose -f docker-compose.dev.yml up -d

# 3. Запустить все сервисы
bun dev

# 4. Открыть дашборд
open http://localhost:3010
# → Залогиниться как MockStreamer (кнопка Login → автоматический mock OAuth)

# 5. Открыть Admin UI для триггеров
open http://localhost:3333/admin
```

---

## Fake-пользователи

Mock-сервер предоставляет двух фиксированных пользователей:

| Роль        | ID    | Login        | Display Name |
| ----------- | ----- | ------------ | ------------ |
| Broadcaster | 12345 | mockstreamer | MockStreamer |
| Bot         | 67890 | mockbot      | MockBot      |

При авторизации через дашборд вы всегда входите как **MockStreamer** (ID=12345).

---

## Тестирование событий через Admin UI

Самый простой способ — открыть **http://localhost:3333/admin**.

Там есть готовые формы для всех основных событий с примерами payload:

| Событие                                               | Что тестирует                   |
| ----------------------------------------------------- | ------------------------------- |
| `channel.follow`                                      | Алерты на фолловы, команды      |
| `channel.subscribe`                                   | Алерты на подписки              |
| `channel.cheer`                                       | Алерты на битсы                 |
| `channel.raid`                                        | Алерты на рейды                 |
| `channel.chat.message`                                | Команды бота, триггеры, фильтры |
| `channel.ban` / `channel.unban`                       | Логика модерации                |
| `channel.poll.begin` / `channel.poll.end`             | Голосования                     |
| `channel.prediction.begin` / `.end`                   | Предсказания                    |
| `channel.channel_points_custom_reward_redemption.add` | Redemption-триггеры             |
| `stream.online` / `stream.offline`                    | Онлайн/офлайн события           |

Нажимаете **Trigger** — событие уходит по WebSocket в `apps/eventsub`, дальше по шине в нужный сервис.

---

## Тестирование чат-сообщений и команд бота

### Как это работает

```
Admin UI → POST /admin/trigger/channel.chat.message
  → WebSocket (port 8081) → apps/eventsub
  → HandleChannelChatMessage → шина сообщений
  → apps/parser — разбирает команды
  → apps/bots — формирует ответ
```

### Через Admin UI

Откройте **http://localhost:3333/admin**, найдите карточку `channel.chat.message`, измените поле
`message.text` на нужную команду и нажмите **Trigger**.

Пример payload для команды `!followage`:

```json
{
	"broadcaster_user_id": "12345",
	"broadcaster_user_login": "mockstreamer",
	"broadcaster_user_name": "MockStreamer",
	"chatter_user_id": "99999",
	"chatter_user_login": "viewer",
	"chatter_user_name": "Viewer",
	"message_id": "msg-test-1",
	"message": {
		"text": "!followage",
		"fragments": [{ "type": "text", "text": "!followage" }]
	},
	"message_type": "text"
}
```

### Через curl

```bash
curl -s -X POST http://localhost:3333/admin/trigger/channel.chat.message \
  -H "Content-Type: application/json" \
  -d '{
    "broadcaster_user_id": "12345",
    "broadcaster_user_login": "mockstreamer",
    "broadcaster_user_name": "MockStreamer",
    "chatter_user_id": "99999",
    "chatter_user_login": "viewer",
    "chatter_user_name": "Viewer",
    "message_id": "msg-test-1",
    "message": {
      "text": "!команда",
      "fragments": [{ "type": "text", "text": "!команда" }]
    },
    "message_type": "text"
  }'
```

### Где смотреть ответ бота

> ⚠️ IRC (реальный чат Twitch) **не замокан**. Бот обработает событие, но его ответ пойдёт в
> настоящий Twitch IRC — в dev-окружении вы его там не увидите.

Ответ бота можно наблюдать в **логах** сервиса `bots`:

```bash
# Смотреть только логи bots в реальном времени
bun dev 2>&1 | grep -i "\[bots\]"

# Или через отдельный терминал если bots запущен отдельно
go run ./apps/bots/cmd/main.go 2>&1 | grep -i "send\|message\|command\|chat"
```

Искать строки типа:

- `sending message` — бот собирается что-то отправить
- `command processed` — команда обработана
- `triggered` — сработал триггер

---

## Тестирование EventSub-событий

Все события из Admin UI идут через WebSocket EventSub (port 8081). Flow:

1. `apps/eventsub` подключается к `ws://localhost:8081/ws` при старте
2. Admin UI POSTит на `http://localhost:3333/admin/trigger/{event_type}`
3. Mock-сервер бродкастит JSON-уведомление всем подключённым WS-клиентам
4. `apps/eventsub` получает и роутит в нужный handler

Проверить что eventsub подключён:

```bash
curl -s http://localhost:7777/helix/eventsub/subscriptions | jq .
# Должен вернуть список активных подписок (заполняется при старте eventsub)
```

---

## Тестирование OAuth / авторизации

Mock-сервер эмулирует полный OAuth2 Authorization Code Flow.

### Ручная проверка authorize

```bash
curl -v "http://localhost:7777/oauth2/authorize?client_id=mock-client-id&redirect_uri=http://localhost:3009/auth/callback&response_type=code&scope=openid"
# → HTTP 302, Location: http://localhost:3009/auth/callback?code=mock_code_1234567890
```

### Ручная проверка token exchange

```bash
curl -s -X POST http://localhost:7777/oauth2/token \
  -d "client_id=mock-client-id&client_secret=mock-client-secret&code=mock_code_1234567890&grant_type=authorization_code&redirect_uri=http://localhost:3009/auth/callback"
# → {"access_token":"mock-user-token","expires_in":99999999,...}
```

### Проверка токена

```bash
curl -s -H "Authorization: OAuth mock-user-token" http://localhost:7777/oauth2/validate
# → {"user_id":"12345","login":"mockstreamer",...}
```

---

## Тестирование Helix API

Mock-сервер отвечает на все основные Helix-эндпоинты. Для проверки используйте app-токен:

```bash
# Получить пользователя по ID
curl -s -H "Authorization: Bearer mock-app-token" \
  "http://localhost:7777/helix/users?id=12345" | jq .

# Получить информацию о стриме
curl -s -H "Authorization: Bearer mock-app-token" \
  "http://localhost:7777/helix/streams?user_id=12345" | jq .

# Список каналов
curl -s -H "Authorization: Bearer mock-app-token" \
  "http://localhost:7777/helix/channels?broadcaster_id=12345" | jq .

# Активные EventSub-подписки
curl -s -H "Authorization: Bearer mock-app-token" \
  "http://localhost:7777/helix/eventsub/subscriptions" | jq .
```

Неизвестные эндпоинты возвращают `{"data":[],"total":0,"pagination":{}}`.

---

## Что НЕ замокано

| Часть                  | Статус  | Причина                                             |
| ---------------------- | ------- | --------------------------------------------------- |
| Twitch IRC / Chat      | ❌      | Бот подключается к `irc.chat.twitch.tv` напрямую    |
| Twitch Player / Embed  | ❌      | Iframe загружается с twitch.tv                      |
| PubSub (legacy)        | ❌      | Не используется в новых фичах                       |
| Редкие Helix-эндпоинты | ⚠️ stub | Возвращают пустой `{"data":[]}`, не вызывают ошибок |

---

## Порты и сервисы

| Порт | Сервис         | URL                         |
| ---- | -------------- | --------------------------- |
| 3010 | Nuxt (web)     | http://localhost:3010       |
| 3009 | api-gql        | http://localhost:3009       |
| 3333 | Mock Admin UI  | http://localhost:3333/admin |
| 7777 | Mock HTTP API  | http://localhost:7777       |
| 8081 | Mock WebSocket | ws://localhost:8081/ws      |

Подробная документация по mock-серверу: [`apps/twitch-mock/README.md`](apps/twitch-mock/README.md)
