# AGENTS.md — libs/types

Shared TypeScript type definitions.

## OVERVIEW

Common TypeScript types used across frontend applications and libraries. Provides shared interfaces for API responses, entities, and domain models.

## STRUCTURE

```
libs/types/
├── src/
│   └── *.ts                 # Type definitions
├── package.json
├── tsconfig.json
└── ...
```

## USAGE

```typescript
import type { User, Channel } from '@twir/types';

const user: User = { ... };
```

## NOTES

- Shared by dashboard, overlays, web
- Keep in sync with Go entities
- No runtime code — types only
