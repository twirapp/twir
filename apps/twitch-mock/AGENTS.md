# AGENTS.md — apps/twitch-mock

Local Twitch API mock server for development (Go). Not deployed to production.

## OVERVIEW

Mimics Twitch OAuth2 + Helix API + EventSub WebSocket so Twir can be developed/tested without a real Twitch account, OAuth flows, or rate limits. See `README.md` for full usage.

## STRUCTURE

```
apps/twitch-mock/
├── cmd/             # main.go entry
├── app/             # HTTP/WS server setup
├── internal/        # Endpoint handlers
└── Dockerfile
```

## KEY FACTS

- Enable via `TWITCH_MOCK_ENABLED=true` in root `.env` (see `.env.mock.example`); started by `docker-compose.dev.yml`.
- Ports: `7777` OAuth2+Helix, `8081` EventSub WebSocket, `3333` Admin UI (`http://localhost:3333/admin`).
- Fake users: broadcaster `12345`/`mockstreamer`, bot `67890`/`mockbot`.
- Events triggerable via Admin UI or `POST http://localhost:3333/admin/trigger/<event.type>`.

## NOTES

- NOT mocked: Twitch IRC/chat protocol, player iframe, niche Helix endpoints (may 404/empty).
- If login unexpectedly hits real Twitch — `TWITCH_MOCK_ENABLED` is missing/false.
