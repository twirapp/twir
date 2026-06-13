# Dashboard → Nuxt Layer Migration Design

## Summary

Migrate `frontend/dashboard` (Vue 3 SPA) into `web/layers/dashboard/` as a Nuxt layer with CSR-only rendering. This unifies the frontend monolith and eliminates the separate dashboard build pipeline.

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| SSR vs CSR | CSR only | Dashboard is an authenticated zone — no SEO benefit from SSR |
| Dependencies | Upgrade to web versions | vee-validate v5, zod 4, urql v2 — single source of truth |
| Routing | File-based + custom auth guards | Nuxt-way routing with middleware for permissions |
| UI Components | Use web's shadcn (Ui prefix) | Consistent component library, migrate custom ones |
| GraphQL codegen | Unified pipeline | Web's Nuxt module handles codegen |
| Workspace packages | Keep as is | `@twir/frontend-chat` etc. stay as workspace packages |
| TanStack Query | Remove | Use urql cache directly, consistent with web |
| Migration strategy | Incremental | Infra → UI → Features one by one |
| Old dashboard | Delete after migration | Clean repo |

## Architecture

### Layer Structure

```
web/layers/dashboard/
├── nuxt.config.ts           # ssr: false, dashboard-specific config
├── components/
│   ├── ui/                  # Custom shadcn components (color picker, etc.)
│   ├── dashboard/           # Dashboard widget components
│   ├── songRequests/        # Song request components
│   └── registry/            # Registry components
├── composables/             # Dashboard composables (useGlobalYoutubePlayer, etc.)
├── features/                # Feature modules
│   ├── admin-panel/
│   ├── alerts/
│   ├── bot-settings/
│   ├── chat-alerts/
│   ├── commands/
│   ├── community-chat-messages/
│   ├── community-emotes-statistic/
│   ├── community-rewards-history/
│   ├── community-roles/
│   ├── community-users/
│   ├── dashboard/
│   ├── events/
│   ├── expiring-vips/
│   ├── games/
│   ├── giveaways/
│   ├── greetings/
│   ├── import/
│   ├── integrations/
│   ├── keywords/
│   ├── moderation/
│   ├── modules/
│   ├── overlay-builder/
│   ├── overlays/
│   ├── timers/
│   └── variables/
├── layouts/
│   ├── dashboard.vue        # Main layout (sidebar + header)
│   ├── popup.vue            # Popup widget layout
│   └── fullscreen.vue       # Overlay editor layout
├── middleware/
│   └── auth.ts              # Auth + permission guards
├── pages/
│   ├── index.vue            # /dashboard
│   ├── admin.vue            # /dashboard/admin
│   ├── alerts.vue           # /dashboard/alerts
│   ├── bot-settings.vue     # /dashboard/bot-settings
│   ├── community.vue        # /dashboard/community
│   ├── community/
│   │   └── roles.vue        # /dashboard/community/roles
│   ├── commands/
│   │   ├── [system].vue     # /dashboard/commands/:system
│   │   └── [system]/
│   │       └── [id].vue     # /dashboard/commands/:system/:id
│   ├── events/
│   │   ├── index.vue        # /dashboard/events
│   │   ├── chat-alerts.vue  # /dashboard/events/chat-alerts
│   │   └── [id].vue         # /dashboard/events/:id
│   ├── expiring-vips.vue    # /dashboard/expiring-vips
│   ├── files.vue            # /dashboard/files
│   ├── games.vue            # /dashboard/games
│   ├── giveaways/
│   │   ├── index.vue        # /dashboard/giveaways
│   │   └── view/
│   │       └── [id].vue     # /dashboard/giveaways/view/:id
│   ├── greetings.vue        # /dashboard/greetings
│   ├── import.vue           # /dashboard/import
│   ├── integrations/
│   │   ├── index.vue        # /dashboard/integrations
│   │   ├── discord.vue      # /dashboard/integrations/discord-settings
│   │   └── callbacks/
│   │       ├── spotify.vue
│   │       ├── donationalerts.vue
│   │       ├── nightbot.vue
│   │       ├── valorant.vue
│   │       ├── discord.vue
│   │       ├── vk.vue
│   │       └── [name].vue
│   ├── keywords.vue         # /dashboard/keywords
│   ├── moderation/
│   │   ├── index.vue        # /dashboard/moderation
│   │   └── [id].vue         # /dashboard/moderation/:id
│   ├── modules.vue          # /dashboard/modules
│   ├── notifications.vue    # /dashboard/notifications
│   ├── overlays/
│   │   ├── index.vue        # /dashboard/overlays
│   │   ├── chat.vue         # /dashboard/overlays/chat
│   │   ├── kappagen.vue     # /dashboard/overlays/kappagen
│   │   ├── brb.vue          # /dashboard/overlays/brb
│   │   ├── tts.vue          # /dashboard/overlays/tts
│   │   ├── obs.vue          # /dashboard/overlays/obs
│   │   ├── dudes.vue        # /dashboard/overlays/dudes
│   │   ├── faceit-stats.vue # /dashboard/overlays/faceit-stats
│   │   └── valorant-stats.vue
│   ├── registry/
│   │   └── overlays/
│   │       └── [id].vue     # /dashboard/registry/overlays/:id
│   ├── settings/
│   │   ├── index.vue        # /dashboard/settings
│   │   └── custom-widgets.vue
│   ├── song-requests.vue    # /dashboard/song-requests
│   ├── timers/
│   │   ├── index.vue        # /dashboard/timers
│   │   └── [id].vue         # /dashboard/timers/:id
│   ├── variables/
│   │   ├── index.vue        # /dashboard/variables
│   │   └── [id].vue         # /dashboard/variables/:id
│   ├── forbidden.vue        # /dashboard/forbidden
│   └── [...slug].vue        # Catch-all 404
├── plugins/
│   ├── monaco.client.ts     # Monaco editor (client-only)
│   ├── youtube.client.ts    # YouTube player (client-only)
│   └── i18n.ts              # vue-i18n setup
└── utils/
```

### Layer nuxt.config.ts

```typescript
export default defineNuxtConfig({
  ssr: false,

  // Dashboard-specific Vite config
  vite: {
    resolve: {
      alias: {
        vue: 'vue/dist/vue.esm-bundler.js',
      },
    },
  },

  // Route rules — all dashboard routes are CSR
  routeRules: {
    '/dashboard/**': { ssr: false },
  },
})
```

### Web nuxt.config.ts Changes

Add layer to extends:
```typescript
extends: [
  './layers/landing',
  './layers/url-shortener',
  './layers/pastebin',
  './layers/public',
  './layers/overlays',
  './layers/dashboard',  // NEW
]
```

Add dashboard-specific deps to optimizeDeps:
```typescript
vite: {
  optimizeDeps: {
    include: [
      // existing...
      '@guolao/vue-monaco-editor',
      'grid-layout-plus',
      'vue-draggable-plus',
    ],
  },
}
```

### Auth Middleware

```typescript
// layers/dashboard/middleware/auth.ts
export default defineNuxtRouteMiddleware(async (to) => {
  if (to.path.startsWith('/dashboard/popup')) return

  const { data: profile } = await useAsyncQuery(profileQueryDocument)

  if (!profile.value?.authenticatedUser) {
    return navigateTo('/', { replace: true })
  }

  if (to.meta.adminOnly && !profile.value.authenticatedUser.isBotAdmin) {
    return navigateTo('/dashboard/forbidden')
  }

  if (to.meta.neededPermission) {
    const hasAccess = await checkPermission(to.meta.neededPermission)
    if (!hasAccess) {
      return navigateTo('/dashboard/forbidden')
    }
  }
})
```

### Layout Migration

**dashboard.vue layout** — main layout with sidebar + header:
- Migrates from `src/layout/layout.vue`
- Uses web's shadcn components (Ui prefix)
- Sidebar from `src/layout/sidebar/`
- Header from `src/layout/header/`

**popup.vue layout** — for popup widgets:
- Migrates from `src/popup-layout/popup-layout.vue`

**fullscreen.vue layout** — for overlay editors:
- Full-screen component editors

### GraphQL Migration

1. Remove `src/plugins/urql.ts` (custom client)
2. Use web's urql client configured via `@bicou/nuxt-urql`
3. Add WebSocket subscription support to web's urql config
4. Migrate all `useQuery`/`useMutation` calls from urql v1 → v2 API
5. Remove `@tanstack/vue-query` — use urql cache directly
6. Move `src/gql/` → `layers/dashboard/gql/` (generated types)
7. Move `src/api/` → `layers/dashboard/api/` (query documents)

### Dependency Changes

**Remove from dashboard** (already in web):
- `@urql/vue` (web has v2 via nuxt-urql)
- `@tanstack/vue-query` (removed)
- `@tanstack/query-broadcast-client-experimental`
- `vee-validate` (web has v5-beta)
- `zod` (web has v4)
- `vue` (web has it)
- `vue-router` (web has it)
- `tailwindcss` (web has it)
- `lucide-vue-next` (web uses @nuxt/icon)
- `reka-ui` (web has it)
- `shadcn-vue` (web has shadcn-nuxt)
- `clsx`, `tailwind-merge`, `class-variance-authority` (web has them)
- `vue-sonner` (web has it)
- `graphql`, `graphql-ws` (web has them)

**Add to web/package.json** (dashboard-specific):
- `@guolao/vue-monaco-editor`
- `@editorjs/editorjs`, `@editorjs/header`, `@editorjs/list`, `@editorjs/paragraph`, `@editorjs/quote`, `@editorjs/simple-image`, `@editorjs/underline`, `@editorjs/delimiter`
- `@formkit/drag-and-drop`
- `vue-draggable-plus`
- `vue-i18n`
- `@vuepic/vue-datepicker`
- `grid-layout-plus`
- `@discord-message-components/vue`
- `@twirapp/kappagen`
- `@unovis/vue`, `@unovis/ts`
- `nanoid`
- `tinycolor2`
- `date-fns`
- `lodash.chunk`
- `nested-css-to-flat`
- `@zero-dependency/utils`
- `vaul-vue`
- `vue3-moveable`

**Keep as workspace packages**:
- `@twir/frontend-chat`
- `@twir/frontend-faceit-stats`
- `@twir/frontend-now-playing`
- `@twir/frontend-valorant-stats`

### Icon Migration

Dashboard uses direct imports: `import { User } from 'lucide-vue-next'`
Web uses Nuxt Icon: `<Icon name="lucide:user" />`

Options:
1. Add `lucide-vue-next` as dep and keep direct imports (simpler migration)
2. Convert all to `<Icon>` component (more work, consistent with web)

Recommendation: Keep `lucide-vue-next` direct imports initially, convert to `<Icon>` in a separate pass.

### i18n

Dashboard uses `vue-i18n` with locale files in `src/locales/`.
Web doesn't use i18n.

Solution: Add `vue-i18n` as a Nuxt plugin in the dashboard layer:
```typescript
// layers/dashboard/plugins/i18n.ts
import { createI18n } from 'vue-i18n'
import en from '../locales/en.json'

export default defineNuxtPlugin((nuxtApp) => {
  const i18n = createI18n({
    legacy: false,
    locale: 'en',
    messages: { en },
  })
  nuxtApp.vueApp.use(i18n)
})
```

### Codegen Integration

Web already has GraphQL codegen via Nuxt module (`modules/gql-codegen`).
Dashboard queries/mutations will be picked up by the same codegen pipeline.

Steps:
1. Move `src/api/**/*.ts` → `layers/dashboard/api/**/*.ts`
2. Move `src/gql/` → `layers/dashboard/gql/`
3. Update codegen config to include dashboard layer paths
4. Remove dashboard's standalone codegen config

## Incremental Migration Steps

### Phase 1: Infrastructure
1. Create `web/layers/dashboard/` skeleton
2. Create `nuxt.config.ts` with `ssr: false`
3. Create layouts (dashboard, popup, fullscreen)
4. Create auth middleware
5. Register layer in web/nuxt.config.ts
6. Add dashboard-specific deps to web/package.json

### Phase 2: GraphQL & State
1. Adapt urql client config for WebSocket subscriptions
2. Move API queries/mutations to layer
3. Remove TanStack Query usage
4. Run codegen to generate types

### Phase 3: UI Components
1. Audit dashboard shadcn components vs web
2. Migrate custom components (color picker, etc.)
3. Add missing shadcn components via CLI
4. Adapt all component imports to Ui prefix

### Phase 4: Features (one by one)
1. Dashboard (main page)
2. Bot settings
3. Commands
4. Integrations
5. Events & chat-alerts
6. Alerts
7. Overlays (all variants)
8. Timers
9. Giveaways
10. Keywords
11. Variables
12. Moderation
13. Community
14. Games
15. Song requests
16. Greetings
17. Expiring VIPs
18. Files
19. Import
20. Notifications
21. Settings
22. Admin panel

### Phase 5: Cleanup
1. Delete `frontend/dashboard/`
2. Remove dashboard from CI/CD matrix
3. Update Dockerfile (web serves everything)
4. Update documentation

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| urql v1→v2 breaking changes | Test each query migration individually |
| vee-validate v4→v5 migration | Form-by-form, test validation behavior |
| zod 3→4 migration | Mostly compatible, test edge cases |
| Nuxt CSR-only quirks | Use `<ClientOnly>` for heavy client components |
| Monaco editor SSR issues | Use `.client.ts` plugin suffix |
| YouTube player SSR issues | Use `.client.ts` plugin suffix |
| Icon migration scope | Keep lucide-vue-next imports initially |
| Large PR size | Incremental migration, merge per-phase |
