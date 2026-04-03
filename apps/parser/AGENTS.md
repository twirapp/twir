# AGENTS.md — apps/parser

Chat command and variable parser.

## OVERVIEW

Parses chat commands and variables in messages. Executes command logic and returns responses. Supports custom commands, timers, and variable substitution.

## STRUCTURE

```
apps/parser/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── commands/            # Command implementations
│   │   ├── stats/           # Stats commands
│   │   ├── dota/            # Dota commands
│   │   └── ...
│   ├── variables/           # Variable processors
│   │   ├── user/            # User variables
│   │   └── ...
│   └── ...
├── locales/                 # i18n files
│   └── en/
│       ├── commands/        # Command descriptions
│       └── variables/       # Variable descriptions
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type | Path          | Purpose                  |
| ---- | ------------- | ------------------------ |
| Main | `cmd/main.go` | Parser service bootstrap |
| gRPC | —             | Command execution API    |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build parser
```

## DEPENDENCIES

- Postgres (commands, variables)
- Redis (caching)
- emotes-cacher (emote lookups)
- External APIs (Steam, Dota, etc.)

## ANTI-PATTERNS

- **TODO** in codebase: refactor parsectx to new chat message struct
- **TODO**: create typed response subtype

## NOTES

- Supports 50+ built-in commands
- Variable syntax: `${variableName}`
- Extensible command system
