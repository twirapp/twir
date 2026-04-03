# AGENTS.md — libs/migrations

Database migration management.

## OVERVIEW

PostgreSQL and ClickHouse migrations. CLI commands for creating and running migrations across environments.

## STRUCTURE

```
libs/migrations/
├── cmd/                     # Migration CLI
├── migrations/              # Migration files
│   ├── postgres/           # Postgres migrations
│   └── clickhouse/         # ClickHouse migrations
├── Dockerfile
├── go.mod
└── ...
```

## KEY COMMANDS

```bash
# Create new migration
bun cli m create --name <name> --db postgres|clickhouse --type sql|go

# Run migrations
bun cli m run
```

## CONVENTIONS

### Migration File Naming

```
{timestamp}_{name}.sql
{timestamp}_{name}.go
```

### SQL Migration

```sql
-- +migrate Up
CREATE TABLE example (...);

-- +migrate Down
DROP TABLE example;
```

## NOTES

- Supports SQL and Go migrations
- Postgres and ClickHouse
- Version tracked in schema_migrations table
