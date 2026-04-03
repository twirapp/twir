# AGENTS.md — cli

Custom Go-based CLI for monorepo orchestration. Central command hub for dev, build, and migrations.

## OVERVIEW

Go CLI tool that orchestrates the entire monorepo. Provides unified commands for development, building, database migrations, and dependency management. Called via Bun scripts from root `package.json`.

## STRUCTURE

```
cli/
├── main.go                  # CLI entry point
├── internal/
│   └── cmds/
│       ├── build/           # Build orchestration
│       │   └── build.go     # Build command implementation
│       ├── dev/             # Dev server management
│       ├── migrations/      # DB migration commands
│       └── deps/            # Dependency installation
└── go.mod
```

## ENTRY POINTS

| Type  | Path                           | Purpose                   |
| ----- | ------------------------------ | ------------------------- |
| Main  | `main.go`                      | CLI bootstrap using cobra |
| Build | `internal/cmds/build/build.go` | Build orchestration logic |

## KEY COMMANDS

```bash
# Run via Bun (recommended)
bun cli <command>

# Or directly
go run ./cli/main.go <command>

# Development
bun cli dev                # Start all services in dev mode

# Building
bun cli build              # Build all apps
bun cli build gql          # Regenerate GraphQL resolvers
bun cli build <app>        # Build specific app

# Migrations
bun cli m create           # Create new migration
bun cli m run              # Run pending migrations

# Dependencies
bun cli deps               # Install binary deps
bun cli deps -skip-node    # Skip Node deps
bun cli deps -skip-go      # Skip Go deps
```

## BUILD ORCHESTRATION

The build command (`internal/cmds/build/build.go`) orchestrates:

1. **JavaScript packages**: `bun --filter` to build workspaces
2. **GraphQL generation**: `gqlgen` for api-gql
3. **Go apps**: Individual `go build` for each app

Build order is determined by dependency graph.

## ANTI-PATTERNS

- **DO NOT** bypass CLI for builds — ensures correct order and env setup
- **DO NOT** add build logic to package.json — belongs in CLI
- **ALWAYS** use CLI for migrations — handles postgres/clickhouse correctly

## COMMANDS REFERENCE

| Command          | Description                        |
| ---------------- | ---------------------------------- |
| `dev`            | Start all services with hot reload |
| `build [target]` | Build apps/packages                |
| `m create`       | Create migration file              |
| `m run`          | Execute pending migrations         |
| `deps`           | Install tool dependencies          |

## NOTES

- Uses cobra for CLI framework
- Respects `.bun-version` for Bun version
- Uses `go.work` for Go module resolution
- Build outputs to `dist/` or app-specific locations
