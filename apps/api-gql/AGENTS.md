# AGENTS.md — apps/api-gql

GraphQL API gateway for Twir. Central orchestrator connecting all backend services.

## OVERVIEW

Main GraphQL API service serving dashboard, overlays, and web clients. Handles authentication, real-time subscriptions via WebSockets, and aggregates data from all microservices.

## STRUCTURE

```
apps/api-gql/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── delivery/
│   │   ├── gql/             # GraphQL resolvers
│   │   │   ├── resolvers/   # Root resolvers
│   │   │   ├── schema/      # GraphQL schema files
│   │   │   └── mappers/     # Entity → DTO mappers
│   │   └── http/            # HTTP handlers (REST)
│   ├── services/            # Business logic layer
│   └── server/              # Server setup (HTTP + WebSocket)
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type               | Path             | Purpose                         |
| ------------------ | ---------------- | ------------------------------- |
| Main               | `cmd/main.go`    | Service bootstrap               |
| GraphQL Playground | `/graphql`       | Interactive query explorer      |
| HTTP API           | `/api/*`         | REST endpoints                  |
| WebSocket          | `/subscriptions` | Real-time GraphQL subscriptions |

## KEY COMMANDS

```bash
# Run locally (requires infra: postgres, redis)
bun cli build gql          # Regenerate resolvers after schema changes
go run ./cmd/main.go       # Start dev server on :3009

# Build for production
bun cli build api-gql
```

## CONVENTIONS

### GraphQL Schema Changes

1. Edit `.graphql` files in `internal/delivery/gql/schema/`
2. Run `bun cli build gql` to regenerate resolvers
3. Implement resolver logic in `internal/delivery/gql/resolvers/`

### Data Flow

```
DB (pgx) → Model → Entity (libs/entities) → DTO (GraphQL) → Client
```

### Services

- One service per domain in `internal/services/`
- Services depend on repositories (from `libs/repositories`)
- Return entities, not models directly

## ANTI-PATTERNS

- **NEVER** use `context.TODO()` in production code — propagate request context
- **NEVER** use GORM — use pgx via `libs/repositories`
- **DO NOT** edit generated resolver files — modify schema and regenerate
- **DO NOT** return database models directly — always map to entities. Always write/use entities from `/libs/entities`

## DEPENDENCIES

- Postgres (migrations in `libs/migrations`)
- Redis (caching)
- Other apps via gRPC (tokens, events, parser)

## PORTS

| Service   | Port | Protocol      |
| --------- | ---- | ------------- |
| HTTP API  | 3009 | HTTP/GraphQL  |
| WebSocket | 3009 | WS (upgraded) |

## NOTES

- Uses gqlgen via own cli tool for GraphQL code generation
- WebSocket subscriptions handled via separate upgrade path
- CORS configured for dashboard (:3006) and web (:3000) origins in dev
