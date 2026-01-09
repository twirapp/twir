# Dashboard to Web Nuxt Monolith Migration Plan

## Overview

Migrate the standalone Vue 3 + Vite dashboard application (`frontend/dashboard/`) into the Nuxt 4 web monolith (`web/`) as a new layer. This consolidates two separate deployments into one unified application while maintaining all existing functionality and routes.

## Key Objectives

- Create `/web/layers/dashboard/` as a new Nuxt layer
- Migrate from Vue Router to Nuxt file-based routing
- Convert TanStack Vue Query to Pinia stores
- Migrate from vue-i18n to @nuxtjs/i18n module
- Audit and merge shadcn-vue components
- Maintain all `/dashboard/*` routes and features
- Single deployment artifact

## Implementation Phases

### Phase 1: Foundation & Setup

#### 1.1 Create Layer Directory Structure

Create the new dashboard layer:

```
web/layers/dashboard/
├── nuxt.config.ts
├── components/
│   ├── dashboard/
│   ├── integrations/
│   └── registry/
├── composables/
├── features/          # Copy entire features/ directory
├── layouts/
│   ├── default.vue
│   └── popup.vue
├── middleware/
│   └── auth.global.ts
├── pages/
│   └── dashboard/
├── stores/
├── utils/
│   └── api/
├── assets/
│   └── css/
├── locales/
└── types/
```

#### 1.2 Update Dependencies

**File: `/home/satont/Projects/twir/web/package.json`**

Add dashboard-specific dependencies:

```json
{
  "dependencies": {
    "@crashmax/json-format-highlight": "1.1.0",
    "@discord-message-components/vue": "0.2.1",
    "@editorjs/delimiter": "1.4.0",
    "@editorjs/editorjs": "2.29.1",
    "@editorjs/header": "2.8.1",
    "@editorjs/list": "1.9.0",
    "@editorjs/paragraph": "2.11.4",
    "@editorjs/quote": "2.6.0",
    "@editorjs/simple-image": "1.6.0",
    "@editorjs/underline": "1.1.0",
    "@formkit/drag-and-drop": "0.0.38",
    "@guolao/vue-monaco-editor": "1.5.1",
    "@protobuf-ts/twirp-transport": "2.9.4",
    "@tabler/icons-vue": "2.46.0",
    "@tanstack/vue-virtual": "3.13.12",
    "@twirapp/kappagen": "1.3.0",
    "@types/lodash.chunk": "4.2.9",
    "@vuepic/vue-datepicker": "12.1.0",
    "@vueuse/components": "catalog:",
    "@vueuse/router": "catalog:",
    "@zero-dependency/utils": "1.7.7",
    "date-fns": "3.6.0",
    "grid-layout-plus": "1.0.4",
    "lightweight-charts": "4.1.3",
    "lodash.chunk": "4.2.0",
    "nanoid": "5.1.6",
    "nested-css-to-flat": "1.0.5",
    "tailwindcss-animate": "1.0.7",
    "tinycolor2": "1.6.0",
    "vaul-vue": "0.4.1",
    "vue-draggable-plus": "0.5.6",
    "vue3-moveable": "0.28.0",
    "reka-ui": "2.7.0"
  },
  "devDependencies": {
    "@nuxtjs/i18n": "^8.0.0",
    "@protobuf-ts/runtime-rpc": "2.9.4",
    "@types/tinycolor2": "1.4.6",
    "sass": "1.75.0"
  }
}
```

**Version Upgrades Required:**
- `vee-validate`: 4.15.1 → 5.0.0-beta.0 (major - will need code updates)
- `@vee-validate/zod`: Update to match vee-validate 5.x
- Update all validation schemas to match new API

Run: `bun install`

#### 1.3 Configure Nuxt

**File: `/home/satont/Projects/twir/web/nuxt.config.ts`**

Add i18n module and dashboard layer:

```typescript
export default defineNuxtConfig({
  // Add to existing extends array (create if doesn't exist)
  extends: [
    './layers/dashboard'
  ],

  modules: [
    // ... existing modules
    '@nuxtjs/i18n', // Add this
  ],

  // Add i18n configuration
  i18n: {
    locales: [
      { code: 'en', name: 'English', file: 'en.json' },
      { code: 'ru', name: 'Russian', file: 'ru.json' },
      { code: 'uk', name: 'Українська', file: 'uk.json' },
      { code: 'de', name: 'Deutsch', file: 'de.json' },
      { code: 'ja', name: '日本語', file: 'ja.json' },
      { code: 'sk', name: 'Slovenčina', file: 'sk.json' },
      { code: 'es', name: 'Español', file: 'es.json' },
      { code: 'pt', name: 'Português', file: 'pt.json' },
    ],
    defaultLocale: 'en',
    strategy: 'no_prefix',
    lazy: true,
    langDir: 'locales',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'twir_locale',
      redirectOn: 'root',
    },
  },

  // ... rest of config
})
```

**File: `/home/satont/Projects/twir/web/layers/dashboard/nuxt.config.ts`**

Create layer config:

```typescript
export default defineNuxtConfig({
  imports: {
    dirs: [
      'composables',
      'utils/**',
      'stores',
    ],
  },

  components: [
    {
      path: '~/components',
      pathPrefix: false,
    },
  ],
})
```

### Phase 2: State Management - Pinia Stores

#### 2.1 Create Auth Store

**File: `/home/satont/Projects/twir/web/layers/dashboard/stores/auth.ts`**

Migrate from `/frontend/dashboard/src/api/auth.ts`:

```typescript
import { defineStore } from 'pinia'
import { graphql } from '~/gql/gql.js'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

export const useDashboardAuth = defineStore('dashboard-auth', () => {
  const router = useRouter()

  const { data, executeQuery, fetching } = useQuery({
    query: graphql(`
      query AuthenticatedUser {
        authenticatedUser {
          id
          isBotAdmin
          isBanned
          isEnabled
          isBotModerator
          hideOnLandingPage
          botId
          apiKey
          twitchProfile {
            description
            displayName
            login
            profileImageUrl
          }
          selectedDashboardId
          availableDashboards {
            id
            flags
            twitchProfile {
              login
              displayName
              profileImageUrl
            }
            apiKey
            plan {
              id
              name
              maxCommands
              # ... copy all plan fields
            }
          }
        }
      }
    `),
  })

  const user = computed(() => data.value?.authenticatedUser)

  const { executeMutation: executeLogout } = useMutation(
    graphql(`mutation userLogout { logout }`)
  )

  async function logout() {
    await executeLogout({})
    await router.push('/')
  }

  function checkPermission(flag: ChannelRolePermissionEnum): boolean {
    if (!user.value) return false
    if (user.value.isBotAdmin) return true
    if (user.value.id === user.value.selectedDashboardId) return true

    const dashboard = user.value.availableDashboards.find(
      d => d.id === user.value.selectedDashboardId
    )
    if (!dashboard) return false

    return dashboard.flags.includes(ChannelRolePermissionEnum.CanAccessDashboard) ||
           dashboard.flags.includes(flag)
  }

  return {
    user,
    isLoading: fetching,
    fetchUser: executeQuery,
    logout,
    checkPermission,
  }
})
```

#### 2.2 Create Dashboard Selection Store

**File: `/home/satont/Projects/twir/web/layers/dashboard/stores/dashboard.ts`**

```typescript
export const useDashboardSelection = defineStore('dashboard-selection', () => {
  const { executeMutation } = useMutation(
    graphql(`
      mutation SetDashboard($dashboardId: String!) {
        authenticatedUserSelectDashboard(dashboardId: $dashboardId)
      }
    `)
  )

  async function selectDashboard(dashboardId: string) {
    await executeMutation({ dashboardId })
    await navigateTo('/dashboard', { replace: true })
    window.location.reload()
  }

  return { selectDashboard }
})
```

### Phase 3: Layouts & Middleware

#### 3.1 Create Default Layout

**File: `/home/satont/Projects/twir/web/layers/dashboard/layouts/default.vue`**

Migrate from `/frontend/dashboard/src/layout/layout.vue`:

```vue
<script setup lang="ts">
import { Toaster } from '~/components/ui/sonner'
import { TooltipProvider } from '~/components/ui/tooltip'

const route = useRoute()
const isFullScreen = computed(() => route.meta?.fullScreen === true)
</script>

<template>
  <TooltipProvider :delay-duration="100">
    <template v-if="isFullScreen">
      <div class="w-full h-full">
        <slot />
      </div>
      <Toaster />
    </template>
    <template v-else>
      <DashboardSidebar>
        <DashboardHeader />
        <div
          :style="{
            padding: route.meta?.noPadding ? undefined : '24px',
            height: '100%',
          }"
          class="bg-[#0b0b0c]"
        >
          <slot />
        </div>
        <Toaster />
      </DashboardSidebar>
    </template>
  </TooltipProvider>
</template>
```

**File: `/home/satont/Projects/twir/web/layers/dashboard/layouts/popup.vue`**

```vue
<template>
  <div class="popup-layout">
    <slot />
  </div>
</template>
```

#### 3.2 Create Auth Middleware

**File: `/home/satont/Projects/twir/web/layers/dashboard/middleware/auth.global.ts`**

```typescript
export default defineNuxtRouteMiddleware(async (to) => {
  // Only run on dashboard routes
  if (!to.path.startsWith('/dashboard')) return

  // Skip for popup routes (if they don't need auth)
  if (to.path.startsWith('/dashboard/popup')) return

  const authStore = useDashboardAuth()

  // Fetch user if not loaded
  if (!authStore.user && !authStore.isLoading) {
    await authStore.fetchUser()
  }

  // Redirect if no user
  if (!authStore.user) {
    return navigateTo('/', { external: true })
  }

  // Check admin-only routes
  if (to.meta.adminOnly && !authStore.user.isBotAdmin) {
    return navigateTo('/dashboard/forbidden')
  }

  // Check permissions
  if (to.meta.neededPermission) {
    const hasAccess = authStore.checkPermission(to.meta.neededPermission)
    if (!hasAccess) {
      return navigateTo('/dashboard/forbidden')
    }
  }
})
```

### Phase 4: Routing Migration

#### 4.1 Route Mapping

Convert Vue Router routes to Nuxt file-based pages. Based on `/frontend/dashboard/src/plugins/router.ts`:

| Old Route | New File | Meta |
|-----------|----------|------|
| `/dashboard` | `pages/dashboard/index.vue` | `noPadding: true` |
| `/dashboard/bot-settings` | `pages/dashboard/bot-settings.vue` | `neededPermission: ViewBotSettings, noPadding: true` |
| `/dashboard/integrations` | `pages/dashboard/integrations/index.vue` | `neededPermission: ViewIntegrations, noPadding: true` |
| `/dashboard/integrations/spotify` | `pages/dashboard/integrations/spotify.vue` | OAuth callback |
| `/dashboard/integrations/donationalerts` | `pages/dashboard/integrations/donationalerts.vue` | OAuth callback |
| `/dashboard/integrations/discord` | `pages/dashboard/integrations/discord.vue` | OAuth callback |
| `/dashboard/integrations/vk` | `pages/dashboard/integrations/vk.vue` | OAuth callback |
| `/dashboard/integrations/valorant` | `pages/dashboard/integrations/valorant.vue` | OAuth callback |
| `/dashboard/integrations/:integrationName` | `pages/dashboard/integrations/[integrationName].vue` | Dynamic callback |
| `/dashboard/integrations/discord-settings` | `pages/dashboard/integrations/discord-settings.vue` | Settings page |
| `/dashboard/commands/:system` | `pages/dashboard/commands/[system].vue` | `neededPermission: ViewCommands, noPadding: true` |
| `/dashboard/commands/:system/:id` | `pages/dashboard/commands/[system]-[id].vue` | `neededPermission: ManageCommands, noPadding: true` |
| `/dashboard/timers` | `pages/dashboard/timers/index.vue` | With timers feature |
| `/dashboard/overlays/*` | `pages/dashboard/overlays/*.vue` | Multiple overlay pages |
| `/dashboard/events/:id` | `pages/dashboard/events/[id].vue` | Event editor |
| `/dashboard/popup/widgets/eventslist` | `pages/dashboard/popup/widgets/eventslist.vue` | `layout: 'popup'` |
| `/dashboard/popup/widgets/audit-log` | `pages/dashboard/popup/widgets/audit-log.vue` | `layout: 'popup'` |

#### 4.2 Page Template Pattern

Each page should use `definePageMeta`:

```vue
<script setup lang="ts">
import { ChannelRolePermissionEnum } from '~/gql/graphql'

definePageMeta({
  layout: 'default', // or 'popup'
  middleware: ['auth'], // auto from auth.global.ts
  noPadding: true, // optional
  fullScreen: false, // optional
  neededPermission: ChannelRolePermissionEnum.ViewCommands, // optional
  adminOnly: false, // optional
})
</script>

<template>
  <!-- Page content -->
</template>
```

### Phase 5: Component Migration

#### 5.1 UI Components Audit

**Process:**
1. List all dashboard UI components: `find frontend/dashboard/src/components/ui -name "*.vue"`
2. List all web UI components: `find web/app/components/ui -name "*.vue"`
3. Compare and identify duplicates vs unique components
4. For duplicates: use web version, update imports
5. For dashboard-only: copy to `/web/app/components/ui/`

**Common components that may overlap:**
- Button, Input, Label, Select, Dropdown, Table, etc.

#### 5.2 Feature Components

Copy entire features directory:

```bash
cp -r frontend/dashboard/src/features/ web/layers/dashboard/features/
```

Then update imports in all feature files:
- Change `@/` to `~/` (for shared components)
- Change `@/` to `#dashboard/` (for layer-specific imports)

#### 5.3 Dashboard-Specific Components

Copy to layer:

```bash
cp -r frontend/dashboard/src/components/dashboard/ web/layers/dashboard/components/dashboard/
cp -r frontend/dashboard/src/components/integrations/ web/layers/dashboard/components/integrations/
cp -r frontend/dashboard/src/components/registry/ web/layers/dashboard/components/registry/
```

Update imports in these components.

### Phase 6: i18n Migration

#### 6.1 Copy Translation Files

```bash
cp -r frontend/dashboard/src/locales/* web/layers/dashboard/locales/
```

Files to copy:
- en.json (1414 lines)
- ru.json, uk.json, de.json, ja.json, sk.json, es.json, pt.json

#### 6.2 Create Locale Composable

**File: `/home/satont/Projects/twir/web/layers/dashboard/composables/use-dashboard-locale.ts`**

```typescript
export function useDashboardLocale() {
  const { locale, setLocale } = useI18n()
  const savedLocale = useLocalStorage('twirLocale', 'en')

  // Sync with localStorage for backward compatibility
  watch(locale, (newLocale) => {
    savedLocale.value = newLocale
  })

  // Apply saved locale on mount
  onMounted(() => {
    if (savedLocale.value && savedLocale.value !== locale.value) {
      setLocale(savedLocale.value)
    }
  })

  return { locale, setLocale }
}
```

#### 6.3 Component Updates

Components using i18n remain mostly the same:
- `{{ $t('key') }}` still works
- `const { t } = useI18n()` is auto-imported by Nuxt

### Phase 7: API Integration

#### 7.1 Copy API Utilities

```bash
cp -r frontend/dashboard/src/api/ web/layers/dashboard/utils/api/
```

Update imports and adapt to Nuxt patterns.

#### 7.2 Twirp Client Setup

**File: `/home/satont/Projects/twir/web/layers/dashboard/utils/api/twirp.ts`**

```typescript
import { TwirpFetchTransport } from '@protobuf-ts/twirp-transport'

export function createTwirpTransport(baseUrl: string) {
  return new TwirpFetchTransport({
    baseUrl,
    fetchInit: {
      credentials: 'include',
    },
  })
}
```

#### 7.3 URQL Error Handling

**File: `/home/satont/Projects/twir/web/layers/dashboard/plugins/urql-errors.ts`**

Add toast error handling for GraphQL mutations similar to dashboard pattern.

### Phase 8: Assets & Styles

#### 8.1 Copy CSS

```bash
cp frontend/dashboard/src/main.css web/layers/dashboard/assets/css/dashboard.css
```

Import in layer's nuxt.config.ts or in layout.

#### 8.2 Copy Other Assets

```bash
cp -r frontend/dashboard/src/assets/* web/layers/dashboard/assets/
```

### Phase 9: Composables & Utilities

#### 9.1 Copy Composables

```bash
cp -r frontend/dashboard/src/composables/* web/layers/dashboard/composables/
```

Update imports and adapt patterns:
- `use-mutation.ts` - adapt for Nuxt/URQL
- `use-pagination.ts` - keep as-is
- Others - review and adapt

#### 9.2 Copy Types

```bash
cp -r frontend/dashboard/src/types/* web/layers/dashboard/types/
```

### Phase 10: Testing & Validation

#### 10.1 Development Testing

1. Start dev server: `cd web && bun run dev`
2. Navigate to `http://localhost:3010/dashboard`
3. Test authentication flow
4. Test each major feature:
   - Dashboard home
   - Bot settings
   - Commands (built-in & custom)
   - Timers
   - Integrations (OAuth flows)
   - Overlays
   - Events
   - Community features
   - Moderation
   - Admin panel

#### 10.2 Route Testing

Verify all dashboard routes are accessible and render correctly.

#### 10.3 i18n Testing

Switch between all 8 languages and verify translations load.

#### 10.4 Permission Testing

Test with different user roles to verify permission checks work.

### Phase 11: Build & Deployment

#### 11.1 Production Build

```bash
cd web
bun run build
```

Verify build completes without errors.

#### 11.2 Build Output

Check `.output/` directory contains unified build including dashboard.

#### 11.3 Test Production Build

```bash
bun run start
```

Navigate to dashboard routes and verify functionality.

### Phase 12: Cleanup

#### 12.1 Archive Old Dashboard

```bash
mv frontend/dashboard frontend/dashboard.old
```

Or delete after confirming migration success.

#### 12.2 Update CI/CD

Remove dashboard-specific build/deployment jobs from CI/CD pipeline.

#### 12.3 Update Documentation

Update project README and documentation to reflect unified web structure.

## Critical Files Reference

### Source Files (Dashboard)
- `/home/satont/Projects/twir/frontend/dashboard/src/plugins/router.ts` - Route definitions
- `/home/satont/Projects/twir/frontend/dashboard/src/api/auth.ts` - Auth patterns
- `/home/satont/Projects/twir/frontend/dashboard/src/layout/layout.vue` - Layout structure
- `/home/satont/Projects/twir/frontend/dashboard/package.json` - Dependencies

### Target Files (Web)
- `/home/satont/Projects/twir/web/nuxt.config.ts` - Main configuration
- `/home/satont/Projects/twir/web/urql.ts` - URQL client setup
- `/home/satont/Projects/twir/web/codegen.ts` - GraphQL codegen
- `/home/satont/Projects/twir/web/package.json` - Dependencies

### New Files to Create
- `/home/satont/Projects/twir/web/layers/dashboard/nuxt.config.ts`
- `/home/satont/Projects/twir/web/layers/dashboard/layouts/default.vue`
- `/home/satont/Projects/twir/web/layers/dashboard/layouts/popup.vue`
- `/home/satont/Projects/twir/web/layers/dashboard/middleware/auth.global.ts`
- `/home/satont/Projects/twir/web/layers/dashboard/stores/auth.ts`
- `/home/satont/Projects/twir/web/layers/dashboard/stores/dashboard.ts`
- All page files in `pages/dashboard/`

## Known Risks & Mitigations

### High Risk
1. **vee-validate 4.x → 5.x upgrade**
   - Mitigation: Update all form validation code to v5 API during migration
   - Test all forms thoroughly

2. **State management conversion**
   - Mitigation: Create Pinia stores that mirror TanStack Query behavior
   - Test data fetching and mutations

### Medium Risk
1. **Component deduplication**
   - Mitigation: Careful audit process, test UI after merging

2. **Route migration**
   - Mitigation: Use route mapping table, test all routes

### Low Risk
1. **i18n migration** - @nuxtjs/i18n is compatible with vue-i18n
2. **URQL integration** - Web already has URQL configured
3. **GraphQL codegen** - Already configured to include layers

## Success Criteria

- [ ] All 50+ dashboard routes accessible at `/dashboard/*`
- [ ] Authentication and permission checks working
- [ ] All features functional (commands, overlays, integrations, etc.)
- [ ] All 8 languages working with translations
- [ ] Single `bun run dev` command starts unified app
- [ ] Production build succeeds
- [ ] No visual regressions
- [ ] HMR working in development
- [ ] All OAuth flows working
- [ ] GraphQL queries and mutations working

## Estimated Timeline

- **Week 1**: Foundation & setup, dependency management
- **Week 2-3**: Core migration (layouts, middleware, stores)
- **Week 4-6**: Feature migration (pages, components)
- **Week 7**: Testing & validation
- **Week 8**: Deployment & cleanup

Total: **6-8 weeks** for complete migration
