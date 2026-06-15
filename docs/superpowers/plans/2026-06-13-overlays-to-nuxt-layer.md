# Frontend Overlays → Nuxt Layer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate `frontend/overlays` (Vue 3 SPA overlay renderer) into the existing `web/layers/overlays/` Nuxt layer, supporting both `/o/:apiKey/*` (canonical) and `/overlays/:apiKey/*` (redirect for backward compatibility) paths. Move frontend workspace packages to `web/lib/` as plain folders.

**Architecture:** Expand the existing `web/layers/overlays/` layer. All overlay pages go under `pages/o/[apiKey]/`. Nuxt `routeRules` disable SSR for overlay routes and redirect `/overlays/:apiKey/**` → `/o/:apiKey/**`. Dashboard's `copyOverlayLink.ts` updated to use `/o/` path. Old `frontend/overlays/` marked deprecated.

**Tech Stack:** Nuxt 4, Vue 3, urql v2, graphql-ws, TMI.js, obs-websocket-js, @twirapp/dudes-vue, @twirapp/kappagen

**Branch:** `refactor/overlays-to-nuxt`

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
| `web/layers/overlays/nuxt.config.ts` | Add routeRules: `/o/:apiKey/**` → `{ ssr: false }`, redirect `/overlays/:apiKey/**` → `/o/:apiKey/**` |
| `web/layers/dashboard/components/overlays/copyOverlayLink.ts` | Change `/overlays/` → `/o/` |
| `web/package.json` | Add overlay-specific deps |

### Directories to move

| Source | Target |
|--------|--------|
| `frontend/overlays/src/composables/` | `web/layers/overlays/composables/` |
| `frontend/overlays/src/components/` | `web/layers/overlays/components/` |
| `frontend/overlays/src/plugins/urql.ts` | Removed (use web's urql) |
| `libs/frontend-chat/` | `web/lib/frontend-chat/` (plain folder, no package.json) |
| `libs/frontend-faceit-stats/` | `web/lib/frontend-faceit-stats/` (plain folder, no package.json) |
| `libs/frontend-now-playing/` | `web/lib/frontend-now-playing/` (plain folder, no package.json) |
| `libs/fontsource/` | `web/lib/fontsource/` (plain folder, no package.json) |

---

## Phase 0: Move Workspace Packages to web/lib/

### Task 0: Move frontend workspace packages to web/lib/

**Goal:** Move `@twir/frontend-chat`, `@twir/frontend-faceit-stats`, `@twir/frontend-now-playing`, `@twir/fontsource` from `libs/` to `web/lib/` as plain folders (no `package.json`). Update all imports.

**Files:**
- Move: `libs/frontend-chat/src/*` → `web/lib/frontend-chat/`
- Move: `libs/frontend-faceit-stats/src/*` → `web/lib/frontend-faceit-stats/`
- Move: `libs/frontend-now-playing/src/*` → `web/lib/frontend-now-playing/`
- Move: `libs/fontsource/src/*` → `web/lib/fontsource/`
- Modify: all files importing from `@twir/frontend-chat`, `@twir/frontend-faceit-stats`, `@twir/frontend-now-playing`, `@twir/fontsource`
- Modify: `web/package.json` (remove workspace deps)
- Modify: root `package.json` (remove from workspaces if needed)

- [ ] **Step 1: Create web/lib/ directories and copy source files**

```bash
mkdir -p web/lib/frontend-chat web/lib/frontend-faceit-stats web/lib/frontend-now-playing web/lib/fontsource
cp -r libs/frontend-chat/src/* web/lib/frontend-chat/
cp -r libs/frontend-faceit-stats/src/* web/lib/frontend-faceit-stats/
cp -r libs/frontend-now-playing/src/* web/lib/frontend-now-playing/
cp -r libs/fontsource/src/* web/lib/fontsource/
```

- [ ] **Step 2: Update imports in web (dashboard layer)**

Replace in all files under `web/`:
- `@twir/frontend-chat` → `~~/lib/frontend-chat`
- `@twir/frontend-faceit-stats` → `~~/lib/frontend-faceit-stats`
- `@twir/frontend-now-playing` → `~~/lib/frontend-now-playing`
- `@twir/fontsource` → `~~/lib/fontsource`

Files to update:
- `web/layers/dashboard/pages/dashboard/overlays/chat/constants.ts`
- `web/layers/dashboard/pages/dashboard/overlays/chat/components/Form.vue`
- `web/layers/dashboard/features/overlays/faceit-stats/builder.vue`
- `web/layers/dashboard/features/overlays/faceit-stats/composables/use-faceit-stats.ts`
- `web/layers/dashboard/pages/dashboard/overlays/now-playing.vue`
- `web/layers/dashboard/lib/fontsource/index.ts`
- `web/layers/dashboard/lib/fontsource/components/FontSelector.vue`

- [ ] **Step 3: Remove workspace deps from web/package.json**

Remove from `web/package.json` dependencies:
- `@twir/frontend-chat`
- `@twir/frontend-faceit-stats`
- `@twir/frontend-now-playing`
- `@twir/fontsource`

- [ ] **Step 4: Remove old libs/ directories**

```bash
rm -rf libs/frontend-chat libs/frontend-faceit-stats libs/frontend-now-playing libs/fontsource
```

- [ ] **Step 5: Verify dashboard still works**

```bash
cd web && bun dev
# Navigate to dashboard overlays pages — chat, faceit-stats, now-playing, fontsource should work
```

- [ ] **Step 6: Commit**

```bash
git add web/lib/ web/package.json libs/
git commit -m "refactor: move frontend workspace packages to web/lib/ as plain folders"
```

---

## Phase 1: Infrastructure

### Task 1: Configure overlays layer

**Files:**
- Modify: `web/layers/overlays/nuxt.config.ts`

- [ ] **Step 1: Update layer nuxt.config.ts**

```ts
// web/layers/overlays/nuxt.config.ts
export default defineNuxtConfig({
  routeRules: {
    '/o/:apiKey/**': { ssr: false },
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
git commit -m "feat(overlays-layer): configure routeRules for ssr:false and /overlays → /o redirect"
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
- `@twirapp/kappagen`
- `emoji-regex`

- [ ] **Step 2: Add missing deps**

```bash
cd web
bun add tmi.js obs-websocket-js @twirapp/dudes-vue @twirapp/kappagen emoji-regex
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

- [ ] **Step 3: Adapt composables to use web's urql and local libs**

For each composable:
- Replace `import { graphql } from '@/gql/graphql'` → `import { graphql } from '~~/app/gql/graphql'`
- Replace `@twir/frontend-chat` → `~~/lib/frontend-chat`
- Replace `@twir/frontend-faceit-stats` → `~~/lib/frontend-faceit-stats`
- Replace `@twir/frontend-now-playing` → `~~/lib/frontend-now-playing`
- Replace `@twir/fontsource` → `~~/lib/fontsource`
- Replace `@/plugins/urql` refs with `useUrqlClient()` from nuxt-urql

Key files to adapt:
- `composables/brb/use-brb-graphql.ts`
- `composables/chat/use-chat-overlay-socket.ts`
- `composables/dudes/use-dudes-socket.ts`
- `composables/dudes/use-dudes.ts`
- `composables/dudes/use-dudes-settings.ts`
- `composables/kappagen/use-kappagen-socket.ts`
- `composables/kappagen/use-kappagen-builder.ts`
- `composables/now-playing/use-now-playing-socket.ts`
- `composables/obs/use-obs-graphql.ts`
- `composables/overlays/use-overlays.ts`
- `composables/overlays/use-custom-overlay.ts`
- `composables/tmi/use-chat-tmi.ts`
- `composables/tmi/use-emotes.ts`
- `composables/tmi/use-message-helpers.ts`
- `composables/tts/use-tts-graphql.ts`
- `components/brb-timer.vue`

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

No shadcn components used — minimal changes needed. Fix relative import paths and update `@twir/fontsource` → `~~/lib/fontsource`.

- [ ] **Step 3: Commit**

```bash
git add web/layers/overlays/components/
git commit -m "feat(overlays-layer): migrate overlay components"
```

---

## Phase 3: GraphQL

### Task 5: Migrate GraphQL queries

**Goal:** No separate gql/ directory needed. Web's codegen scans `./layers/**/*.{ts,vue}` and auto-discovers overlay queries.

**Files:**
- No files to copy — composables already have `graphql()` calls that web's codegen will pick up
- Composables need import path updates (done in Task 3)

- [ ] **Step 1: Run codegen to generate types**

```bash
cd web && bun run graphql-codegen
```

Verify that `web/app/gql/` now contains types for overlay queries.

- [ ] **Step 2: Verify composables import from correct path**

All composables should import from `~~/app/gql/graphql`, not from any local `gql/` directory.

- [ ] **Step 3: Commit**

```bash
git add web/app/gql/
git commit -m "feat(overlays-layer): run codegen for overlay GraphQL queries"
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
- `@/gql/` → `~~/app/gql/`
- `@/plugins/urql` → use `useUrqlClient()` from nuxt-urql
- `@twir/frontend-chat` → `~~/lib/frontend-chat`
- `@twir/frontend-faceit-stats` → `~~/lib/frontend-faceit-stats`
- `@twir/frontend-now-playing` → `~~/lib/frontend-now-playing`
- `@twir/fontsource` → `~~/lib/fontsource`

- [ ] **Step 6: Commit**

```bash
git add web/layers/overlays/pages/
git commit -m "feat(overlays-layer): create all overlay pages under /o/[apiKey]"
```

---

### Task 7: Verify valorant-stats coexistence

**Files:**
- `web/layers/overlays/pages/o/[apiKey]/valorant-stats.client.vue` — keep as-is

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
Phase 0 (Workspace Pkg):   Task 0
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
| TMI.js SSR issues | `routeRules` with `{ ssr: false }` — same pattern as dashboard |
| OBS WebSocket in browser | Already browser-compatible, no changes needed |
| GraphQL subscription setup | Web's urql already has subscriptionExchange configured |
| Local lib imports | Use `~~/lib/` prefix (Nuxt alias) — consistent with existing patterns |
| Redirect performance | Nuxt routeRules redirect is server-side (301), fast |
| API key in URL (security) | Same as current — API key is designed for browser-source URLs |
