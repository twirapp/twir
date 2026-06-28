# Secret Store & Custom Executron Design

## Overview

Implement a per-channel secret store accessible from custom variable scripts via `twir.secrets.get('NAME')`. Build a new NATS-based executron service (Bun + isolated-vm) to replace the external HTTP-based one. Remove Python script support entirely.

## Requirements

| Requirement | Decision |
|---|---|
| Executron runtime | Bun + isolated-vm, JS only |
| Communication | NATS bus (like integrations service) |
| Secret storage | AES-256-GCM encrypted in PostgreSQL |
| Secret UI | Tab on variables page, reveal on click |
| twir global API | `twir.secrets.get(name)`, `twir.channel.id` |
| Python support | Remove entirely |
| Resource limits | 5s timeout, 128MB memory limit |

## Components

### 1. Database: `channels_secrets` table

```sql
CREATE TABLE channels_secrets (
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    value       TEXT NOT NULL,  -- AES-256-GCM encrypted, base64 encoded
    "channelId" TEXT NOT NULL REFERENCES channels(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE("channelId", name)
);
```

Encryption key from `SECRETS_ENCRYPTION_KEY` env var.

### 2. GraphQL API

**Schema** (`apps/api-gql/internal/delivery/gql/schema/secrets.graphql`):

```graphql
extend type Query {
    secrets: [Secret!]! @isAuthenticated @hasAccessToSelectedDashboard
    secretValue(id: UUID!): String! @isAuthenticated @hasAccessToSelectedDashboard
}

extend type Mutation {
    secretCreate(opts: SecretCreateInput!): Secret! @isAuthenticated @hasAccessToSelectedDashboard
    secretUpdate(id: UUID!, opts: SecretUpdateInput!): Secret! @isAuthenticated @hasAccessToSelectedDashboard
    secretDelete(id: UUID!): Boolean! @isAuthenticated @hasAccessToSelectedDashboard
}

type Secret {
    id: UUID!
    name: String!
    description: String
}

input SecretCreateInput {
    name: String! @validate(constraint: "max=100")
    description: String @validate(constraint: "max=500")
    value: String! @validate(constraint: "max=10000")
}

input SecretUpdateInput {
    name: String @validate(constraint: "max=100")
    description: String @validate(constraint: "max=500")
    value: String @validate(constraint: "max=10000")
}
```

- `secrets` query returns name+description (no value)
- `secretValue` query returns decrypted value (for reveal-on-click in UI)
- Value is write-only from list queries

### 3. Go Backend

**New packages:**

| Package | Purpose |
|---|---|
| `libs/entities/secret/entity.go` | Domain entity |
| `libs/repositories/channels_secret/pgx/pgx.go` | Data access |
| `apps/api-gql/internal/services/channels_secret/` | Business logic with encrypt/decrypt |
| `apps/api-gql/internal/delivery/gql/resolvers/secrets.resolver.go` | GraphQL resolvers |

**Encryption**: AES-256-GCM. Each secret gets a random IV stored alongside the ciphertext. Encryption key from config.

### 4. Executron Service

**New app: `apps/executron/`**

```
apps/executron/
├── src/
│   ├── index.ts              # Entry point, NATS connection
│   ├── nats-handler.ts       # NATS request handler
│   ├── executor.ts           # isolated-vm runner
│   ├── secrets-provider.ts   # Fetches secrets from DB, decrypts
│   ├── twir-global.ts        # Builds the twir object
│   └── types.ts              # Shared types
├── Dockerfile
├── package.json
└── tsconfig.json
```

**NATS Request/Response:**

```typescript
// Subject: "executron.execute"
interface ExecuteRequest {
    channelId: string
    language: "javascript"
    code: string
    userId?: string
}

interface ExecuteResponse {
    result: string
    error: string
}
```

**JS Runtime (isolated-vm):**
- Create isolated context with 5s timeout, 128MB memory limit
- Inject `twir` global:
  - `twir.secrets.get(name: string): string | null` — sync, reads from pre-fetched secrets
  - `twir.channel.id: string` — current channel ID
- Extensible for future APIs (`twir.http`, `twir.variables`, etc.)

**Secrets Provider:**
- Receives `channelId` with each request
- Fetches secrets from PostgreSQL, decrypts with AES-256-GCM
- No caching — reads from DB every time

**Integration with api-gql:**
- Remove old HTTP-based `executron` package from `apps/parser/pkg/executron/`
- Add NATS queue in `libs/bus-core/` for executron
- `apps/api-gql/internal/services/variables/` calls executron via NATS instead of HTTP

### 5. Dashboard UI

**Variables page becomes tabbed** (`web/layers/dashboard/pages/dashboard/variables/index.vue`):

```vue
<script setup lang="ts">
import type { PageLayoutTab } from '~~/layers/dashboard/layout/page-layout.vue'
import VariablesList from '~~/layers/dashboard/features/variables/variables.vue'
import SecretsList from '~~/layers/dashboard/features/channels-secret/secrets.vue'

const tabs = computed<PageLayoutTab[]>(() => [
    { title: 'Variables', component: VariablesList, name: 'variables' },
    { title: 'Secrets', component: SecretsList, name: 'secrets' },
])
</script>

<template>
    <PageLayout :tabs="tabs" active-tab="variables">
        <template #title>{{ t('sidebar.variables') }}</template>
    </PageLayout>
</template>
```

**Secrets feature** (`web/layers/dashboard/features/channels-secret/`):
- `secrets.vue` — main component with table
- `composables/use-secrets-table.ts` — table config
- `composables/use-secrets-api.ts` — GraphQL queries/mutations
- `ui/secret-create-button.vue` — create dialog
- `ui/secret-edit-dialog.vue` — edit dialog with masked value + reveal button
- `ui/secret-actions.vue` — row actions (edit, delete)

Secret value field shows `*****` by default. Clicking eye icon calls `secretValue` query.

### 6. Python Removal

| File | Change |
|---|---|
| `apps/api-gql/internal/entity/customvar.go` | Remove `ScriptLanguagePython` constant |
| `apps/api-gql/internal/delivery/gql/mappers/variables.go` | Remove Python mapping |
| `libs/repositories/variables/model/model.go` | Remove `ScriptLanguagePython` constant |
| `apps/api-gql/internal/delivery/gql/schema/variables.graphql` | Remove `PYTHON` from enum |
| `web/layers/dashboard/features/variables/variables-edit.vue` | Remove Python example, language selector |
| `web/layers/landing/components/index/features/features-data.ts` | Remove Python mention |

### 7. Infrastructure

| File | Change |
|---|---|
| `cli/internal/cmds/deploy/deploy.go` | Add `executron` to `serviceImages` and `releaseServices` |
| `.github/workflows/dockerv3.yml` | Add executron to build matrix (`runtime: js`) |
| `docker-compose.stack.yml` | Add executron service definition |
| `docker-compose.dev.yml` | Add executron for local development |

## Execution Order

1. **Database migration** — create `channels_secrets` table
2. **Go backend** — entity, repository, service, GraphQL schema + resolvers
3. **Executron service** — new Bun app with isolated-vm + NATS
4. **Wire executron** — replace HTTP calls with NATS in api-gql
5. **Dashboard UI** — secrets tab, CRUD dialogs
6. **Python removal** — clean up all Python references
7. **Infrastructure** — CLI, CI, docker-compose
