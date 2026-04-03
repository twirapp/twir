# AGENTS.md — apps/integrations

JavaScript service for 3rd party integrations.

## OVERVIEW

Node.js/Bun service handling integrations with external services (Discord, Spotify, etc.). Manages OAuth flows and API polling.

## STRUCTURE

```
apps/integrations/
├── src/
│   └── ...                  # Integration handlers
├── package.json
├── tsconfig.json
└── Dockerfile
```

## ENTRY POINTS

| Type | Path           | Purpose       |
| ---- | -------------- | ------------- |
| Main | `src/index.ts` | Service entry |

## KEY COMMANDS

```bash
# Development
bun dev

# Build
bun run build

# Or via CLI
bun cli build integrations
```

## DEPENDENCIES

- Postgres (integration tokens)
- External APIs (Discord, Spotify, etc.)

## NOTES

- Bun runtime for JS execution
- Separate from Go services for flexibility
