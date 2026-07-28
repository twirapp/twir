# AGENTS.md — web

Nuxt 3 application with layered architecture: public website, dashboard SPA, docs, and public services.

## OVERVIEW

Multi-layer Nuxt 3 application. Serves public-facing content (landing page, documentation, URL shortener, pastebin) **and** the authenticated dashboard (`layers/dashboard`, migrated from the deleted `frontend/dashboard`).

## STRUCTURE

```
web/
├── app/                     # Main app code
│   ├── components/          # Vue components
│   ├── layouts/             # Nuxt layouts
│   ├── pages/               # Route pages
│   └── assets/              # Static assets
├── layers/                  # Nuxt layers (feature modules)
│   ├── dashboard/           # Dashboard SPA (own AGENTS.md)
│   ├── landing/             # Marketing site
│   ├── public/              # Public utilities
│   ├── url-shortener/       # Short link service
│   ├── pastebin/            # Code sharing
│   ├── overlays/            # Public overlay previews
│   └── widgets/             # Public widgets
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

Layers under `layers/` are auto-registered by Nuxt; explicit `extends` in main `nuxt.config.ts`
is only needed for layers outside that directory.

### Icons

Use the Nuxt `<Icon />` component everywhere, including shared UI components and layers:

```vue
<template>
	<Icon name="lucide:user" class="h-4 w-4" />
	<Icon name="simple-icons:twitch" class="h-4 w-4 text-[#9146FF]" />
</template>
```

- `lucide:*` — UI chrome (default).
- `simple-icons:*` — brand/platform logos (Twitch, Kick, VK, ...); locally installed collection.
- `twir-integrations:*` / `twir-overlays:*` / `twir-compare:*` — local SVG custom collections
  (dirs under `layers/dashboard/assets/*` and `layers/landing/assets/compare`).

Do not import icon components from `@lucide/vue` or `lucide-vue-next` in `web`.

### Styling

Same as dashboard: Tailwind CSS with theme colors.

### Vue Auto-Imports

Nuxt auto-imports Vue reactivity primitives. Do NOT import them manually:

```vue
<!-- WRONG -->
<script setup>
import { ref, computed, watch, onMounted } from 'vue'
</script>

<!-- CORRECT — just use them directly -->
<script setup>
const count = ref(0)
const doubled = computed(() => count.value * 2)
watch(count, (val) => { /* ... */ })
onMounted(() => { /* ... */ })
</script>
```

Auto-imported: `ref`, `computed`, `watch`, `watchEffect`, `onMounted`, `onUnmounted`, `nextTick`, `toRaw`, `unref`, `shallowRef`, etc.

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
| `dashboard`     | Dashboard SPA    | `/dashboard/*`|
| `landing`       | Marketing site   | `/`           |
| `url-shortener` | Link shortening  | `/s/*`        |
| `pastebin`      | Code sharing     | `/paste/*`    |
| `public`        | Public utilities | `/public/*`   |
| `overlays`      | Overlay previews | `/overlays/*` |
| `widgets`       | Public widgets   | `/w/*`        |

Layers in `layers/` are auto-registered by Nuxt (no manual `extends` needed).

## PORTS

| Service    | Port                         |
| ---------- | ---------------------------- |
| Dev Server | 3000                         |
| API Proxy  | 3000/api/\* → localhost:3009 |

## Playwright MCP

- For browser QA, open the dev site through Caddy using the current root `.env` `SITE_BASE_URL`.
- If `SITE_BASE_URL` is not set, use `http://localhost:3005`.
- Do not use `localhost:3000` or `localhost:3010` directly for authenticated dashboard QA; cookies,
  redirects, API proxying, and sockets must match the Caddy-served dev URL.

## NOTES

- Uses `nitro.preset: 'bun'` for Bun runtime
- GraphQL endpoint configured for SSR: `/api/query`
- Telemetry enabled for Nuxt insights
- Image optimization via `@nuxt/image`
