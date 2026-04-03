# AGENTS.md — libs/repositories

Data access layer for PostgreSQL using pgx. Repository pattern implementation for all backend services.

## OVERVIEW

Shared data access layer using pgx (PostgreSQL driver). Provides type-safe database access for Go services. Each domain has its own repository package with pgx implementation.

## STRUCTURE

```
libs/repositories/
├── {entity}/                # One directory per entity. But is it legacy. Instead write entities in `/libs/entities`
│   ├── pgx/
│   │   └── pgx.go          # pgx implementation
│   └── model.go            # Repository model (if needed)
├── go.mod
└── ...

# Example:
channels/
├── pgx/
│   └── pgx.go              # ChannelsRepository implementation
└── model.go                # Channel model
```

## CONVENTIONS

### Repository Implementation

```go
package pgx

type ChannelsRepository struct {
    db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *ChannelsRepository {
    return &ChannelsRepository{db: db}
}

func (r *ChannelsRepository) GetByID(ctx context.Context, id string) (channels.Channel, error) {
    // Implementation
}
```

### Model Pattern

Models include nil-checking for empty results:

```go
type Channel struct {
    ID        string
    Name      string
    // ...
    isNil     bool
}

func (c Channel) IsNil() bool {
    return c.isNil
}

var Nil = Channel{isNil: true}
```

### Entity Mapping

Data flows through layers:

```
DB Row → Model → Entity (libs/entities) → DTO (GraphQL/HTTP)
```

Example:

```go
func (r *ChannelsRepository) GetByID(ctx context.Context, id string) (channel.Channel, error) {
    var model Channel
    err := r.db.QueryRow(ctx, query, id).Scan(&model.ID, &model.Name)
    if err != nil {
        return channel.Nil, err
    }

    // Map to entity
    return channel.Channel{
        ID:   model.ID,
        Name: model.Name,
    }, nil
}
```

## ANTI-PATTERNS

- **NEVER** use GORM or other ORMs — use pgx only
- **NEVER** return models directly — map to entities
- **NEVER** use `sql.Null*` types — use pgx native types
- **ALWAYS** accept `context.Context` as first parameter
- **ALWAYS** propagate context to pgx calls

## DEPENDENCIES

- `libs/entities/{domain}` — Domain entities
- `libs/gomodels` — Legacy model definitions (migration in progress)
- `libs/cache` — Caching layer (optional)

## TESTING

```bash
# Run repository tests if they exist
go test ./...

# With coverage
go test -cover ./...
```

## NOTES

- Uses pgx v5
- Connection pool from `libs/baseapp`
- Follow PostgreSQL best practices (see `skills/supabase-postgres-best-practices`)
- Migrations live in `libs/migrations`
