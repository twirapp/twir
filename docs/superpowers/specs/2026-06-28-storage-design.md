# Storage Feature Design

## Overview

A per-channel key-value storage system allowing channel owners to persist JSON/plain text data via the dashboard and scripts. Data is stored as JSONB in PostgreSQL with a 30MB per-channel size limit.

## Data Model

**Table: `channels_storage`**

```sql
CREATE TABLE channels_storage (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX channels_storage_channel_id_key_idx ON channels_storage(channel_id, key);
CREATE INDEX channels_storage_channel_id_idx ON channels_storage(channel_id);
```

- One row per key per channel
- `value` stores any JSON-valid data (string, number, array, object, null)
- Unique constraint on `(channel_id, key)` enables upsert via `INSERT ... ON CONFLICT DO UPDATE`
- Cascade delete when channel is removed
- Size check: `SELECT COALESCE(SUM(pg_column_size(value)), 0) FROM channels_storage WHERE channel_id = $1`

**Key constraints**: alphanumeric, underscores, hyphens, dots. Max 255 chars. Case-sensitive. Regex: `^[a-zA-Z0-9._-]+$`

## Backend

### Repository Layer

Location: `libs/repositories/channels_storage/`

```
channels_storage/
├── channels_storage.go      # Interface + inputs + errors
├── model/
│   └── model.go             # ChannelStorage struct
└── pgx/
    └── pgx.go               # pgx implementation
```

**Interface methods**:
- `GetAllByChannelID(ctx, channelID) ([]model.ChannelStorage, error)`
- `GetByKey(ctx, channelID, key) (model.ChannelStorage, error)`
- `Set(ctx, input SetInput) (model.ChannelStorage, error)` — upsert
- `Delete(ctx, channelID, key) error`
- `DeleteAllByChannelID(ctx, channelID) error`
- `GetTotalSizeByChannelID(ctx, channelID) (int64, error)` — `pg_column_size` sum

**SetInput**: `ChannelID`, `Key`, `Value` (json.RawMessage)

**Model**:
```go
type ChannelStorage struct {
    ID        uuid.UUID
    ChannelID string
    Key       string
    Value     json.RawMessage
    CreatedAt time.Time
    UpdatedAt time.Time
    isNil     bool
}
```

### Service Layer

Location: `apps/api-gql/internal/services/channels_storage/`

- Wraps repository
- `Set` checks total size before upserting: if current size + new value size > 30MB, return error
- `GetAllByChannelID` returns all entries for the channel
- No encryption needed

### GraphQL Schema

Location: `apps/api-gql/internal/delivery/gql/schema/storage.graphql`

```graphql
extend type Query {
    storageKeys: [StorageEntry!]!
        @isAuthenticated
        @hasAccessToSelectedDashboard
        @hasChannelRolesDashboardPermission(permission: VIEW_VARIABLES)
}

extend type Mutation {
    storageSet(key: String! @validate(constraint: "min=1|max=255|regex=^[a-zA-Z0-9._-]+$"), value: JSON!): StorageEntry!
        @isAuthenticated
        @hasAccessToSelectedDashboard
        @hasChannelRolesDashboardPermission(permission: MANAGE_VARIABLES)
    storageDelete(key: String!): Boolean!
        @isAuthenticated
        @hasAccessToSelectedDashboard
        @hasChannelRolesDashboardPermission(permission: MANAGE_VARIABLES)
    storageDeleteAll: Boolean!
        @isAuthenticated
        @hasAccessToSelectedDashboard
        @hasChannelRolesDashboardPermission(permission: MANAGE_VARIABLES)
}

type StorageEntry {
    key: String!
    value: JSON!
    createdAt: DateTime!
    updatedAt: DateTime!
}
```

- Uses `JSON` scalar (already exists in gqlgen)
- Reuses `VIEW_VARIABLES` / `MANAGE_VARIABLES` permissions

### Mapper

Location: `apps/api-gql/internal/delivery/gql/mappers/storage.go`

Maps `entity.ChannelStorage` → `gqlmodel.StorageEntry`.

### Resolver

Location: `apps/api-gql/internal/delivery/gql/resolvers/storage.resolver.go`

- `StorageKeys` query → `service.GetAllByChannelID(dashboardId)`
- `StorageSet` mutation → `service.Set(dashboardId, key, value)`
- `StorageDelete` mutation → `service.Delete(dashboardId, key)`
- `StorageDeleteAll` mutation → `service.DeleteAllByChannelID(dashboardId)`

### Entity

Location: `apps/api-gql/internal/entity/storage.go`

```go
type ChannelStorage struct {
    ID        uuid.UUID
    ChannelID string
    Key       string
    Value     json.RawMessage
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Executron Sandbox Integration

### Architecture

Uses capability module pattern (same as `twir:fetch`). Storage operations are async — sandbox sends request to host, host executes DB call, returns result.

### Files

- `src/libs/sandbox/storage.ts` — injected JS that wraps `twir:storage` capability calls
- `src/libs/host-storage.ts` — host-side handler: receives operation, executes against Postgres, returns result
- Update `src/libs/executor.ts` — register `twir:storage` capability

### Host-side storage handler

- Connects to Postgres (reuses connection from `db.ts`)
- Implements all operations as direct SQL queries
- Size check on `set` operations

### Available methods

All methods are `async`.

**Core CRUD**:
| Method | Description |
|---|---|
| `twir.storage.get(key)` | Returns parsed value or `null` |
| `twir.storage.set(key, value)` | Upserts value |
| `twir.storage.delete(key)` | Deletes key, returns boolean |
| `twir.storage.has(key)` | Returns boolean |
| `twir.storage.keys()` | Returns array of all keys |
| `twir.storage.clear()` | Deletes all keys for the channel |

**Array helpers** (value at key must be an array):
| Method | Description |
|---|---|
| `twir.storage.push(key, ...items)` | Appends items to array |
| `twir.storage.pop(key)` | Removes + returns last element |
| `twir.storage.find(key, predicate)` | Finds first matching element |
| `twir.storage.filter(key, predicate)` | Returns filtered array |
| `twir.storage.splice(key, start, deleteCount, ...items)` | Array splice |

**Object helpers** (value at key must be an object, `path` is dot-notation like `"address.city"`):
| Method | Description |
|---|---|
| `twir.storage.getProperty(key, path)` | Gets nested value |
| `twir.storage.setProperty(key, path, value)` | Sets nested value |
| `twir.storage.deleteProperty(key, path)` | Deletes nested property |
| `twir.storage.hasProperty(key, path)` | Checks if path exists |
| `twir.storage.merge(key, partial)` | Shallow merge partial object |

**Implementation notes**:
- Array/object helpers: host fetches the full value, sandbox applies predicate/manipulation, host saves back
- Predicates for `find`/`filter` run inside the sandbox (user code)
- `merge` does shallow `Object.assign` on the host side

## Frontend

### Location

Tab "Storage" in the Variables page (`/dashboard/variables?tab=storage`).

### Layout

File browser style:
- **Left panel**: list of all keys, searchable, shows type badge (string/number/array/object)
- **Right panel**: selected key's value editor
  - Monaco JSON editor for arrays/objects
  - Simple input for strings/numbers
  - "Save" button (calls `storageSet` mutation)
  - "Delete Key" button with confirmation dialog
  - Shows entry size
- **Bottom bar**: total storage usage out of 30MB
- **"+ Add Entry"**: dialog with key input + value editor

### Files

**New files**:
- `features/storage/storage.vue` — main component (left panel + right panel)
- `features/storage/composables/use-storage-api.ts` — GraphQL queries/mutations
- `features/storage/composables/use-storage-table.ts` — key list columns
- `features/storage/ui/storage-key-actions.vue` — row actions
- `features/storage/ui/storage-editor.vue` — value editor (Monaco or simple input)
- `features/storage/ui/storage-create-dialog.vue` — add entry dialog
- `api/storage.ts` — API composable (createGlobalState)

**Updated files**:
- `pages/dashboard/variables/index.vue` — add Storage tab alongside Variables and Secrets
- Navigation config: no change needed (storage is a sub-tab, not a separate page)
- GraphQL codegen: regenerate after schema changes

### API Composable

```typescript
// api/storage.ts
export const useStorageApi = createGlobalState(() => {
    const storageQuery = useQuery({
        variables: {},
        context: { additionalTypenames: ['StorageInvalidateKey'] },
        query: graphql(`query GetStorageKeys { storageKeys { key value createdAt updatedAt } }`),
    })

    const useMutationStorageSet = () =>
        useMutation(graphql(`mutation StorageSet($key: String!, $value: JSON!) { storageSet(key: $key, value: $value) { key value createdAt updatedAt } }`), ['StorageInvalidateKey'])

    const useMutationStorageDelete = () =>
        useMutation(graphql(`mutation StorageDelete($key: String!) { storageDelete(key: $key) }`), ['StorageInvalidateKey'])

    const useMutationStorageDeleteAll = () =>
        useMutation(graphql(`mutation StorageDeleteAll { storageDeleteAll }`), ['StorageInvalidateKey'])

    return { storageQuery, useMutationStorageSet, useMutationStorageDelete, useMutationStorageDeleteAll }
})
```

## Size Limit Enforcement

- 30MB = `30 * 1024 * 1024` = 31,457,280 bytes
- Checked before every `set` operation (both API and sandbox)
- Query: `SELECT COALESCE(SUM(pg_column_size(value)), 0) FROM channels_storage WHERE channel_id = $1`
- If current size + new value size > limit, return error with current usage and limit
- Frontend shows usage bar at the bottom of the storage panel
