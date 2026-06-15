# AGENTS.md — web/layers

Nuxt layer modules for the public website.

## OVERVIEW

Feature-isolated Nuxt layers that compose the public website. Each layer follows the Nuxt 4 layer structure with app code inside an `app/` subdirectory.

## STRUCTURE

Each layer follows the standard Nuxt 4 layer structure:

```
layers/{layer-name}/
├── nuxt.config.ts       # Layer config (required)
├── app/
│   ├── components/      # Auto-merged components
│   ├── composables/     # Auto-merged composables
│   ├── layouts/         # Auto-merged layouts
│   ├── pages/           # Auto-merged routes
│   ├── plugins/         # Auto-merged plugins
│   ├── api/             # Client-side API utilities
│   ├── features/        # Feature modules
│   ├── stores/          # Pinia stores
│   └── ...              # Other app code
├── server/              # Layer-specific server routes (if any)
└── package.json         # Layer dependencies (if any)
```

## LAYERS

| Layer           | Route         | Purpose                  |
| --------------- | ------------- | ------------------------ |
| `landing`       | `/`           | Marketing site           |
| `dashboard`     | `/dashboard/*`| Dashboard application    |
| `url-shortener` | `/s/*`        | Short link redirect page |
| `pastebin`      | `/h/*`        | Code sharing UI          |
| `public`        | `/p/*`        | Public utilities         |
| `overlays`      | `/o/*`        | Browser source overlays  |

## CONVENTIONS

Each layer:

- Has its own `nuxt.config.ts` at the layer root
- App code lives inside `app/` subdirectory
- Server code (if any) lives at layer root in `server/`
- Gets merged by parent `web/nuxt.config.ts`

### Import Paths

Use `~~/layers/{layer}/app/{path}` for absolute imports within layers:

```typescript
import { useSomething } from '~~/layers/dashboard/app/api/something'
import MyComponent from '~~/layers/dashboard/app/components/my-component.vue'
```

Use `#layers/{layer}/app/{path}` for Nuxt layer aliases:

```typescript
import { useSomething } from '#layers/landing/app/api/stats'
```

## USAGE

Layers registered in `web/nuxt.config.ts`:

```typescript
export default defineNuxtConfig({
	extends: [
		"./layers/landing",
		"./layers/dashboard",
		"./layers/url-shortener",
		"./layers/pastebin",
		"./layers/public",
		"./layers/overlays",
	],
});
```

## NOTES

- Layers don't run independently
- Parent config takes precedence
- Shared components go in `web/app/components`
