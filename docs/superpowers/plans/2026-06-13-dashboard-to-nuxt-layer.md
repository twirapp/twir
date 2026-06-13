# Dashboard → Nuxt Layer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate `frontend/dashboard` (Vue 3 SPA) into `web/layers/dashboard/` as a Nuxt layer with CSR-only rendering, then delete the old dashboard app.

**Architecture:** Nuxt layer at `web/layers/dashboard/` with `ssr: false`. File-based routing replaces vue-router. Auth via Nuxt middleware. Uses web's urql client (already has WebSocket subscriptions). `@nuxtjs/i18n` module for i18n.

**Tech Stack:** Nuxt 4, Vue 3, urql v2, vee-validate v5, zod 4, shadcn-vue (Ui prefix), @nuxtjs/i18n, Tailwind CSS 4

**Branch:** All work must be done on branch `refactor/dashboard-move-to-nuxt`. Create it from `main` before starting Task 1.

---

## File Map

### New files to create

| Path | Responsibility |
|------|----------------|
| `web/layers/dashboard/nuxt.config.ts` | Layer config (ssr: false) |
| `web/layers/dashboard/layouts/dashboard.vue` | Main layout (sidebar + header) |
| `web/layers/dashboard/layouts/popup.vue` | Popup widget layout |
| `web/layers/dashboard/layouts/fullscreen.vue` | Overlay editor layout |
| `web/layers/dashboard/middleware/auth.ts` | Auth + permission guards |
| `web/layers/dashboard/plugins/monaco.client.ts` | Monaco editor (client-only) |
| `web/layers/dashboard/plugins/youtube.client.ts` | YouTube IFrame API (client-only) |
| `web/layers/dashboard/pages/**/*.vue` | ~40 page files (file-based routing) |

### Directories to move (source → target)

| Source | Target |
|--------|--------|
| `frontend/dashboard/src/features/` | `web/layers/dashboard/features/` |
| `frontend/dashboard/src/components/` (non-ui) | `web/layers/dashboard/components/` |
| `frontend/dashboard/src/composables/` | `web/layers/dashboard/composables/` |
| `frontend/dashboard/src/config/` | `web/layers/dashboard/config/` |
| `frontend/dashboard/src/helpers/` | `web/layers/dashboard/helpers/` |
| `frontend/dashboard/src/types/` | `web/layers/dashboard/types/` |
| `frontend/dashboard/src/locales/` | `web/layers/dashboard/locales/` |
| `frontend/dashboard/src/assets/` | `web/layers/dashboard/assets/` |
| `frontend/dashboard/src/api/` | `web/layers/dashboard/api/` |
| `frontend/dashboard/src/gql/` | `web/layers/dashboard/gql/` |
| `frontend/dashboard/src/lib/` | `web/layers/dashboard/lib/` |

### Files to modify

| Path | Change |
|------|--------|
| `web/nuxt.config.ts` | Add dashboard layer + @nuxtjs/i18n module |
| `web/package.json` | Add dashboard-specific deps |
| `web/urql.ts` | Add API key support from URL query |
| `web/codegen.ts` | Include dashboard gql output |

### Files to delete (Phase 5)

| Path | Reason |
|------|--------|
| `frontend/dashboard/` (entire directory) | Replaced by layer |

---

## Phase 1: Infrastructure

### Task 1: Create layer skeleton

**Files:**
- Create: `web/layers/dashboard/nuxt.config.ts`
- Create: `web/layers/dashboard/package.json`
- Modify: `web/nuxt.config.ts`

- [ ] **Step 1: Create dashboard layer directory**

```bash
mkdir -p web/layers/dashboard/{layouts,middleware,pages,pages/integrations,pages/integrations/callbacks,pages/commands,pages/events,pages/giveaways,pages/overlays,pages/moderation,pages/timers,pages/variables,pages/community,pages/settings,pages/registry/overlays,plugins}
```

- [ ] **Step 2: Create layer nuxt.config.ts**

```ts
// web/layers/dashboard/nuxt.config.ts
export default defineNuxtConfig({
  ssr: false,
  routeRules: {
    '/dashboard/**': { ssr: false },
  },
  vite: {
    resolve: {
      alias: {
        vue: 'vue/dist/vue.esm-bundler.js',
      },
    },
  },
})
```

- [ ] **Step 3: Register layer in web/nuxt.config.ts**

Add `'./layers/dashboard'` to the `extends` array in `web/nuxt.config.ts`.

- [ ] **Step 4: Commit**

```bash
git add web/layers/dashboard/nuxt.config.ts web/nuxt.config.ts
git commit -m "feat(dashboard-layer): create layer skeleton with ssr:false"
```

---

### Task 2: Add dashboard dependencies to web

**Files:**
- Modify: `web/package.json`

- [ ] **Step 1: Add dashboard-specific deps**

Run these from `web/`:

```bash
cd web
bun add @guolao/vue-monaco-editor @editorjs/editorjs @editorjs/header @editorjs/list @editorjs/paragraph @editorjs/quote @editorjs/simple-image @editorjs/underline @editorjs/delimiter @formkit/drag-and-drop vue-draggable-plus @vuepic/vue-datepicker grid-layout-plus @discord-message-components/vue @twirapp/kappagen @unovis/vue @unovis/ts nanoid tinycolor2 date-fns lodash.chunk nested-css-to-flat @zero-dependency/utils vaul-vue vue3-moveable vue-i18n @nuxtjs/i18n
```

Also add dev deps:
```bash
cd web
bun add -d @intlify/unplugin-vue-i18n @types/tinycolor2 @types/lodash.chunk
```

- [ ] **Step 2: Verify install succeeds**

```bash
cd web && bun install
```

- [ ] **Step 3: Commit**

```bash
git add web/package.json web/bun.lock
git commit -m "feat(dashboard-layer): add dashboard-specific dependencies"
```

---

### Task 3: Configure @nuxtjs/i18n module

**Files:**
- Modify: `web/nuxt.config.ts`
- Move: `frontend/dashboard/src/locales/` → `web/layers/dashboard/locales/`

- [ ] **Step 1: Copy locale files**

```bash
cp -r frontend/dashboard/src/locales web/layers/dashboard/locales
```

- [ ] **Step 2: Add i18n module to web/nuxt.config.ts**

Add `'@nuxtjs/i18n'` to the `modules` array and add i18n config:

```ts
i18n: {
  locales: [
    { code: 'en', file: 'en.json' },
    { code: 'ru', file: 'ru.json' },
    { code: 'de', file: 'de.json' },
    { code: 'es', file: 'es.json' },
    { code: 'ja', file: 'ja.json' },
    { code: 'pt', file: 'pt.json' },
    { code: 'sk', file: 'sk.json' },
    { code: 'uk', file: 'uk.json' },
  ],
  defaultLocale: 'en',
  lazy: true,
  langDir: './layers/dashboard/locales/',
},
```

- [ ] **Step 3: Commit**

```bash
git add web/layers/dashboard/locales/ web/nuxt.config.ts
git commit -m "feat(dashboard-layer): configure @nuxtjs/i18n with dashboard locales"
```

---

### Task 4: Create auth middleware

**Files:**
- Create: `web/layers/dashboard/middleware/auth.ts`
- Copy: `frontend/dashboard/src/api/auth.ts` → `web/layers/dashboard/api/auth.ts`

- [ ] **Step 1: Copy auth API file**

```bash
cp frontend/dashboard/src/api/auth.ts web/layers/dashboard/api/auth.ts
```

- [ ] **Step 2: Adapt auth.ts for Nuxt**

Key changes needed in `web/layers/dashboard/api/auth.ts`:
- Replace `import { useQuery, useMutation } from '@urql/vue'` with urql v2 imports
- Replace `import { graphql } from '@/gql/graphql.js'` with correct path
- Replace TanStack Query usage with urql directly
- Remove `import { urqlClient } from '@/plugins/urql.js'` — use `useUrqlClient()` from nuxt-urql

- [ ] **Step 3: Create auth middleware**

```ts
// web/layers/dashboard/middleware/auth.ts
export default defineNuxtRouteMiddleware(async (to) => {
  // Popup routes skip auth
  if (to.path.startsWith('/dashboard/popup')) return

  const { client } = useUrqlClient()
  const { data } = await client.query(profileQueryDocument, {}).toPromise()

  if (!data?.authenticatedUser) {
    return navigateTo('/', { replace: true })
  }

  if (to.meta.adminOnly && !data.authenticatedUser.isBotAdmin) {
    return navigateTo('/dashboard/forbidden', { replace: true })
  }

  if (to.meta.neededPermission) {
    const hasAccess = await checkUserPermission(data, to.meta.neededPermission)
    if (!hasAccess) {
      return navigateTo('/dashboard/forbidden', { replace: true })
    }
  }
})
```

- [ ] **Step 4: Commit**

```bash
git add web/layers/dashboard/api/auth.ts web/layers/dashboard/middleware/auth.ts
git commit -m "feat(dashboard-layer): add auth middleware and API"
```

---

### Task 5: Create layouts

**Files:**
- Create: `web/layers/dashboard/layouts/dashboard.vue`
- Create: `web/layers/dashboard/layouts/popup.vue`
- Create: `web/layers/dashboard/layouts/fullscreen.vue`

- [ ] **Step 1: Copy layout source files**

```bash
cp frontend/dashboard/src/layout/layout.vue web/layers/dashboard/layouts/dashboard.vue
cp frontend/dashboard/src/popup-layout/popup-layout.vue web/layers/dashboard/layouts/popup.vue
```

- [ ] **Step 2: Create fullscreen layout**

```vue
<!-- web/layers/dashboard/layouts/fullscreen.vue -->
<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'
</script>

<template>
  <div class="w-full h-full">
    <slot />
    <Toaster />
  </div>
</template>
```

- [ ] **Step 3: Adapt dashboard.vue layout for Nuxt**

Key changes:
- Replace `<RouterView />` with `<slot />`
- Replace `useRoute()` / `useRouter()` with `useRoute()` from Nuxt
- Add `definePageMeta({ layout: 'dashboard' })` pattern
- Replace `@/components/ui/sonner` import (use web's `@/components/ui/sonner`)
- Replace `@/components/ui/tooltip` import
- Replace sidebar/header imports to point to layer paths

- [ ] **Step 4: Adapt popup.vue layout for Nuxt**

Replace `<RouterView />` with `<slot />`.

- [ ] **Step 5: Commit**

```bash
git add web/layers/dashboard/layouts/
git commit -m "feat(dashboard-layer): create dashboard, popup, fullscreen layouts"
```

---

### Task 6: Migrate sidebar and header components

**Files:**
- Move: `frontend/dashboard/src/layout/sidebar/` → `web/layers/dashboard/layout/sidebar/`
- Move: `frontend/dashboard/src/layout/header/` → `web/layers/dashboard/layout/header/`
- Move: `frontend/dashboard/src/layout/page-layout.vue` → `web/layers/dashboard/layout/page-layout.vue`
- Move: `frontend/dashboard/src/layout/shadcn-layout.vue` → `web/layers/dashboard/layout/shadcn-layout.vue`
- Move: `frontend/dashboard/src/layout/stream-info-editor.vue` → `web/layers/dashboard/layout/stream-info-editor.vue`
- Move: `frontend/dashboard/src/layout/use-public-page-href.ts` → `web/layers/dashboard/layout/use-public-page-href.ts`
- Move: `frontend/dashboard/src/layout/use-sidebar-collapse.ts` → `web/layers/dashboard/layout/use-sidebar-collapse.ts`

- [ ] **Step 1: Copy layout subdirectories**

```bash
cp -r frontend/dashboard/src/layout/sidebar web/layers/dashboard/layout/sidebar
cp -r frontend/dashboard/src/layout/header web/layers/dashboard/layout/header
cp frontend/dashboard/src/layout/*.vue frontend/dashboard/src/layout/*.ts web/layers/dashboard/layout/
```

- [ ] **Step 2: Adapt imports in sidebar files**

For each file in `web/layers/dashboard/layout/sidebar/`:
- Replace `@/components/ui/sidebar` → `@/components/ui/sidebar` (same path, web has sidebar component)
- Replace `@/layout/sidebar/*` imports to relative imports
- Replace `lucide-vue-next` imports (keep as-is initially)

- [ ] **Step 3: Adapt imports in header files**

For each file in `web/layers/dashboard/layout/header/`:
- Replace `@/components/ui/*` imports (web has matching components)
- Replace `@/api/*` imports to relative paths

- [ ] **Step 4: Commit**

```bash
git add web/layers/dashboard/layout/
git commit -m "feat(dashboard-layer): migrate sidebar and header components"
```

---

### Task 7: Migrate config/navigation.ts

**Files:**
- Move: `frontend/dashboard/src/config/` → `web/layers/dashboard/config/`

- [ ] **Step 1: Copy config directory**

```bash
cp -r frontend/dashboard/src/config web/layers/dashboard/config
```

- [ ] **Step 2: Adapt icon imports**

Navigation uses `lucide-vue-next` icons. Keep these imports as-is (we're keeping lucide-vue-next as a dep).

- [ ] **Step 3: Commit**

```bash
git add web/layers/dashboard/config/
git commit -m "feat(dashboard-layer): migrate navigation config"
```

---

### Task 8: Migrate composables

**Files:**
- Move: `frontend/dashboard/src/composables/` → `web/layers/dashboard/composables/`

- [ ] **Step 1: Copy composables**

```bash
cp -r frontend/dashboard/src/composables web/layers/dashboard/composables
```

- [ ] **Step 2: Adapt use-mutation.ts**

Replace TanStack Query wrapper with direct urql usage:

```ts
// web/layers/dashboard/composables/use-mutation.ts
import { useMutation } from '@urql/vue'

export function useTypedMutation<T, V>(document: any) {
  const { executeMutation } = useMutation<T>(document)
  return { executeMutation }
}
```

- [ ] **Step 3: Adapt use-pagination.ts**

Remove TanStack Table dependency if any, keep as pure state composable.

- [ ] **Step 4: Commit**

```bash
git add web/layers/dashboard/composables/
git commit -m "feat(dashboard-layer): migrate composables"
```

---

### Task 9: Migrate shared UI components (non-shadcn)

**Files:**
- Move: `frontend/dashboard/src/components/` (excluding `ui/`) → `web/layers/dashboard/components/`

- [ ] **Step 1: List dashboard non-ui components**

```bash
ls frontend/dashboard/src/components/ | grep -v ui
```

- [ ] **Step 2: Copy non-ui components**

```bash
# Copy everything except ui/ from dashboard components
for dir in $(ls frontend/dashboard/src/components/ | grep -v ui); do
  cp -r "frontend/dashboard/src/components/$dir" "web/layers/dashboard/components/$dir"
done
# Copy standalone files
cp frontend/dashboard/src/components/*.vue web/layers/dashboard/components/ 2>/dev/null || true
```

- [ ] **Step 3: Adapt imports in each component**

Replace `@/components/ui/*` — web uses `Ui` prefix. This is the largest mechanical change. For each component:
- `Button` → `UiButton`
- `Card` → `UiCard`
- `Dialog` → `UiDialog`
- etc.

- [ ] **Step 4: Commit**

```bash
git add web/layers/dashboard/components/
git commit -m "feat(dashboard-layer): migrate shared components"
```

---

### Task 10: Add missing shadcn components to web

**Files:**
- Modify: `web/app/components/ui/` (add missing components)

Dashboard has 45 shadcn components, web has 23. Missing from web:
- `copy-input`, `date-picker`, `drawer`, `editorjs`, `hamburger-menu`, `InputWithIcon`, `kbd`, `multi-select`, `progress`, `scroll-area`, `settings-modal`, `stepper`, `tabs`, `tags-input`, `textarea`

Also dashboard has `action-confirm.vue` and some components that may overlap.

- [ ] **Step 1: Diff component lists**

Compare `frontend/dashboard/src/components/ui/` vs `web/app/components/ui/` to get exact missing list.

- [ ] **Step 2: Add missing shadcn components via CLI**

```bash
cd web
bun run shadcn-vue add tabs textarea scroll-area progress drawer date-picker
```

- [ ] **Step 3: Copy custom components that aren't in shadcn registry**

These need manual copy + Ui prefix adaptation:
- `copy-input/`
- `editorjs/`
- `kbd/`
- `multi-select/`
- `tags-input/`
- `stepper/`
- `settings-modal.vue`
- `action-confirm.vue`
- `hamburger-menu.vue`
- `InputWithIcon.vue`

- [ ] **Step 4: Commit**

```bash
git add web/app/components/ui/
git commit -m "feat(dashboard-layer): add missing shadcn components to web"
```

---

## Phase 2: GraphQL & State

### Task 11: Adapt urql client for API key support

**Files:**
- Modify: `web/urql.ts`

Dashboard supports `?apiKey=` URL parameter for API key auth. Web's urql client needs this.

- [ ] **Step 1: Add API key to fetchOptions**

In `web/urql.ts`, modify the client-side `fetchOptions`:

```ts
fetchOptions: {
  credentials: 'include',
  headers: {
    ...headers,
    ...(typeof window !== 'undefined' ? getApiKeyHeader() : {}),
  },
},
```

Add helper:
```ts
function getApiKeyHeader(): Record<string, string> {
  if (typeof window === 'undefined') return {}
  const params = new URLSearchParams(window.location.search)
  const apiKey = params.get('apiKey')
  return apiKey ? { 'Api-Key': apiKey } : {}
}
```

- [ ] **Step 2: Commit**

```bash
git add web/urql.ts
git commit -m "feat(dashboard-layer): add API key support to urql client"
```

---

### Task 12: Migrate GraphQL queries and mutations

**Files:**
- Move: `frontend/dashboard/src/api/` → `web/layers/dashboard/api/`
- Move: `frontend/dashboard/src/gql/` → `web/layers/dashboard/gql/`

- [ ] **Step 1: Copy api and gql directories**

```bash
cp -r frontend/dashboard/src/api web/layers/dashboard/api
cp -r frontend/dashboard/src/gql web/layers/dashboard/gql
```

- [ ] **Step 2: Update codegen config if needed**

Web's `codegen.ts` already scans `./layers/**/*.{ts,vue}`. The dashboard API files will be picked up automatically. No changes needed to codegen config.

- [ ] **Step 3: Run codegen**

```bash
cd web && bun run graphql-codegen
```

This regenerates `web/app/gql/` including all dashboard queries.

- [ ] **Step 4: Update import paths in api/ files**

For each file in `web/layers/dashboard/api/`:
- Replace `import { graphql } from '@/gql/graphql.js'` → `import { graphql } from '~/app/gql/graphql'`
- Replace `import { urqlClient } from '@/plugins/urql.js'` → use `useUrqlClient()` from nuxt-urql

- [ ] **Step 5: Commit**

```bash
git add web/layers/dashboard/api/ web/layers/dashboard/gql/ web/app/gql/
git commit -m "feat(dashboard-layer): migrate GraphQL queries and run codegen"
```

---

### Task 13: Remove TanStack Query usage

**Files:**
- Modify: `web/layers/dashboard/api/auth.ts`
- Modify: various feature files that use `useQuery` from TanStack

- [ ] **Step 1: Find all TanStack Query imports**

```bash
grep -r "tanstack/vue-query" web/layers/dashboard/
grep -r "useQuery\|useMutation\|QueryClient" web/layers/dashboard/ --include="*.ts" --include="*.vue"
```

- [ ] **Step 2: Replace each usage with urql**

Pattern: `useQuery(queryKey, queryFn)` → `useQuery({ query: document })` from `@urql/vue`

- [ ] **Step 3: Remove broadcastQueryClient**

No longer needed — single app.

- [ ] **Step 4: Commit**

```bash
git add web/layers/dashboard/
git commit -m "feat(dashboard-layer): remove TanStack Query, use urql directly"
```

---

## Phase 3: Plugins

### Task 14: Create Monaco editor plugin

**Files:**
- Create: `web/layers/dashboard/plugins/monaco.client.ts`

- [ ] **Step 1: Create client-only plugin**

```ts
// web/layers/dashboard/plugins/monaco.client.ts
import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor'

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(VueMonacoEditorPlugin, {
    paths: {
      vs: 'https://cdn.jsdelivr.net/npm/monaco-editor@0.52.0/min/vs',
    },
  })
})
```

- [ ] **Step 2: Commit**

```bash
git add web/layers/dashboard/plugins/monaco.client.ts
git commit -m "feat(dashboard-layer): add Monaco editor client plugin"
```

---

### Task 15: Create YouTube player plugin

**Files:**
- Move: `frontend/dashboard/src/composables/useGlobalYoutubePlayer.ts` → `web/layers/dashboard/composables/`
- Create: `web/layers/dashboard/plugins/youtube.client.ts`

- [ ] **Step 1: Copy YouTube composable**

```bash
cp frontend/dashboard/src/composables/useGlobalYoutubePlayer.ts web/layers/dashboard/composables/
```

- [ ] **Step 2: Create YouTube IFrame API loader plugin**

```ts
// web/layers/dashboard/plugins/youtube.client.ts
export default defineNuxtPlugin(() => {
  // Load YouTube IFrame API
  if (typeof window !== 'undefined' && !window.YT) {
    const tag = document.createElement('script')
    tag.src = 'https://www.youtube.com/iframe_api'
    const firstScriptTag = document.getElementsByTagName('script')[0]
    firstScriptTag.parentNode?.insertBefore(tag, firstScriptTag)
  }
})
```

- [ ] **Step 3: Commit**

```bash
git add web/layers/dashboard/plugins/youtube.client.ts web/layers/dashboard/composables/useGlobalYoutubePlayer.ts
git commit -m "feat(dashboard-layer): add YouTube player composable and plugin"
```

---

## Phase 4: Pages & Features

### Task 16: Create main dashboard page + layout wiring

**Files:**
- Create: `web/layers/dashboard/pages/index.vue`
- Modify: `web/layers/dashboard/layouts/dashboard.vue` (final wiring)

- [ ] **Step 1: Copy dashboard main page**

```bash
cp frontend/dashboard/src/pages/Dashboard.vue web/layers/dashboard/pages/index.vue
```

- [ ] **Step 2: Add layout meta to page**

```vue
<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})
</script>
```

- [ ] **Step 3: Adapt component imports**

Replace `@/features/dashboard/*` → `../features/dashboard/*` or `~/layers/dashboard/features/dashboard/*`

- [ ] **Step 4: Test page renders**

Start dev server: `cd web && bun dev`
Navigate to `/dashboard` — should show main dashboard page with sidebar.

- [ ] **Step 5: Commit**

```bash
git add web/layers/dashboard/pages/index.vue
git commit -m "feat(dashboard-layer): create main dashboard page"
```

---

### Task 17: Migrate feature pages (batch 1: core)

Features: bot-settings, commands, integrations, modules, variables

- [ ] **Step 1: Copy feature directories**

```bash
cp -r frontend/dashboard/src/features/bot-settings web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/commands web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/integrations web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/modules web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/variables web/layers/dashboard/features/
```

- [ ] **Step 2: Create page files**

```bash
cp frontend/dashboard/src/pages/Integrations.vue web/layers/dashboard/pages/integrations/index.vue
cp frontend/dashboard/src/pages/Keywords.vue web/layers/dashboard/pages/keywords.vue
```

Create `web/layers/dashboard/pages/bot-settings.vue`:
```vue
<script setup lang="ts">
import BotSettings from '~/layers/dashboard/features/bot-settings/bot-settings.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', neededPermission: 'ViewBotSettings', noPadding: true })
</script>

<template><BotSettings /></template>
```

Repeat pattern for each feature page.

- [ ] **Step 3: Adapt imports in each feature**

Replace `@/` imports with correct relative or `~/layers/dashboard/` paths.
Replace shadcn component names with Ui-prefixed versions.

- [ ] **Step 4: Commit**

```bash
git add web/layers/dashboard/features/ web/layers/dashboard/pages/
git commit -m "feat(dashboard-layer): migrate bot-settings, commands, integrations, modules, variables"
```

---

### Task 18: Migrate feature pages (batch 2: events & alerts)

Features: events, chat-alerts, alerts, moderation

- [ ] **Step 1: Copy feature directories**

```bash
cp -r frontend/dashboard/src/features/events web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/chat-alerts web/layers/dashboard/features/ 2>/dev/null || true
cp -r frontend/dashboard/src/features/alerts web/layers/dashboard/features/ 2>/dev/null || true
cp -r frontend/dashboard/src/features/moderation web/layers/dashboard/features/
```

- [ ] **Step 2: Create page files**

```bash
cp frontend/dashboard/src/pages/chat-alerts.vue web/layers/dashboard/pages/events/chat-alerts.vue
cp frontend/dashboard/src/pages/alerts.vue web/layers/dashboard/pages/alerts.vue
```

Create pages for events list and event form:
```
web/layers/dashboard/pages/events/index.vue
web/layers/dashboard/pages/events/[id].vue
web/layers/dashboard/pages/moderation/index.vue
web/layers/dashboard/pages/moderation/[id].vue
```

- [ ] **Step 3: Adapt imports**

- [ ] **Step 4: Commit**

```bash
git add web/layers/dashboard/features/ web/layers/dashboard/pages/
git commit -m "feat(dashboard-layer): migrate events, chat-alerts, alerts, moderation"
```

---

### Task 19: Migrate feature pages (batch 3: overlays)

Features: overlays (chat, kappagen, brb, tts, obs, dudes, faceit-stats, valorant-stats), overlay-builder

- [ ] **Step 1: Copy feature directories**

```bash
cp -r frontend/dashboard/src/features/overlays web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/overlay-builder web/layers/dashboard/features/
```

- [ ] **Step 2: Copy overlay pages**

```bash
cp frontend/dashboard/src/pages/Overlays.vue web/layers/dashboard/pages/overlays/index.vue
cp frontend/dashboard/src/pages/overlays/chat/Chat.vue web/layers/dashboard/pages/overlays/chat.vue
cp frontend/dashboard/src/pages/overlays/dudes/dudes-settings.vue web/layers/dashboard/pages/overlays/dudes.vue
```

Create remaining overlay pages:
```
web/layers/dashboard/pages/overlays/kappagen.vue
web/layers/dashboard/pages/overlays/brb.vue
web/layers/dashboard/pages/overlays/tts.vue
web/layers/dashboard/pages/overlays/obs.vue
web/layers/dashboard/pages/overlays/faceit-stats.vue
web/layers/dashboard/pages/overlays/valorant-stats.vue
```

- [ ] **Step 3: Create registry overlay edit page**

```bash
mkdir -p web/layers/dashboard/pages/registry/overlays
```

Create `web/layers/dashboard/pages/registry/overlays/[id].vue` from `frontend/dashboard/src/components/registry/overlays/edit.vue`.

- [ ] **Step 4: Adapt imports**

- [ ] **Step 5: Commit**

```bash
git add web/layers/dashboard/features/ web/layers/dashboard/pages/
git commit -m "feat(dashboard-layer): migrate overlays and overlay builder"
```

---

### Task 20: Migrate feature pages (batch 4: remaining)

Features: timers, giveaways, games, song-requests, greetings, community, expiring-vips, import, admin-panel, dashboard-widgets

- [ ] **Step 1: Copy remaining feature directories**

```bash
cp -r frontend/dashboard/src/features/timers web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/giveaways web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/games web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/greetings web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/expiring-vips web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/import web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/admin-panel web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/dashboard web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/keywords web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/community-chat-messages web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/community-emotes-statistic web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/community-rewards-history web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/community-roles web/layers/dashboard/features/
cp -r frontend/dashboard/src/features/community-users web/layers/dashboard/features/
```

- [ ] **Step 2: Create remaining page files**

```
web/layers/dashboard/pages/timers/index.vue
web/layers/dashboard/pages/timers/[id].vue
web/layers/dashboard/pages/giveaways/index.vue
web/layers/dashboard/pages/giveaways/view/[id].vue
web/layers/dashboard/pages/games.vue
web/layers/dashboard/pages/song-requests.vue
web/layers/dashboard/pages/greetings.vue
web/layers/dashboard/pages/community.vue
web/layers/dashboard/pages/community/roles.vue
web/layers/dashboard/pages/expiring-vips.vue
web/layers/dashboard/pages/import.vue
web/layers/dashboard/pages/admin.vue
web/layers/dashboard/pages/keywords.vue
web/layers/dashboard/pages/files.vue
web/layers/dashboard/pages/notifications.vue
web/layers/dashboard/pages/settings/index.vue
web/layers/dashboard/pages/settings/custom-widgets.vue
web/layers/dashboard/pages/forbidden.vue
web/layers/dashboard/pages/[...slug].vue  (404 catch-all)
```

- [ ] **Step 3: Create integration callback pages**

```
web/layers/dashboard/pages/integrations/callbacks/spotify.vue
web/layers/dashboard/pages/integrations/callbacks/donationalerts.vue
web/layers/dashboard/pages/integrations/callbacks/nightbot.vue
web/layers/dashboard/pages/integrations/callbacks/valorant.vue
web/layers/dashboard/pages/integrations/callbacks/discord.vue
web/layers/dashboard/pages/integrations/callbacks/vk.vue
web/layers/dashboard/pages/integrations/callbacks/[name].vue
web/layers/dashboard/pages/integrations/discord.vue
```

- [ ] **Step 4: Adapt imports in all features**

Bulk replace:
- `@/components/ui/*` → Ui-prefixed component names
- `@/features/*` → `~/layers/dashboard/features/*`
- `@/api/*` → `~/layers/dashboard/api/*`
- `@/composables/*` → `~/layers/dashboard/composables/*`
- `@/gql/*` → `~/layers/dashboard/gql/*` or `~/app/gql/*`

- [ ] **Step 5: Commit**

```bash
git add web/layers/dashboard/features/ web/layers/dashboard/pages/
git commit -m "feat(dashboard-layer): migrate all remaining features and pages"
```

---

### Task 21: Migrate remaining assets and helpers

**Files:**
- Move: `frontend/dashboard/src/assets/` → `web/layers/dashboard/assets/`
- Move: `frontend/dashboard/src/helpers/` → `web/layers/dashboard/helpers/`
- Move: `frontend/dashboard/src/lib/` → `web/layers/dashboard/lib/`
- Move: `frontend/dashboard/src/types/` → `web/layers/dashboard/types/`
- Move: `frontend/dashboard/src/plugins/i18n.ts` → adapted for @nuxtjs/i18n

- [ ] **Step 1: Copy directories**

```bash
cp -r frontend/dashboard/src/assets web/layers/dashboard/assets
cp -r frontend/dashboard/src/helpers web/layers/dashboard/helpers
cp -r frontend/dashboard/src/lib web/layers/dashboard/lib
cp -r frontend/dashboard/src/types web/layers/dashboard/types
```

- [ ] **Step 2: Remove i18n plugin (handled by @nuxtjs/i18n module)**

Do NOT copy `frontend/dashboard/src/plugins/i18n.ts` — the `@nuxtjs/i18n` module handles this.

- [ ] **Step 3: Commit**

```bash
git add web/layers/dashboard/assets/ web/layers/dashboard/helpers/ web/layers/dashboard/lib/ web/layers/dashboard/types/
git commit -m "feat(dashboard-layer): migrate assets, helpers, lib, types"
```

---

## Phase 5: Cleanup

### Task 22: Delete old dashboard and update CI

**Files:**
- Delete: `frontend/dashboard/` (entire directory)
- Modify: `.github/workflows/dockerv3.yml` (remove dashboard from matrix)
- Modify: `.github/workflows/build-and-lint.yml` (remove dashboard)
- Modify: Root `package.json` (remove dashboard workspace if listed)

- [ ] **Step 1: Verify web serves dashboard correctly**

```bash
cd web && bun dev
# Navigate to /dashboard, test auth flow, test a few pages
```

- [ ] **Step 2: Delete frontend/dashboard**

```bash
rm -rf frontend/dashboard
```

- [ ] **Step 3: Update CI workflows**

Remove `frontend/dashboard` from Docker build matrix and lint workflow.

- [ ] **Step 4: Update root package.json if needed**

Remove `frontend/dashboard` from workspaces if explicitly listed.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "chore: remove frontend/dashboard, fully migrated to web/layers/dashboard"
```

---

## Execution Order Summary

```
Phase 1 (Infrastructure):  Tasks 1-10
Phase 2 (GraphQL & State): Tasks 11-13
Phase 3 (Plugins):         Tasks 14-15
Phase 4 (Pages & Features): Tasks 16-21
Phase 5 (Cleanup):         Task 22
```

Each phase produces working, testable software. After Phase 1, the layer exists. After Phase 2, data flows. After Phase 3, client-only features work. After Phase 4, all pages are migrated. Phase 5 removes the old app.
