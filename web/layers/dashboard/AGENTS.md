# AGENTS.md — web/layers/dashboard

Main dashboard SPA (Nuxt layer, migrated from the deleted `frontend/dashboard`). Authenticated area under `/dashboard/*`.

## STRUCTURE

```
layers/dashboard/
├── pages/                 # Routes; thin wrappers only (see Page pattern)
├── features/              # Feature-per-dir: <name>/<name>.vue + ui/ + composables/ + api.ts + *.spec.ts
├── layout/                # App chrome: sidebar/, header/, page-layout.vue (PageLayout)
├── layouts/               # Nuxt layouts: dashboard.vue, fullscreen.vue, popup.vue
├── components/            # Shared dashboard components (dialog-or-sheet.vue, command-menu/, ...)
├── api/                   # Shared urql queries/mutations (auth.ts, dashboard.ts, openapi.ts, ...)
├── composables/           # use-theme, etc.
├── config/                # navigation.ts (sidebar menu), i18n-locales.ts
├── locales/               # en/ru/uk/de/es/pt/sk/ja JSON; sidebar.* keys drive menu labels
└── assets/                # SVG icon collections: overlays/ (twir-overlays), integrations/ (twir-integrations)
```

## WHERE TO LOOK

| Task                       | Location                                             |
| -------------------------- | ---------------------------------------------------- |
| Add dashboard page         | `pages/dashboard/<name>.vue` + `features/<name>/`    |
| Add sidebar menu item      | `config/navigation.ts` `baseNavigationItems`         |
| Shared GraphQL queries     | `api/*.ts`                                           |
| Feature GraphQL queries    | `features/<name>/api.ts`                             |
| Menu/page i18n labels      | `locales/en.json` (`sidebar.*`) + other locales      |
| Responsive dialog/sheet    | `components/dialog-or-sheet.vue`                     |

## CONVENTIONS

- **Page pattern**: `pages/dashboard/*.vue` only does `definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })` and renders one feature component. All logic lives in `features/<name>/`.
- **Feature layout**: pages render `layout/page-layout.vue` (`PageLayout`); tabs via its `PageLayoutTab[]`.
- **Navigation items**: use `translationKey: 'sidebar.*'` when a locale key exists, otherwise hardcoded `name` (both patterns present).
- **GraphQL**: typed documents from `~/gql/graphql.js` (codegen). After changing queries run `bun run graphql-codegen` in `web/`.
- **Icons**: `lucide:*` for UI chrome, `simple-icons:*` for platform/brand logos (locally installed collection), `twir-integrations:*` for integration cards.
- **Destructive confirmations**: use shared `web/app/components/ui/action-confirm/ActionConfirm.vue`; do not create new confirm dialogs.

## TESTING

```bash
# from web/
bun run test                                      # all vitest specs
bun run test -- layers/dashboard/features/<name>  # scoped
```

- Specs colocate in `features/**/*.spec.ts` (vitest + happy-dom + `@vue/test-utils`).
- Mock the feature's `api.ts` via `vi.mock('../api.js')` + `vi.hoisted`; see `features/channel-platforms/*.spec.ts`.

## NOTES

- Layer is auto-registered by Nuxt from `layers/` (not in `extends`).
- `@twir/frontend-chat` ships SFC source — it is in `build.transpile` of `web/nuxt.config.ts`.
