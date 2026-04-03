# AGENTS.md — libs/api

GraphQL API client for frontend.

## OVERVIEW

Shared GraphQL client configuration and utilities. Wraps urql with project-specific settings for authentication and error handling.

## STRUCTURE

```
libs/api/
├── src/
│   └── *.ts                 # API client, urql config
├── package.json
├── tsconfig.json
└── ...
```

## USAGE

```typescript
import { createClient } from "@twir/api";

const client = createClient({
	url: "/api/query",
});
```

## NOTES

- Used by dashboard and overlays
- Configures urql client
- Handles auth tokens
