# Frontend Overlays → Nuxt Layer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate `frontend/overlays` (Vue 3 SPA overlay renderer) into the existing `web/layers/overlays/` Nuxt layer, supporting both `/o/:apiKey/*` (canonical) and `/overlays/:apiKey/*` (redirect for backward compatibility) paths.

**Architecture:** Expand the existing `web/layers/overlays/` layer. All overlay pages go under `pages/o/[apiKey]/`. Nuxt `routeRules` redirect `/overlays/:apiKey/**` → `/o/:apiKey/**`. Dashboard's `copyOverlayLink.ts` updated to use `/o/` path. Old `frontend/overlays/` marked deprecated.

**Tech Stack:** Nuxt 4, Vue 3, urql v2, graphql-ws, TMI.js, obs-websocket-js, @twirapp/dudes-vue, @twirapp/kappagen

**Branch:** All work must be done on branch `refactor/overlays-move-to-nuxt`. Create it from `main` before starting Task 1.

---

## Context

### What frontend/overlays does

A standalone SPA that renders browser-source overlays for OBS. Each overlay type has its own route with an `:apiKey` parameter. The app connects to the backend via GraphQL subscriptions and WebSockets for real-time data.

### Overlay types (10 routes)

| Route | Component | Purpose |
|-------|-----------|---------|
| `/overlays/:apiKey/overlays` | `overlays.vue` | Custom overlay renderer (HTML/IMAGE layers) |
| `/overlays/:apiKey/chat` | `chat.vue` | Chat overlay (TMI.js) |
| `/overlays/:apiKey/dudes` | `dudes.vue` | Dudes avatar overlay (canvas) |
| `/overlays/:apiKey/kappagen` | `kappagen.vue` | Emote animation overlay |
| `/overlays/:apiKey/brb` | `be-right-back.vue` | BRB timer overlay |
| `/overlays/:apiKey/now-playing` | `now-playing.vue` | Now-playing widget |
| `/overlays/:apiKey/faceit-stats` | `faceit-stats.vue` | FACEIT stats widget |
| `/overlays/:apiKey/tts` | `tts.vue` | Text-to-speech (invisible) |
| `/overlays/:apiKey/obs` | `obs.vue` | OBS WebSocket controller (invisible) |
| `/overlays/:apiKey/alerts` | `alerts.vue` | Alert audio player (invisible) |

### Existing web/layers/overlays/

Already contains:
- `pages/o/[apiKey]/valorant-stats.client.vue` — Valorant stats overlay
- Empty `nuxt.config.ts`

### Key difference from dashboard

- **No auth** — API key in URL params, not cookies
- **No UI framework** — no shadcn, no layout with sidebar
- **Real-time heavy** — GraphQL subscriptions, TMI.js, OBS WebSocket, 7TV EventAPI
- **Client-only** — all overlays are client-rendered (canvas, Web Audio, Shadow DOM)

---

## File Map

### New files to create

| Path | Responsibility |
|------|----------------|
| `web/layers/overlays/pages/o/[apiKey]/index.vue` | Custom overlay renderer |
| `web/layers/overlays/pages/o/[apiKey]/chat.vue` | Chat overlay |
| `web/layers/overlays/pages/o/[apiKey]/dudes.vue` | Dudes overlay |
| `web/layers/overlays/pages/o/[apiKey]/kappagen.vue` | Kappagen overlay |
| `web/layers/overlays/pages/o/[apiKey]/brb.vue` | BRB overlay |
| `web/layers/overlays/pages/o/[apiKey]/now-playing.vue` | Now-playing widget |
| `web/layers/overlays/pages/o/[apiKey]/faceit-stats.vue` | FACEIT stats widget |
| `web/layers/overlays/pages/o/[apiKey]/tts.vue` | TTS overlay |
| `web/layers/overlays/pages/o/[apiKey]/obs.vue` | OBS controller |
| `web/layers/overlays/pages/o/[apiKey]/alerts.vue` | Alert audio player |
| `web/layers/overlays/composables/**/*.ts` | All composables from frontend/overlays |
| `web/layers/overlays/components/**/*.vue` | Shared overlay components |
| `web/layers/overlays/helpers.ts` | Utility functions |
| `web/layers/overlays/types.ts` | Type definitions |
| `web/layers/overlays/api.ts` | API client setup |

### Files to modify

| Path | Change |
|------|--------|
| `web/layers/overlays/nuxt.config.ts` | Add routeRules for redirect |
| `web/nuxt.config.ts` | Add dashboard deps if needed |
| `web/layers/dashboard/components/overlays/copyOverlayLink.ts` | Change `/overlays/` → `/o/` |
| `web/package.json` | Add overlay-specific deps |

### Directories to move

| Source | Target |
|--------|--------|
| `frontend/overlays/src/composables/` | `web/layers/overlays/composables/` |
| `frontend/overlays/src/components/` | `web/layers/overlays/components/` |
| `frontend/overlays/src/plugins/urql.ts` | Removed (use web's urql) |
| `frontend/overlays/src/gql/` | `web/layers/overlays/gql/` |

---

## Phase 1: Infrastructure

### Task 1: Configure overlays layer

**Files:**
- Modify: `web/layers/overlays/nuxt.config.ts`

- [ ] **Step 1: Update layer nuxt.config.ts**

```ts
// web/layers/overlays/nuxt.config.ts
export default defineNuxtConfig({
  ssr: false,
  routeRules: {
    '/overlays/:apiKey/**': {
      redirect: (to) => {
        return `/o/${to.params.apiKey}/${(to.params as any).pathMatch || ''}`
      },
    },
  },
})
```

- [ ] **Step 2: Verify redirect works**

Start dev server, navigate to `/overlays/test123/chat` — should redirect to `/o/test123/chat`.

- [ ] **Step 3: Commit**

```bash
git add web/layers/overlays/nuxt.config.ts
git commit -m "feat(overlays-layer): configure ssr:false and /overlays → /o redirect"
```

---

### Task 2: Add overlay dependencies to web

**Files:**
- Modify: `web/package.json`

- [ ] **Step 1: Check which deps are missing**

Compare `frontend/overlays/package.json` dependencies vs `web/package.json`. Missing from web:
- `tmi.js`
- `obs-websocket-js`
- `@twirapp/dudes-vue`
- `@twirapp/kappagen` (may already be in web from dashboard migration)
- `emoji-regex`

Workspace packages already shared: `@twir/api`, `@twir/fontsource`, `@twir/frontend-chat`, `@twir/frontend-faceit-stats`, `@twir/frontend-now-playing`, `@twir/grpc`

- [ ] **Step 2: Add missing deps**

```bash
cd web
bun add tmi.js obs-websocket-js @twirapp/dudes-vue emoji-regex
```

- [ ] **Step 3: Commit**

```bash
git add web/package.json web/bun.lock
git commit -m "feat(overlays-layer): add overlay-specific dependencies"
```

---

## Phase 2: Move Shared Code

### Task 3: Migrate composables, helpers, types

**Files:**
- Move: `frontend/overlays/src/composables/` → `web/layers/overlays/composables/`
- Move: `frontend/overlays/src/helpers.ts` → `web/layers/overlays/helpers.ts`
- Move: `frontend/overlays/src/types.ts` → `web/layers/overlays/types.ts`
- Move: `frontend/overlays/src/api.ts` → `web/layers/overlays/api.ts`

- [ ] **Step 1: Copy composables**

```bash
cp -r frontend/overlays/src/composables web/layers/overlays/composables
cp frontend/overlays/src/helpers.ts web/layers/overlays/helpers.ts
cp frontend/overlays/src/types.ts web/layers/overlays/types.ts
cp frontend/overlays/src/api.ts web/layers/overlays/api.ts
```

- [ ] **Step 2: Remove urql plugin (use web's urql)**

Do NOT copy `frontend/overlays/src/plugins/urql.ts` — use web's `@bicou/nuxt-urql`.

- [ ] **Step 3: Adapt composables to use web's urql**

For each composable that imports from `@/plugins/urql`:
- Replace with `useUrqlClient()` from nuxt-urql
- Replace `import { graphql } from '@/gql/graphql'` with `import { graphql } from '~/app/gql/graphql'`

Key files to adapt:
- `composables/brb/use-brb-graphql.ts`
- `composables/chat/use-chat-overlay-socket.ts`
- `composables/dudes/use-dudes-socket.ts`
- `composables/kappagen/use-kappagen-socket.ts`
- `composables/now-playing/use-now-playing-socket.ts`
- `composables/obs/use-obs-graphql.ts`
- `composables/overlays/use-overlays.ts`
- `composables/overlays/use-custom-overlay.ts`
- `composables/tts/use-tts-graphql.ts`

- [ ] **Step 4: Commit**

```bash
git add web/layers/overlays/composables/ web/layers/overlays/helpers.ts web/layers/overlays/types.ts web/layers/overlays/api.ts
git commit -m "feat(overlays-layer): migrate composables, helpers, types, api"
```

---

### Task 4: Migrate overlay components

**Files:**
- Move: `frontend/overlays/src/components/` → `web/layers/overlays/components/`

- [ ] **Step 1: Copy components**

```bash
cp -r frontend/overlays/src/components web/layers/overlays/components
```

Components:
- `html-layer.vue` — Shadow DOM HTML renderer
- `image-layer.vue` — Image layer renderer
- `brb-timer.vue` — BRB countdown/countup timer
- `brb-text-with-emotes.vue` — Text with inline emotes

- [ ] **Step 2: Adapt imports**

No shadcn components used — minimal changes needed. Just fix relative import paths.

- [ ] **Step 3: Commit**

```bash
git add web/layers/overlays/components/
git commit -m "feat(overlays-layer): migrate overlay components"
```

---

## Phase 3: GraphQL

### Task 5: Migrate GraphQL types and queries

**Files:**
- Move: `frontend/overlays/src/gql/` → `web/layers/overlays/gql/` (temporary)
- Modify: `web/codegen.ts` (if needed)

- [ ] **Step 1: Copy gql directory**

```bash
cp -r frontend/overlays/src/gql web/layers/overlays/gql
```

- [ ] **Step 2: Verify codegen picks up overlay queries**

Web's `codegen.ts` already scans `./layers/**/*.{ts,vue}`. The overlay composable files with `graphql()` calls will be picked up. Run codegen:

```bash
cd web && bun run graphql-codegen
```

- [ ] **Step 3: Update import paths in composables**

For each composable using GraphQL:
- Replace `import { graphql } from '@/gql/graphql'` → `import { graphql } from '~/app/gql/graphql'`
- Replace `import { graphql } from '../gql/graphql'` → `import { graphql } from '~/app/gql/graphql'`

- [ ] **Step 4: Remove copied gql/ directory**

After codegen generates types in `web/app/gql/`, the local `gql/` directory is no longer needed:

```bash
rm -rf web/layers/overlays/gql/
```

- [ ] **Step 5: Commit**

```bash
git add web/layers/overlays/ web/app/gql/
git commit -m "feat(overlays-layer): migrate GraphQL queries, run codegen"
```

---

## Phase 4: Pages

### Task 6: Create overlay pages

**Files:**
- Create: `web/layers/overlays/pages/o/[apiKey]/index.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/chat.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/dudes.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/kappagen.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/brb.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/now-playing.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/faceit-stats.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/tts.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/obs.vue`
- Create: `web/layers/overlays/pages/o/[apiKey]/alerts.vue`

- [ ] **Step 1: Create page directory**

```bash
mkdir -p web/layers/overlays/pages/o/\[apiKey\]
```

- [ ] **Step 2: Create index.vue (custom overlays)**

```vue
<script setup lang="ts">
import OverlaysPage from '~/layers/overlays/pages-src/overlays.vue'

definePageMeta({ layout: false })
</script>

<template>
  <OverlaysPage />
</template>
```

Actually — each page should be the full component. Copy and adapt:

```bash
cp frontend/overlays/src/pages/overlays.vue web/layers/overlays/pages/o/\[apiKey\]/index.vue
```

Then adapt imports in the file.

- [ ] **Step 3: Create remaining pages**

```bash
cp frontend/overlays/src/pages/overlays/chat.vue web/layers/overlays/pages/o/\[apiKey\]/chat.vue
cp frontend/overlays/src/pages/overlays/dudes.vue web/layers/overlays/pages/o/\[apiKey\]/dudes.vue
cp frontend/overlays/src/pages/overlays/kappagen.vue web/layers/overlays/pages/o/\[apiKey\]/kappagen.vue
cp frontend/overlays/src/pages/overlays/be-right-back.vue web/layers/overlays/pages/o/\[apiKey\]/brb.vue
cp frontend/overlays/src/pages/overlays/now-playing.vue web/layers/overlays/pages/o/\[apiKey\]/now-playing.vue
cp frontend/overlays/src/pages/overlays/faceit-stats.vue web/layers/overlays/pages/o/\[apiKey\]/faceit-stats.vue
cp frontend/overlays/src/pages/tts.vue web/layers/overlays/pages/o/\[apiKey\]/tts.vue
cp frontend/overlays/src/pages/obs.vue web/layers/overlays/pages/o/\[apiKey\]/obs.vue
cp frontend/overlays/src/pages/alerts.vue web/layers/overlays/pages/o/\[apiKey\]/alerts.vue
```

- [ ] **Step 4: Add layout: false to each page**

Each page needs `definePageMeta({ layout: false })` since overlays render without any layout wrapper.

- [ ] **Step 5: Adapt imports in each page**

Replace:
- `@/composables/` → `~/layers/overlays/composables/`
- `@/components/` → `~/layers/overlays/components/`
- `@/helpers` → `~/layers/overlays/helpers`
- `@/types` → `~/layers/overlays/types`
- `@/gql/` → `~/app/gql/`
- `@/plugins/urql` → use `useUrqlClient()` from nuxt-urql

- [ ] **Step 6: Commit**

```bash
git add web/layers/overlays/pages/
git commit -m "feat(overlays-layer): create all overlay pages under /o/[apiKey]"
```

---

### Task 7: Remove old overlays layer page (valorant-stats)

**Files:**
- Delete: `web/layers/overlays/pages/o/[apiKey]/valorant-stats.client.vue` (moved into overlay pages)

Wait — the existing valorant-stats page is already at `pages/o/[apiKey]/valorant-stats.client.vue`. This should remain as-is since it's a separate workspace package component. No changes needed.

- [ ] **Step 1: Verify valorant-stats page still works**

The existing page at `pages/o/[apiKey]/valorant-stats.client.vue` uses `@twir/frontend-valorant-stats` workspace package. It should coexist with the new pages.

- [ ] **Step 2: No action needed**

The file stays. The `.client.vue` suffix ensures it's client-only rendered.

---

## Phase 5: Dashboard Link Update

### Task 8: Update overlay URL generation in dashboard

**Files:**
- Modify: `web/layers/dashboard/components/overlays/copyOverlayLink.ts`

- [ ] **Step 1: Change overlay URL path**

In `web/layers/dashboard/components/overlays/copyOverlayLink.ts`, line 26:

Change:
```ts
return `${window.location.origin}/overlays/${overlayApiKey.value}/${overlayPath}`
```

To:
```ts
return `${window.location.origin}/o/${overlayApiKey.value}/${overlayPath}`
```

- [ ] **Step 2: Search for other /overlays/ references in dashboard**

```bash
rg "/overlays/" -g "*.ts" -g "*.vue" web/layers/dashboard/ | grep -v "node_modules"
```

Check if any other files generate overlay URLs. Update them all.

- [ ] **Step 3: Commit**

```bash
git add web/layers/dashboard/components/overlays/copyOverlayLink.ts
git commit -m "fix(dashboard): update overlay URL generation to use /o/ path"
```

---

## Phase 6: Cleanup

### Task 9: Deprecate frontend/overlays

**Files:**
- Modify: `frontend/overlays/AGENTS.md` (add deprecation notice)
- Modify: `.github/workflows/dockerv3.yml` (remove from matrix or mark deprecated)
- Modify: `.github/workflows/build-and-lint.yml` (remove from matrix)

- [ ] **Step 1: Add deprecation notice to AGENTS.md**

Add at the top of `frontend/overlays/AGENTS.md`:

```markdown
> **DEPRECATED:** This app has been migrated to `web/layers/overlays/`. 
> Do not add new features here. This directory will be removed after 
> all production traffic is routed through the Nuxt layer.
```

- [ ] **Step 2: Remove from CI matrix**

Remove `frontend/overlays` from Docker build matrix and lint workflow.

- [ ] **Step 3: Verify web serves all overlays correctly**

```bash
cd web && bun dev
# Test each overlay type:
# /o/{apiKey}/chat
# /o/{apiKey}/dudes
# /o/{apiKey}/kappagen
# /o/{apiKey}/brb
# /o/{apiKey}/now-playing
# /o/{apiKey}/faceit-stats
# /o/{apiKey}/tts
# /o/{apiKey}/obs
# /o/{apiKey}/alerts
# /o/{apiKey}/overlays (custom)
# /o/{apiKey}/valorant-stats (existing)
```

- [ ] **Step 4: Test backward compatibility redirect**

```bash
# Test redirect from old path:
# /overlays/{apiKey}/chat → /o/{apiKey}/chat
```

- [ ] **Step 5: Commit**

```bash
git add frontend/overlays/AGENTS.md .github/workflows/
git commit -m "chore: deprecate frontend/overlays, remove from CI"
```

---

## Execution Order Summary

```
Phase 1 (Infrastructure):  Tasks 1-2
Phase 2 (Shared Code):     Tasks 3-4
Phase 3 (GraphQL):         Task 5
Phase 4 (Pages):           Tasks 6-7
Phase 5 (Dashboard Links): Task 8
Phase 6 (Cleanup):         Task 9
```

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| TMI.js SSR issues | All overlay pages use `definePageMeta({ layout: false })` + client-only rendering |
| OBS WebSocket in browser | Already browser-compatible, no changes needed |
| GraphQL subscription setup | Web's urql already has subscriptionExchange configured |
| Workspace package compatibility | Packages are shared, no changes needed |
| Redirect performance | Nuxt routeRules redirect is server-side (301), fast |
