# AGENTS.md — web

Nuxt 3 public website with layered architecture. Landing pages, docs, and public services.

## OVERVIEW

Multi-layer Nuxt 3 application. Serves public-facing content including landing page, documentation, URL shortener, and pastebin. Uses Nuxt's layer system for feature separation.

## STRUCTURE

```
web/
├── app/                     # Main app code
│   ├── components/          # Vue components
│   ├── layouts/             # Nuxt layouts
│   ├── pages/               # Route pages
│   └── assets/              # Static assets
├── layers/                  # Nuxt layers (feature modules)
│   ├── landing/             # Marketing site
│   ├── public/              # Public utilities
│   ├── url-shortener/       # Short link service
│   ├── pastebin/            # Code sharing
│   └── overlays/            # Public overlay previews
├── nuxt.config.ts           # Main Nuxt config
├── package.json
└── Dockerfile
```

## ENTRY POINTS

| Type          | Path                      | Purpose               |
| ------------- | ------------------------- | --------------------- |
| Main Config   | `nuxt.config.ts`          | Nuxt configuration    |
| App Entry     | `app/app.vue`             | Root component        |
| Layer Configs | `layers/*/nuxt.config.ts` | Layer-specific config |

## KEY COMMANDS

```bash
# Development (runs on :3000)
bun dev                    # nuxt dev --no-fork

# Build for production
bun run build             # nuxt build

# Start production server
bun run start             # node .output/server/index.mjs
```

## CONVENTIONS

### Nuxt Layer Pattern

Each layer is a self-contained Nuxt app:

```
layers/url-shortener/
├── components/             # Layer components
├── pages/                  # Layer routes
├── nuxt.config.ts          # Layer config
└── package.json            # Layer deps (if any)
```

Layers are registered in main `nuxt.config.ts`:

```typescript
export default defineNuxtConfig({
	extends: ["./layers/landing", "./layers/url-shortener", "./layers/pastebin", "./layers/public"],
});
```

### Icons

Use Nuxt Icon component with Lucide:

```vue
<template>
	<Icon name="lucide:user" class="h-4 w-4" />
</template>
```

### Styling

Same as dashboard: Tailwind CSS with theme colors.

### GraphQL (urql)

```typescript
// SSR-aware urql client
const { data } = await useAsyncQuery(gql`
  query GetData {
    ...
  }
`);
```

## ANTI-PATTERNS

- **DO NOT** bypass layer boundaries — keep features in their layer
- **DO NOT** use client-side only features without `<ClientOnly>`
- **DO NOT** ignore SSR implications in async data fetching

## LAYERS

| Layer           | Purpose          | Routes        |
| --------------- | ---------------- | ------------- |
| `landing`       | Marketing site   | `/`           |
| `url-shortener` | Link shortening  | `/s/*`        |
| `pastebin`      | Code sharing     | `/paste/*`    |
| `public`        | Public utilities | `/public/*`   |
| `overlays`      | Overlay previews | `/overlays/*` |

## PORTS

| Service    | Port                         |
| ---------- | ---------------------------- |
| Dev Server | 3000                         |
| API Proxy  | 3000/api/\* → localhost:3009 |

## NOTES

- Uses `nitro.preset: 'bun'` for Bun runtime
- GraphQL endpoint configured for SSR: `/api/query`
- Telemetry enabled for Nuxt insights
- Image optimization via `@nuxt/image`
